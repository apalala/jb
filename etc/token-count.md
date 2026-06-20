Here is the table formatted in clean, standard Markdown:

| Name | Brief Description | Installation Command |
| --- | --- | --- |
| **`cl100k`** | An ultra-lightweight, standalone Go-based CLI implementation for counting tokens against modern OpenAI vocabularies. | `port install cl100k` *(MacPorts)*<br>`brew install cl100k` *(Homebrew)* |
| **`ttok`** | A dedicated, production-grade CLI token counter and text truncator that wraps `tiktoken`. Perfect for piping data. Also available as a Python-based utility. | `brew install simonw/ttok/ttok` *(Homebrew)*<br>`pipx install ttok`<br>`pip install --user ttok` *(Pip)* |
| **`tiktoken` (CLI)** | OpenAI's official, blazing-fast BPE tokenizer utility. Ideal for gpt-4o, gpt-4, and cl100k_base models. Excellent for script integration but less tailored for terminal piping than `ttok`. | `pip install tiktoken` *(Custom/Pip)* |
| **`llm` (with ttok plugin)** | A comprehensive CLI prompt manager that can count, slice, and pipe text through its tokenization sub-commands. | `brew install llm` *(Homebrew)*<br>`llm install llm-token-counter` |
| **`wc -w`** | The native Unix word counter. **Warning:** It counts white-space words, not BPE tokens, but serves as a quick, zero-install rough approximation (typically 1 word $\approx$ 1.3 tokens). | *Pre-installed natively on macOS* |
| **`tc` (token-counter)** | A lightweight, blazing-fast Unix-style `wc` equivalent for tokens written in Rust. Uses the HuggingFace tokenizers engine to count local files or `stdin`. | `cargo install token-counter` *(Custom/Cargo)* |
| **`rtk` (Rust Token Killer)** | A specialized developer CLI proxy that sits between your shell and agent tools (like Claude Code) to strip and compress token waste. | `brew install rtk-ai/tap/rtk` *(Homebrew)* |
