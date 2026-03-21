package main

import (
	"code"
	"context"
	"errors"
	"fmt"
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
			if file1 == "" || file2 == "" {
				return errors.New("command requires two arguments")
			}
			result, err := code.GenDiff(file1, file2, cmd.String("format"))
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
