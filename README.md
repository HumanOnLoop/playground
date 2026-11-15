# playground

Part of HOTL (Human On The Loop). A deliberately small Go module that HOTL's first real agent runs work on: agents
plan, change, test and merge here before they touch a repo that matters. It is public, so that GitHub can protect
`main` with a ruleset on the org's plan.

```sh
go test ./...
go run .                                     # serves on :8080; -addr to change
go build -ldflags "-X main.version=1.2.3" .  # set the version GET /version reports (default "dev")
```

## CI and protection

- [.github/workflows/ci.yml](.github/workflows/ci.yml) is the `ci` check: gofmt -l, go vet and go test. It runs on
  pushes to `main`, `run/**` (run and fix branches) and `revert-**` (GitHub's revert branches, which hold the run branch's slash).
- [ruleset/main.json](ruleset/main.json) protects `main`: `ci` is required (not strict), and there are no force
  pushes or deletions. Apply it with `ruleset/apply.sh` (needs admin rights through `gh`).
- HOTL's CI/CD adapter merges run branches through pull requests once `ci` is green, so merges also pass GitHub's
  own enforcement (ADR 0008).
