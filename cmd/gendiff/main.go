package main

import (
	"code/files"
	"context"
	"errors"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "gendiff",
		Usage: "Compares two configuration files and shows a difference.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "stylish",
				Usage:   "output format",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			file1 := cmd.Args().Get(0)
			file2 := cmd.Args().Get(1)
			if file2 == "" {
				return errors.New("requires two arguments, but only one was passed")
			}
			if err := files.ParseFileContent(file1); err != nil {
				return err
			}
			if err := files.ParseFileContent(file2); err != nil {
				return err
			}
			return nil
		},
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
