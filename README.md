# PAMILANGA Proof Engine

**Evidence-aware verification for software test artifacts.**

Proof Engine is an experimental open-source Go CLI for checking whether a supplied test-evidence manifest is structurally complete and whether referenced artifact hashes match the local files.

> **Security boundary:** A matching SHA-256 hash does **not** prove that a test actually ran, that the recorded exit code is truthful, or that the manifest comes from a trusted CI runner. Version 0.1 only verifies self-reported evidence consistency. Trusted provenance is a planned feature.

## Quick start

Requires Go 1.22+.

```bash
go test ./...
go run ./cmd/proof-engine verify --manifest proof.json
```

Example `proof.json` (replace the digest with the actual SHA-256 of `test.log`):

```json
{
  "schema_version": 1,
  "command": "go test ./...",
  "exit_code": 0,
  "artifacts": [
    { "path": "test.log", "sha256": "<64-character SHA-256 hex digest>" }
  ]
}
```

Possible statuses: `PASS` (self-reported manifest and local artifacts are consistent), `FAIL` (reported failure or integrity violation), `INCOMPLETE` (required evidence absent), `UNVERIFIED` (manifest unavailable).

**Do not use a v0.1 PASS result as proof of trusted test execution.**

## Roadmap

See [docs/ROADMAP.md](docs/ROADMAP.md). Planned: CI provenance, signed attestations, GitHub Action, test-runner adapters and adversarial verification.

## Development

```bash
go test ./...
go vet ./...
go build ./cmd/proof-engine
```

## Project

Maintained as an open-source PAMILANGA LIMITED initiative. Contributions and security reports are welcome; please avoid publishing exploit details in public issues.

MIT License.
