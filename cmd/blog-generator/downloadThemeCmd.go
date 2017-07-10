package main

import (
	"flag"
	"github.com/RomanosTrechlis/blog-generator/datasource"
)

const downloadThemeShortHelp = `Downloads theme from given datasource`
const downloadThemeLongHelp = `
Downloads theme from a given datasource inside the config.json file.

`

type downloadThemeCmd struct {
	sourceType, source, destination string
}

func (cmd *downloadThemeCmd) Name() string      { return "fetch-theme" }
func (cmd *downloadThemeCmd) Args() string      { return "" }
func (cmd *downloadThemeCmd) ShortHelp() string { return downloadThemeShortHelp }
func (cmd *downloadThemeCmd) LongHelp() string  { return downloadThemeLongHelp }
func (cmd *downloadThemeCmd) Hidden() bool      { return false }

func (cmd *downloadThemeCmd) Register(fs *flag.FlagSet) {
}

func (cmd *downloadThemeCmd) Run(ctx *Ctx, args []string) error {
	var err error
	switch cmd.sourceType {
	case "git":
		ds := datasource.NewGitDataSource()
		_, err = ds.Fetch(cmd.source,
			cmd.destination)
	case "local":
		ds := datasource.NewLocalDataSource()
		_, err = ds.Fetch(cmd.source,
			cmd.destination)
	case "":
		ctx.Err.Fatal("please provide a datasource in the configuration file")
	}

	return err
}
