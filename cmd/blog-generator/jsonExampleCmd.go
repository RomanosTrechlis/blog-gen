package main

import "flag"

const jsonExampleShortHelp = `Run: "blog-generator json-exmple" to see a config.json example`
const jsonExampleLongHelp = `
{
  "Author": "Romanos Trechlis",
  "BlogURL": "romanostrechlis.github.io",
  "BlogLanguage": "en-us",
  "BlogDescription": "This is my personal blog.",
  "DateFormat": "2006-01-02 15:04:05",
  "Theme": {
    "Type": "git",
    "Repository": "https://github.com/RomanosTrechlis/BlogThemeBlueSimple.git"
  },
  "BlogTitle": "Romanos-Antonios Trechlis",
  "NumPostsFrontPage": 10,
  "DataSource": {
    "Type": "git",
    "Repository": "https://github.com/RomanosTrechlis/blog.git"
  },
  "DestFolder": "./public",
  "TempFolder": "./tmp",
  "ThemeFolder": "./static/",
  "StaticPages": [
    {
      "File": "favicon.ico",
      "To": "favicon.ico",
      "IsTemplate": false
    },
    {
      "File": "robots.txt",
      "To": "robots.txt",
      "IsTemplate": false
    },
    {
      "File": "about.png",
      "To": "about.png",
      "IsTemplate": false
    },
    {
      "File": "style.min.css",
      "To": "style.min.css",
      "IsTemplate": false
    },
    {
      "File": "google.min.css",
      "To": "google.min.css",
      "IsTemplate": false
    },
    {
      "File": "about.html",
      "To": "about/index.html",
      "IsTemplate": true
    }
  ],
  "Upload": {
    "URL": "https://github.com/RomanosTrechlis/romanostrechlis.github.io.git",
    "Username": "RomanosTrechlis",
    "Password": ""
  }
}
`


type jsonExampleCmd struct {
}

func (cmd *jsonExampleCmd) Name() string      { return "json-example" }
func (cmd *jsonExampleCmd) Args() string      { return "" }
func (cmd *jsonExampleCmd) ShortHelp() string { return jsonExampleShortHelp }
func (cmd *jsonExampleCmd) LongHelp() string  { return jsonExampleLongHelp }
func (cmd *jsonExampleCmd) Hidden() bool      { return false }

func (cmd *jsonExampleCmd) Register(fs *flag.FlagSet) {
}

func (cmd *jsonExampleCmd) Run(ctx *Ctx, args []string) error {
	ctx.Out.Print(jsonExampleLongHelp)
	return nil
}
