package proof

import (
 "crypto/sha256"
 "encoding/hex"
 "fmt"
 "os"
 "path/filepath"
 "testing"
)

func TestVerify(t *testing.T) {
 dir:=t.TempDir()
 artifact:=[]byte("ok\n")
 if err:=os.WriteFile(filepath.Join(dir,"test.log"),artifact,0600);err!=nil {t.Fatal(err)}
 sum:=sha256.Sum256(artifact)
 manifest:=filepath.Join(dir,"proof.json")
 write:=func(code int, hash string) {
  t.Helper()
  content:=fmt.Sprintf(`{"schema_version":1,"command":"go test ./...","exit_code":%d,"artifacts":[{"path":"test.log","sha256":"%s"}]}`,code,hash)
  if err:=os.WriteFile(manifest,[]byte(content),0600);err!=nil {t.Fatal(err)}
 }
 write(0,hex.EncodeToString(sum[:]))
 if got:=Verify(manifest);got.Status!="PASS" {t.Fatalf("expected PASS, got %+v",got)}
 write(1,hex.EncodeToString(sum[:]))
 if got:=Verify(manifest);got.Status!="FAIL" {t.Fatalf("expected FAIL, got %+v",got)}
 write(0,fmt.Sprintf("%064x",0))
 if got:=Verify(manifest);got.Status!="FAIL" {t.Fatalf("expected FAIL for hash mismatch, got %+v",got)}
 if err:=os.WriteFile(manifest,[]byte("{}"),0600);err!=nil {t.Fatal(err)}
 if got:=Verify(manifest);got.Status!="INCOMPLETE" {t.Fatalf("expected INCOMPLETE, got %+v",got)}
}
