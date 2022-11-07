package cli

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/lumisxh/schef/pkg/config"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"time"
)

func RunCli() {
	app := &cli.App{
		Name:  "Schef",
		Usage: "Schef - the new language agnostic monorepo tool",
		Action: func(*cli.Context) error {

			c := color.New(color.FgCyan, color.OpBold, color.OpBlink)
			c.Println("Schef - the new language agnostic monorepo tool")
			color.Info.Prompt("Type `schef --help`  to see all possible commands and options")
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "member",
				Usage:    "Monorepo member to perform actions on",
				Required: false,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "build",
				Aliases: []string{"b"},
				Usage:   "Builds the entire monorepo (or use flag --member to specify one monorepo member)",
				Action: func(cCtx *cli.Context) error {
					color.Notice.Prompt("Starting build process...\n")

					start := time.Now()

					config.Glue(cCtx.String("member"), "build")

					color.Notice.Prompt("Build process finished. Execution time: %s", time.Since(start))
					return nil
				},
			},
			{
				Name:    "dev",
				Aliases: []string{"d"},
				Usage:   "Executes the dev commands in the entire monorepo (or use flag --member to specify one monorepo member)",
				Action: func(cCtx *cli.Context) error {
					color.Notice.Prompt("Starting dev process...\n")

					start := time.Now()

					config.Glue(cCtx.String("member"), "dev")

					color.Notice.Prompt("Dev process finished. Execution time: %s", time.Since(start))
					return nil
				},
			},
			{
				Name:    "clear",
				Aliases: []string{"c"},
				Usage:   "Clears the internal cache",
				Action: func(cCtx *cli.Context) error {
					dir, userErr := os.UserHomeDir()
					if userErr != nil {
						color.Error.Block(userErr.Error())
						os.Exit(1)
					}
					conf, err := filepath.Abs(fmt.Sprintf("%s/.schef.json", dir))
					if err != nil {
						color.Error.Block(userErr.Error())
						os.Exit(1)

					}

					if osErr := os.WriteFile(conf, []byte("{}"), 0644); osErr != nil {
						color.Warn.Block("%s\n%s", "The cache could not be cleared", osErr.Error())
					}

					color.Notice.Prompt("Schef cache was cleared successfully")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		color.Error.Block(err.Error())
	}
}
