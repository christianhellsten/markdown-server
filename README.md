# Markdown Server: browse and render Markdown files

Markdown Server is a simple web server written in Go that serves `.md` files as HTML and
lists Markdown files and directories in a navigable menu. It uses the
`blackfriday` library to convert Markdown files to HTML and `Pico.css` for
basic styling.

## Features

- Serve Markdown files as HTML
- Display a navigable menu of directories and Markdown files
- Use `Pico.css` for modern, minimal styling
- Syntax highlighting for code blocks using `highlight.js`
- Custom HTML templates
- Image handling

## Installation with Homebrew

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

## Manual Installation

Binaries for MacOS, Windows, and Linux can be found on Github:
https://github.com/christianhellsten/markdown-server/releases

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
