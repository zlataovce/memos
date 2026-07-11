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

type generationContext struct {
	context.Context

	outDir string
	memos  []*generator.Memo
}

func generate(cmdCtx context.Context, cmd *cli.Command) error {
	outDir := filepath.Clean(cmd.String("out"))
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	inputFiles := cmd.StringArgs("input")
	if len(inputFiles) == 0 {
		inputFiles = []string{"memos"}
	}

	ctx := &generationContext{Context: cmdCtx, outDir: outDir}
	for _, file := range inputFiles {
		if err := ctx.processPath(file); err != nil {
			return fmt.Errorf("failed to process path %s: %w", file, err)
		}
	}

	index, err := generator.GenerateIndex(ctx.memos)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(outDir, "index.html"), []byte(index), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *generationContext) processPath(inFile string) error {
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

			return ctx.processPath(path)
		})
		if err != nil {
			return err
		}
	} else {
		if err := ctx.processFile(inFile); err != nil {
			return fmt.Errorf("failed to process file %s: %w", inFile, err)
		}
	}

	return nil
}

func (ctx *generationContext) processFile(inFile string) error {
	in, err := os.ReadFile(inFile)
	if err != nil {
		return err
	}

	id := strings.TrimSuffix(filepath.Base(inFile), filepath.Ext(inFile))

	memo, err := generator.ParseMemo(id, string(in))
	if err != nil {
		return err
	}

	ctx.memos = append(ctx.memos, memo)

	out, err := memo.Generate()
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(ctx.outDir, id+".html"), []byte(out), 0644)
}
