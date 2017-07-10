package main

import (
	"flag"
	"github.com/RomanosTrechlis/blog-generator/util/fs"
	"github.com/RomanosTrechlis/blog-generator/generator"
	"github.com/RomanosTrechlis/blog-generator/config"
)

const generateShortHelp = `Generates blog from existing resources`
const generateLongHelp = `
Generates blog from existing resources. Before this command runs the
fetch-posts and fetch-theme should run.

`

type generateCmd struct {
	siteInfo *config.SiteInformation
}

func (cmd *generateCmd) Name() string      { return "generate" }
func (cmd *generateCmd) Args() string      { return "" }
func (cmd *generateCmd) ShortHelp() string { return generateShortHelp }
func (cmd *generateCmd) LongHelp() string  { return generateLongHelp }
func (cmd *generateCmd) Hidden() bool      { return false }

func (cmd *generateCmd) Register(fs *flag.FlagSet) {
}

func (cmd *generateCmd) Run(ctx *Ctx, args []string) error {
	dirs, err := fs.GetContentFolders(cmd.siteInfo.TempFolder)
	if err != nil {
		return err
	}
	g := generator.NewSiteGenerator(&generator.SiteConfig{
		Sources:  dirs,
		SiteInfo: cmd.siteInfo,
	})

	err = g.Generate()
	return err
}
