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
	Decompress bool     `help:"Unseal and decompress the input (auto-triggered if file ends in .bmx)" short:"d"`
	Width      int      `help:"Column width for text wrapping during sealing (default: 80)" short:"w" default:"80"`
	Inputs     []string `arg:"" optional:"" help:"Paths to input files (reads from stdin if omitted)"`
}

func main() {
	var cli CLI
	kong.Parse(&cli,
		kong.Name("bmx"),
		kong.Description("Pack or unpack text files using a compressed Z85 validation envelope."),
		kong.UsageOnError(),
	)

	if len(cli.Inputs) == 0 {
		processStdin(cli.Decompress, cli.Width)
		return
	}

	for _, inputPath := range cli.Inputs {
		processFile(inputPath, cli.Decompress, cli.Width)
	}
}

func processStdin(decompress bool, width int) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	content := string(data)

	if strings.TrimSpace(content) == "" {
		return
	}

	shouldDecompress := decompress || strings.HasPrefix(content, bmx.Header)

	var result string
	var err2 error
	if shouldDecompress {
		result, err2 = bmx.UnsealText(content)
	} else {
		result, err2 = bmx.SealText(content, width)
	}
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err2)
		os.Exit(1)
	}

	fmt.Print(result)
	if !strings.HasSuffix(result, "\n") {
		fmt.Println()
	}
}

func processFile(inputPath string, decompress bool, width int) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	content := string(data)

	if strings.TrimSpace(content) == "" {
		return
	}

	shouldDecompress := decompress ||
		filepath.Ext(inputPath) == ".bmx" ||
		strings.HasPrefix(content, bmx.Header)

	if !shouldDecompress && filepath.Ext(inputPath) == ".bmx" {
		fmt.Fprintf(os.Stderr, "Error: File is already a .bmx matrix. Aborting to avoid double compression.\n")
		os.Exit(1)
	}

	var result string
	if shouldDecompress {
		result, err = bmx.UnsealText(content)
	} else {
		result, err = bmx.SealText(content, width)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

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
}
