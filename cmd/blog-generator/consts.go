package main

const (
	getPostsShortHelp = `Downloads posts from given datasource`
	getPostsLongHelp  = `
Downloads posts from a given datasource inside the config.json file.

The folowing part of config.json controls the behavior of "fetch-posts"
command.

"DataSource": {
    "Type": "git",
    "Repository": "https://github.com/RomanosTrechlis/blog.git"
},
"TempFolder": "./tmp"

The "Type" can also be "local" and the "Repository" local folder.
The "TempFolder" is were the posts will be cloned for generation.
`

	getThemeShortHelp = `Downloads theme from given datasource`
	getThemeLongHelp  = `
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

	generateShortHelp = `Generates blog from existing resources`
	generateLongHelp  = `
Generates blog from existing resources. Before this command runs the
fetch-posts and fetch-theme should run.

Generate requires most of the fields inside config.json.

The following three specify where the artifacts for the generation of
the blog are and where to put the generated files.
"DestFolder": "./public",
"TempFolder": "./tmp",
"ThemeFolder": "./static/"

The following field specifies how many posts will be per page.
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

	jsonExampleShortHelp = `Run: "blog-generator json-exmple" to see a config.json example`
	jsonExampleLongHelp  = `
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
  	"Type": "git",
    "URL": "https://github.com/RomanosTrechlis/romanostrechlis.github.io.git",
    "Username": "RomanosTrechlis",
    "Password": ""
  }
}
`

	runShortHelp = `Runs a web server for the generated blog`
	runLongHelp  = `
Runs a web server for the generated blog.

The default port for the server is 8080.
`

	execAllShortHelp = `Sequential execution of all commands`
	execAllLongHelp = `Sequential execution of all commands.

Downloads posts, and theme, then generates the blog and runs a web server.
`
)
