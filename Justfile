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

# ─── Aliases ────────────────────────────────────────────────────────────────

alias fmt := ruff-format
alias lint := ruff-lint
