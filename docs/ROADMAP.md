# Roadmap

## v0.1 — Evidence manifest verification
- [x] Go CLI skeleton
- [x] Manifest schema v1, SHA-256 artifact checks, explicit status
- [x] Basic unit tests and CI configuration
- [ ] Verify CI test execution and build logs

## v0.2 — Trusted execution evidence
- Bind evidence to commit SHA, workflow identity and run metadata
- Validate signed provenance / attestations and trust policy
- Prevent untrusted manifest authors from self-certifying test execution
- Harden against symlink races, large files and resource exhaustion

## v0.3 — Developer adoption
- GitHub Action and machine-readable reporting
- Python, Node.js and Go test-runner adapters
- Documentation, examples and security review

No feature is marked verified until its tests have run and evidence is recorded.
