package proof

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Artifact struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}
type Manifest struct {
	SchemaVersion int        `json:"schema_version"`
	Command       string     `json:"command"`
	ExitCode      *int       `json:"exit_code"`
	Artifacts     []Artifact `json:"artifacts"`
}
type Result struct {
	Status   string   `json:"status"`
	Checks   []string `json:"checks"`
	Warnings []string `json:"warnings,omitempty"`
}

func Verify(manifestPath string) Result {
	r := Result{Status: "UNVERIFIED", Checks: []string{}}
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		r.Checks = append(r.Checks, "manifest unreadable: "+err.Error())
		return r
	}
	var m Manifest
	if err = json.Unmarshal(data, &m); err != nil {
		r.Status = "FAIL"
		r.Checks = append(r.Checks, "invalid JSON: "+err.Error())
		return r
	}
	if m.SchemaVersion != 1 || strings.TrimSpace(m.Command) == "" || m.ExitCode == nil || len(m.Artifacts) == 0 {
		r.Status = "INCOMPLETE"
		r.Checks = append(r.Checks, "required evidence fields missing or unsupported schema")
		return r
	}
	if *m.ExitCode != 0 {
		r.Status = "FAIL"
		r.Checks = append(r.Checks, fmt.Sprintf("reported command exit code: %d", *m.ExitCode))
		return r
	}
	base := filepath.Dir(manifestPath)
	for _, a := range m.Artifacts {
		if a.Path == "" || len(a.SHA256) != 64 {
			r.Status = "INCOMPLETE"
			r.Checks = append(r.Checks, "artifact path or SHA-256 missing")
			return r
		}
		expected, err := hex.DecodeString(a.SHA256)
		if err != nil || len(expected) != sha256.Size {
			r.Status = "INCOMPLETE"
			r.Checks = append(r.Checks, "invalid SHA-256 format")
			return r
		}
		if filepath.IsAbs(a.Path) {
			r.Status = "FAIL"
			r.Checks = append(r.Checks, "absolute artifact path prohibited")
			return r
		}
		cleaned := filepath.Clean(a.Path)
		if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
			r.Status = "FAIL"
			r.Checks = append(r.Checks, "artifact path escapes manifest directory")
			return r
		}
		full := filepath.Join(base, cleaned)
		info, err := os.Lstat(full)
		if err != nil || !info.Mode().IsRegular() {
			r.Status = "FAIL"
			r.Checks = append(r.Checks, "artifact missing or not regular: "+a.Path)
			return r
		}
		content, err := os.ReadFile(full)
		if err != nil {
			r.Status = "FAIL"
			r.Checks = append(r.Checks, "artifact unreadable: "+a.Path)
			return r
		}
		sum := sha256.Sum256(content)
		if !strings.EqualFold(hex.EncodeToString(sum[:]), a.SHA256) {
			r.Status = "FAIL"
			r.Checks = append(r.Checks, "hash mismatch: "+a.Path)
			return r
		}
		r.Checks = append(r.Checks, "SHA-256 matched: "+a.Path)
	}
	r.Status = "PASS"
	r.Warnings = []string{"PASS means manifest consistency and artifact integrity only; execution provenance and trust are not independently verified."}
	return r
}
