# static-site-generator
For this project I used the core ideas of [Mario Zupan's](https://zupzup.org/) blog generator and altered them to fit my needs.

Zupzup explains everythin in [this](https://zupzup.org/static-blog-generator-go/) blog post, so I will continue from where he stopped.

## 1. Use a configuration file.

Even though his implementation is solid, changing variables in the code in order to produce a different blog and then compiling the project again wasn't efficient for me.

Additionaly, because I change blog posts on my computer and then push them to git I've implemented a local variation of datasource so I can use that local folder instead of a git repository. 

```go
    switch config.SiteInfo.DataSource.Type {
	case "git":
		ds := datasource.NewGitDataSource()
		_, err = ds.Fetch(config.SiteInfo.DataSource.Repository,
			config.SiteInfo.TempFolder)
	case "local":
		ds := datasource.NewLocalDataSource()
		_, err = ds.Fetch(config.SiteInfo.DataSource.Repository,
			config.SiteInfo.TempFolder)
	case "":
		log.Fatal("please provide a datasource in the configuration file")
	}
```

The configuration file looks like this:

```javascript
{
  "Author": "Romanos Trechlis",
  "BlogURL": "romanostrechlis.github.io",
  "BlogLanguage": "en-us",
  "BlogDescription": "Desc",
  "DateFormat": "2006-01-02 15:04:05",
  "ThemePath": "./static/",
  "BlogTitle": "RTB",
  "NumPostsFrontPage": 10,
  "DataSource": {
    "Type": "git",
    "Repository": "https://github.com/RomanosTrechlis/blog.git"
  },
  "TempFolder": "./tmp",
  "DestFolder": "./public"
}
```

The DataSource Type can also be *local* and the Repository can be a folder

```javascript
{
    "Type": "local",
    "Repository": "C:/Users/Romanos/Desktop/testLocal/blog/"
}

## 2. Use of cli to break functionality

Another issue I had was that I wish to download the blog content once and then generate the site multiple times with different templates. 

I also wish to see what my blog would look like before I push it. 

In order to achieve that I used flags and broke the functionality to different steps. So, if I want to generate, download and see the results I will run the command:

```bash
site-generator -fetch -generate -run
```

This will download my content, generate the site and run a local server.

## 3. Paging

Mario shows a fixed number of blog posts on the frontpage and uses archive to show all the blog posts.

For my blog I wished to show a fixed number of posts and give the visitor the ability to navigate to the next page with the next fixed number of blog posts.

I did that using the following code:

```go
    // frontpage
	paging := config.SiteInfo.NumPostsFrontPage
	numOfPages := getNumberOfPages(posts)
	for i := 0; i < numOfPages; i++ {
		to := destination
		if i != 0 {
			to = fmt.Sprintf("%s/%d", destination, i+1)
		}
		toP := (i + 1) * paging
		if (i + 1) == numOfPages {
			toP = len(posts)
		}
		generators = append(generators, &ListingGenerator{&ListingConfig{
			Posts:       posts[i*paging : toP],
			Template:    t,
			Destination: to,
			PageTitle:   "",
			PageNum:     i + 1,
			MaxPageNum:  numOfPages,
		}})
	}
```

and in the template:

```html
    <div id="paging">
        {{if ne .PageNum 0}}
        {{ if ne .PrevPageNum 0}}
        <a href="/{{if ne .PrevPageNum 1}}{{.PrevPageNum}}/{{end}}">
            <!-- Less Than icon by Icons8 -->
            <img style="vertical-align:middle" class="icon icons8-Less-Than" width="30" height="30" src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABoAAAAaCAYAAACpSkzOAAAA+0lEQVRIS9WVyQ3CMBBFv6EAoBcER9IFKQMqYOgmEU2EI8uBUkIBjpGRIkUm3uJJJHy1PM//zyYw0RETcfCfoHVOmZCoHxd6mU6xKNKAOXACkCmF+l7SihXUBXQD3wr6ETBIkQ2gYUrhfS9pmaTIBfhCgKuSOAzOUQigAehZUGVrF6d1HIAW3AviBDhBmz3VQmBh2qBz4LMoyrptTqr7IAXgVGSCAFQSOLuS7ZuZvTmyWZcCtBbDDNDtvbP8NFqht7y5gEEjSJd7KjAI1NoXAmwkjoNHkJknF3C0NdFnKduacClkWRO+phx9lfs+oO+jqi4kYNRQTQloe/sBjSWMG74ndYYAAAAASUVORK5CYII=">
        </a>
        {{end}}
        <span><strong>{{.PageNum}}</strong></span>
        {{ if ne .NextPageNum 0 }}
        <a href="/{{.NextPageNum}}/">
            <!-- More Than icon by Icons8 -->
            <img style="vertical-align:middle" class="icon icons8-More-Than" width="30" height="30" src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABoAAAAaCAYAAACpSkzOAAABBElEQVRIS8XWyxGCMBAG4D9SgNgBJXjwdZMmHClDK3DtBsYm8OjjYCl4F+LgDDOQCSGQHUwB+fh3NwkCIy0xkoP/QpuIgg8QPGNKuRJrE632lAmBKYA0B84coBZaRySVJM5gVyK1coNBLbSIKJwAJIBtS496g8ap4wStxpsDtIKq8tmARY7j40IvteS9IBtQSmT3hGYsUBd4i8s5aq5BieqQB5wAhPVt2aCyVzqgxKTE+56Q75TIBPwQ4CpzHAYPgw1QAGS6EzsPbFuJqgRdQFXC1iuICzBCtWei0dOyB7YJrIZBfSZcAKtEHIARWu5oLj34HC+rEeL6T6jv43QF9fmg0aAvNqiLGy9Yk3kAAAAASUVORK5CYII=">
        </a>
        {{end}}
        {{end}}
    </div>
```

This code creates numbered folders that contain a new page with the next fixed amount of blog posts.

## 4. Share buttons

I also like share buttons to my blog posts. So, here is the template code:

```html
    {{if .IsPost}}
    <div id="share-buttons">
        <div class="share-button share-text">
            <span>Share</span>
        </div>
        <div class="share-button share-button-facebook" data-share-url="{{.URL}}">
            <div class="box">
                <a href="https://www.facebook.com/sharer/sharer.php?u={{.URL}}">
                    <!-- Facebook icon by Icons8 -->
                    <img class="icon icons8-Facebook" width="48" height="48" src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAABqElEQVRoQ2NkGOKAcYi7n2HUAwMdg6MxMOhjwDZgWwATw798BkZGB7o69v//A/8YmCYe3uC1AZ+9eJOQfcDW+QyMDAl0dTiaZf//MzQe2uDdgMsNOD0ADnnG/+sH0vEwu//+/+d4ZIPvAWxuwekBu8CtBxgZGOwHgwf+MzAcPLTeG2sSxukB+8Ct/weD42FuOLjeG6tbRz1Ar1gaNDHg4SjD4O4ozWCoI4zhd4egbTjDY1B4oKXCmMHGTBynIwe1B0AhX5GrhzfFDWoPTGi2YDDQFhq6HjiwzgvD8TUdZxmOnHpJVDkw4HkAmwfwJRl0X416gKh4xqNoSMYANkfj8uOXb38YfGJ2Da56gBQPXLj6jqGg9sTQ9cCaLQ8Ypsy7NnQ9sGDlbQYQxgVGS6ERWQohe3pIFqOjHkAKgdFMPJqJsTSnR1ujDAy45weoPS40WowO9MDWaAyQGgNDfnB3yA+vgyouu4AtCxgZGeMprcQo0U/2BAfMUpuAzQ5MjEwN9J4rAM0J/Pv/rwHXxAbMfaOzlJQkD2roHY0BaoQiJWaMxgAloUcNvQDVMwpAikLd1gAAAABJRU5ErkJggg==">
                </a>
            </div>
        </div>

        <div class="share-button share-button-twitter" data-share-url="{{.URL}}">
            <div class="box">
                <a href="http://twitter.com/intent/tweet?source=sharethiscom&text={{.PageTitle}}&url={{.URL}}&via=r_trechlis">
                    <!-- Twitter icon by Icons8 -->
                    <img class="icon icons8-Twitter" width="48" height="48" src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAADTklEQVRoQ+1Y3VHbQBD+9kTMZCwGdxB3EFEBTgWBFwxPmAowFQAVxB3gPGH7BaeCKBXEqQDTAYwFCZm528yd+XGwJZ1OAoYZ6VXan+/bb/f2RHjjD73x/FECeO0KlhUoK5CTgVJCOQnMbf6yFTj9XRdC7gNoEBDo7BkYAQhVpXqMTbqcQzSIAmz5+puFTzqA00kDOythXqpE//qQwEdxfhh8SUQHcsvv4oxruL1Z98BtJlyopt9yA3DGNfE3OgfTUG37e64gRC86IUJsEv/5ZYQgNKbV4StV8esLK3NnlFyB00nDE/TdOGN0XUB4g6gFxklW8CZ5ogbeVcf4EwVxKkgEIPqTIwIdPgbnoaz4e0mMPE1U9KJzItSzAgB4CKYaE68ZIDF9kBGAqcRYMe9Z9cUgCjzGz+zJP1pI0Caa1aFbD8xIaN4BD6VaOsDO+3Gcc1f5POg/gfn7mMk9MG3iMYFW41nkoWTRxXb125x85iRoX4s05u0AABC9qEMEPbsTHzMGgZCZQiW8ESCvPD3rHRpYB5KKP9nINP0c0FW4vR4R4UMaiCLfFwOgf73hgVel8n4IkkMifCwyySRfkrCWdAJbSUj0ojYRvrxU0rNxZNNPVweQcqUsYAy6gGfGL7Xtm10p7UlF6fUnIUDraY6KfM/A16T9ZzZWKgC8QhVsR6gGkg4AQJ4DKWtlzA7UXKnZ2lkBMM5MJbjz3HJi8LFqrsSu3U+B2QPQ9wJviQXLDQLatgxl+c6wr5aCpPXEGcD8ZpolNbtvs7Jv3QMmvNVeZJfooq+YcaGWq0GWVT0bgLs+EMxh8nLnBsJ2dXCW0IPh9GJ+RMCuW6rzVi7Sufdi38RP4w6iQCj9d4EbIPrsCibLobUohjuA6aq9D2J97bSe27NJ5E0+ew/cRff60S4D7ft/Oy7s55GN3SqhfyjBe7iJCVYNggrY/JRyY1wHNtOGuWVzWbEhJl5C07GpD6x2EVNHH1IAOqrid7KOyiQg6T2gpw7JDQAtlwuNZtwkvlztFpm42xTS+xAQMJvpU2em+uxVU7NMwIiZRiQwktILs6wFNpLJfw64RHlGm3QJPWPwIlyXAIpgMY+PsgJ52CvCtqxAESzm8VFWIA97Rdj+A7EwUUCG28orAAAAAElFTkSuQmCC">
                </a>
            </div>
        </div>

        <div class="share-button share-button-gplus" data-share-url="{{.URL}}">
            <div class="box">
                <a href="https://plus.google.com/share?url={{.URL}}">
                    <!-- Google Plus icon by Icons8 -->
                    <img class="icon icons8-Google-Plus" width="48" height="48" src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAAEU0lEQVRoQ+1ZQWwbRRR93ykpqdfBRy42ybVBAm7l1NJbEYgNogo30vaMcFuOhCQKVxpXPbdNb0RFeCsh9VbKid5IJNJrgt0Dgovr2DEJjT+asddx1zu737vrhKLOcffPn/f+vD/zZ4bwgjd6wfHjJYGjnsHEZqBpvz6xj9GPGHyGgCxAWRDe1gQZawBXGagS6OEI9u6NOX9sJUE+FgG2s9kdjH/BgN0FK0XFWCPAOYHadXKqVWk3r10kAi7wFlAgUtGO3phRTQHFqEQGJlC3c7MMWo4L3Eu5TYQvpJ2KM0g4BiLQsPPzTFgYZIBBbYmxkHbKi9J+IgJKMg1klkE0K3Ucy455JY3ty5LcEBGo27nbhwbeZc68YjmVC2GBCCWwbecXiDAf5mgY/5mxmHHKgZINJKASFkS3hwFO6pOYp4MS20hA6b6O8c2kVxsp8AMloWqhNmnKByOBo5SOzxJrlJIvAWn0R6ZO4dh7nyA1eRKpySk9bmtzA63Nx3j20/fY33g0aMB97dUeYZoFXwJh0U9NnMToxa8x8ua7gQD3f/sFuze+BP/1JDYRU0L7Eqjb+V9NtY0CP7b0HZB+TQSqtfEIzbkZkW2gEWPNcsrvhNZCuqqk0U0/Z37g//nxFp49uIvW1mPdRdm88uElLS31rfnVDLBTi08AwAjvTXqr2L4ZqNv5AgjLfiOOLa0iNXWq/WunpsG5wL32ikjrzydi8OkfftcuGh+/YSbLuGw55WKvQR+B7encQwKd9npRCfvq0mr3c/PKOSP4KOGWEGDwvUypYkcicPzzb7UsVFOy2bslrrdEfIQEfs6UKmcCCZgSeOzafa1v1UzRd0GYEKsltnn1/e7vMPs+Ofkkcn8OTOfZD0DvYCadhgHyajzM3m8cq1R+DvOREugNlERCyl5AILcG0FveWZBIyG/mepNfL6tXzvkqTEaA161SpX1R0GniVShqEvf2U+XF7o2rkQkwODyJk1xG9cZ37X4X8N9zM7HqI9EyKt7IGk/RnPs0cCPrLTkSKSkkG1loKfHNKnBivBtVVUaofeG5UuKDizh29vyBVEJ2bdFGIS0llLP6tH8iq39aFh4SgQASAg/0J7AaN3I5ffzS/EFdZGCgZLN7czGRkmOgclofaGh8i4DAmlkfaM6ebx9oOru0kpI+0Dy4Gythe2PCwFOLaxN+x8r/55FSsZfOgjQBo9oFRd+YA+5gDTtnM1Ep6uBJ9It8reIOHnY+TgKkyUfsiy3Xcd3OrYDos2GC7fPNfMdyKqF3saFXi24+NJApHhoJ5jtpbBcSu9w9TDlJZBNYjYbJRD9wEBXD9ogwP97/arVJMc8O9YHDHbRzc6duLwpxiSjgYBQt1IoSyXiJi3LAuEq0L4ALRGz7HYKCZ4HXmcmJCtz1HYtAL8DOM6vNxOrWINt5au2c7HhdPbFCPbOyfmZ1/hPPrIPqfBj2ic3AMMBJfL4kIInSMG3+BXh1IU/9tA3qAAAAAElFTkSuQmCC">
                </a>
            </div>
        </div>
    </div>
    {{end}}
```

IsPost is a boolean and URL the blog's url.

## 5. Future functionality: Upload

I also would like to automatically upload the generated blog to github.

For that end I made the **Endpoint** interface:

```go
type Endpoint interface {
	Upload(to string) error
}
```

This inteface must implement the Upload function.
