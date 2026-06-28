lint: ruff-lint vet

fmt: ruff-format gofmt


ruff-lint *args="":
    uv run ruff check {{args}}

ruff-format *args="":
    uv run ruff format {{args}}

ruff *args="":
    uv run ruff check {{args}}
    uv run ruff format --check {{args}}

pytest *args="":
    uv run pytest {{args}}

test:
    uv run pytest jb

git:
    go build -o bin/git cmd/git/safe_git.go

bash:
    go build -o bin/bash cmd/bash/log_bash.go

python3:
    go build -o bin/python3 cmd/python/safe_python.go

head:
    go build -o bin/head cmd/head/safe_head.go

bmx:
    go build -o bin/bmx cmd/bmx/bmx.go

libz:
    test -f lib/zlib/libz.a || (cd lib/zlib && ./configure && make)

czlib-bmx: libz
    go build -tags czlib -o bin/bmx-czlib cmd/bmx/bmx.go

czlib-test *args="": libz
    go test -tags czlib ./pkg/bmx/... {{args}}

jb:
    go build -o bin/jb cmd/jb/jb.go

cmd: python3 bash git head bmx jb

# === Go ===

gofmt:
    gofmt -l -w -s ./cmd ./pkg

golint: gofmt vet
    golangci-lint run ./... --exclude-dirs ./tmp

deps:
    go mod download

vendor: tidy
    go mod vendor -o _vendor

gotest *args="":
    go test ./pkg/... {{args}}

vet:
    go vet -structtag=false ./...

tidy:
    go mod tidy

update:
    go get -u ./...
    go mod tidy

tools:
    go install golang.org/x/tools/cmd/goimports@latest
    go install gotest.tools/gotestsum@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8
