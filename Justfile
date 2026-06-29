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

test *args="":
    uv run pytest {{args}}

git:
    go build -o target/debug/git cmd/git/main.go

release-git:
    mkdir -p target/release
    go build -ldflags="-s -w" -o target/release/git cmd/git/main.go

bash:
    go build -o target/debug/bash cmd/bash/main.go

release-bash:
    mkdir -p target/release
    go build -ldflags="-s -w" -o target/release/bash cmd/bash/main.go

python3:
    go build -o target/debug/python3 cmd/python/main.go

release-python3:
    mkdir -p target/release
    go build -ldflags="-s -w" -o target/release/python3 cmd/python/main.go

head:
    go build -o target/debug/head cmd/head/main.go

release-head:
    mkdir -p target/release
    go build -ldflags="-s -w" -o target/release/head cmd/head/main.go

safe:
    go build -o target/debug/safe ./cmd/safe

release-safe:
    mkdir -p target/release
    go build -ldflags="-s -w" -o target/release/safe ./cmd/safe

bmx:
    go build -o target/debug/bmx cmd/bmx/bmx.go

release-bmx:
    mkdir -p target/release
    go build -ldflags="-s -w" -o target/release/bmx cmd/bmx/bmx.go

libz:
    test -f lib/zlib/libz.a || (cd lib/zlib && ./configure && make)

czlib-bmx: libz
    go build -tags czlib -o target/debug/bmx-czlib cmd/bmx/bmx.go

czlib-test *args="": libz
    go test -tags czlib ./pkg/bmx/... {{args}}

jb:
    go build -o target/debug/jb ./cmd/jb

release-jb:
    mkdir -p target/release
    go build -ldflags="-s -w" -o target/release/jb ./cmd/jb

build: lint vet python3 bash git head bmx jb safe

release: release-python3 release-bash release-git release-head release-bmx release-jb release-safe

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

vet: gofmt
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
