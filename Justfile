# ─── Development Commands ───────────────────────────────────────────────────

# Run ruff linter (check)
ruff-lint *args="":
    uv run ruff check {{args}}

# Run ruff formatter
ruff-format *args="":
    uv run ruff format {{args}}

# Run both ruff lint and format (in check mode only)
ruff *args="":
    uv run ruff check {{args}}
    uv run ruff format --check {{args}}

# Run tests with pytest
test *args="":
    uv run pytest {{args}}

git:
    go build -o bin/git cmd/git/safe_git.go

bash:
    go build -o bin/bash cmd/bash/log_bash.go

python3:
    go build -o bin/python3 cmd/python/safe_python.go

cmd: python3 bash git

# ─── Aliases ────────────────────────────────────────────────────────────────

alias fmt := ruff-format
alias lint := ruff-lint

gofmt:
    gofmt -l -w -s ./cmd ./pkg

golint: gofmt vet
    golangci-lint run ./... --exclude-dirs ./tmp

deps:
    go mod download

vendor: tidy
    go mod vendor -o _vendor

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
