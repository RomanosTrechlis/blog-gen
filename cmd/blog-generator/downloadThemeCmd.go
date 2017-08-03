package main

import (
	"flag"
	"github.com/RomanosTrechlis/blog-generator/datasource"
)

const downloadThemeShortHelp = `Downloads theme from given datasource`
const downloadThemeLongHelp = `
Downloads theme from a given datasource inside the config.json file.
The folowing part of config.json controls the behavior of "fetch-posts"
command.

"ThemeFolder": "./static/",
"Theme": {
    "Type": "git",
    "Repository": "https://github.com/RomanosTrechlis/BlogThemeBlueSimple.git"
},

The "Type" can also be "local" and the "Repository" a local folder.
The "ThemeFolder" is were the static pages of the theme will be
cloned for use in the blog generation phase.
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
	ds, err := datasource.New(cmd.sourceType)
	if err != nil {
		ctx.Err.Fatal("please provide a datasource in the configuration file:", err)
		return err
	}

	_, err = ds.Fetch(cmd.source, cmd.destination)
	if err != nil {
		ctx.Err.Fatal("failure to fetch theme:", err)
	}
	return err
}
