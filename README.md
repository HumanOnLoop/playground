# playground

Part of HOTL (Human On The Loop). A deliberately small Go module that HOTL's first real agent runs work on: agents
plan, change, test and merge here before they touch a repo that matters.

```sh
go test ./...
go run .                     # serves GET /version on :8080 (-addr to change)
go build -ldflags "-X main.version=1.2.3" .
```
