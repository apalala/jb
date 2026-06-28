package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/apalala/jb/pkg/bmx"
)

type CLI struct {
	Decompress bool   `help:"Unseal and decompress the input (auto-triggered if file ends in .bmx)" short:"d"`
	Width      int    `help:"Column width for text wrapping during sealing (default: 80)" short:"w" default:"80"`
	Input      string `arg:"" optional:"" help:"Path to the input file (reads from stdin if omitted)"`
	Output     string `arg:"" optional:"" help:"Path to the output file (writes to stdout or replaces input if omitted)"`
}

func main() {
	var cli CLI
	kong.Parse(&cli,
		kong.Name("bmx"),
		kong.Description("Pack or unpack text files using a compressed Z85 validation envelope."),
		kong.UsageOnError(),
	)

	var inputPath string
	if cli.Input != "" {
		inputPath = cli.Input
	}

	var content string
	if inputPath != "" {
		data, err := os.ReadFile(inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		content = string(data)
	} else {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		content = string(data)
	}

	if strings.TrimSpace(content) == "" {
		return
	}

	shouldDecompress := cli.Decompress ||
		(inputPath != "" && filepath.Ext(inputPath) == ".bmx") ||
		strings.HasPrefix(content, bmx.Header)

	if !shouldDecompress && inputPath != "" && filepath.Ext(inputPath) == ".bmx" {
		fmt.Fprintf(os.Stderr, "Error: File is already a .bmx matrix. Aborting to avoid double compression.\n")
		os.Exit(1)
	}

	var result string
	var err error
	if shouldDecompress {
		result, err = bmx.UnsealText(content)
	} else {
		result, err = bmx.SealText(content, cli.Width)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if cli.Output != "" {
		if err := os.WriteFile(cli.Output, []byte(result), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else if inputPath != "" {
		var outPath string
		if shouldDecompress {
			if filepath.Ext(inputPath) == ".bmx" {
				outPath = strings.TrimSuffix(inputPath, filepath.Ext(inputPath))
			} else {
				dir := filepath.Dir(inputPath)
				base := filepath.Base(inputPath)
				outPath = filepath.Join(dir, base+".out")
			}
		} else {
			dir := filepath.Dir(inputPath)
			base := filepath.Base(inputPath)
			outPath = filepath.Join(dir, base+".bmx")
		}
		if err := os.WriteFile(outPath, []byte(result), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		os.Remove(inputPath)
	} else {
		fmt.Print(result)
		if !strings.HasSuffix(result, "\n") {
			fmt.Println()
		}
	}
}
