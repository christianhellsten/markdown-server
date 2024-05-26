# Markdown Server: browse and render Markdown files

Markdown Server is a simple web server written in Go that serves `.md` files as HTML and
lists Markdown files and directories in a navigable menu. It uses the
`blackfriday` library to convert Markdown files to HTML and `Pico.css` for
basic styling.

## Features

- Displays a menu listing all Markdown files in the current directory
- Converts Markdown files to HTML using the `goldmark` library
- `Pico.css` for modern, minimal styling
- Syntax highlighting for code blocks using `highlight.js`
- HTML templates
- Image support
- Ignore files similar to .gitignore

## Installation

Binaries for MacOS, Windows, and Linux can be found on Github:
https://github.com/christianhellsten/markdown-server/releases

You can also use brew to install the binaries from Github with Homebrew:

```bash
brew tap christianhellsten/markdown-server https://github.com/christianhellsten/markdown-server.git
brew install christianhellsten/markdown-server/markdown-server
```

To uninstall, use:

```bash
brew uninstall christianhellsten/markdown-server/markdown-server
brew untap christianhellsten/markdown-server
```

## Running

To start the server, run the **markdown-server** command:

```bash
markdown-server --help
Usage of markdown-server:
  -dir string
    	Base directory to serve files from (default ".")
  -host string
    	Host to listen on (default "localhost")
  -port int
    	Port to listen on (default 8080)
```

## Ignoring files

Run the following commands to ignore, for example, Markdown files in the node_modules directory:

```bash
mkdir .markdown-server
echo 'node_modules' >> .markdown-server/ignore
```

## Theme

To customize the default theme, run the following commands:


```bash
mkdir .markdown-server
touch .markdown-server/index.html
touch .markdown-server/menu.html
```

Edit the files to create your own theme, for example:

**.markdown-server/index.html**:

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{ .Title }}</title>
    <link
      rel="stylesheet"
      href="https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.slate.min.css"
    />
    <link
      rel="stylesheet"
      href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/default.min.css"
    />
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
    <script>
      hljs.highlightAll();
    </script>
  </head>
  <body>
    <main class="container">
      <header>{{ .Menu }}</header>
      <article>{{ .Content | safeHTML }}</article>
    </main>
  </body>
</html>
```

**.markdown-server/menu.html**:

```html
<details class="dropdown">
  <summary role="button" class="contrast">📁 {{ .UrlPath }}</summary>
  {{ .Menu | safeHTML }}
</details>
```


## Contributing

Clone the repository:

```sh
git clone https://github.com/christianhellsten/markdown-server.git
cd markdown-server
./run.sh
# or, use ./watch.sh
# or, use go run main.go
```

## Screenshot

![Screenshot](markdown-server-screenshot.png)
