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

Generate requires most of the fields inside config.json.

The following three specify where the artifacts for the generation of
the blog are and where to put the generated files.
"DestFolder": "./public",
"TempFolder": "./tmp",
"ThemeFolder": "./static/"

The folowing field specifies how many posts will be per page.
"NumPostsFrontPage": 10,

The following snippet specifies the static pages and other artifacts like .css
.js images etc to be copied, or generated as templates, but are not posts.
"StaticPages": [
    {
      "File": "favicon.ico",
      "To": "favicon.ico",
      "IsTemplate": false
    },
    {
      "File": "about.html",
      "To": "about/index.html",
      "IsTemplate": true
    }
}

Finally, the following fields contain information about the site.
"Author": "Romanos Trechlis",
"BlogURL": "romanostrechlis.github.io",
"BlogLanguage": "en-us",
"BlogDescription": "This is my personal blog.",
"DateFormat": "2006-01-02 15:04:05",

To see a config.json example run: blog-generator json-example
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
