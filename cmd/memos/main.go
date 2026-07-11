package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
	"go.kcra.me/memos/generator"
)

func main() {
	cmd := &cli.Command{
		Name:  "memos",
		Usage: "generate a blog-ish static site out of a bunch of Markdown files",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "out",
				Value: "docs",
				Usage: "the generated output directory",
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name: "input",
				Max:  -1,
			},
		},
		Action: generate,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func generate(_ context.Context, cmd *cli.Command) error {
	outDir := filepath.Clean(cmd.String("out"))
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	inputFiles := cmd.StringArgs("input")
	if len(inputFiles) == 0 {
		inputFiles = []string{"memos"}
	}

	for _, file := range inputFiles {
		if err := processPath(file, outDir); err != nil {
			return fmt.Errorf("failed to process path %s: %w", file, err)
		}
	}

	return nil
}

func processPath(inFile, outDir string) error {
	inFile = filepath.Clean(inFile)

	stat, err := os.Stat(inFile)
	if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", inFile, err)
	}

	if stat.IsDir() {
		err := filepath.WalkDir(inFile, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if path == inFile {
				return nil // skip root
			}

			return processPath(path, outDir)
		})
		if err != nil {
			return err
		}
	} else {
		if err := processFile(inFile, outDir); err != nil {
			return fmt.Errorf("failed to process file %s: %w", inFile, err)
		}
	}

	return nil
}

func processFile(inFile, outDir string) error {
	in, err := os.ReadFile(inFile)
	if err != nil {
		return err
	}

	memo, err := generator.ParseMemo(string(in))
	if err != nil {
		return err
	}

	out, err := memo.Generate()
	if err != nil {
		return err
	}

	name := strings.TrimSuffix(filepath.Base(inFile), filepath.Ext(inFile))
	return os.WriteFile(filepath.Join(outDir, name+".html"), []byte(out), 0644)
}
