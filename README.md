# gog

go to GitHub in go

## Installation

- Clone the repository
- Run `go build -o gog`
- Run `sudo mv gog /usr/local/bin/gog`

## Usage

```sh
GitHub CLI tool for quickly opening repositories and profiles

Usage:
  gog [shortcut or repo] [flags]
  gog [command]

Available Commands:
  add-shortcut Add a shortcut for a repository
  completion   Generate the autocompletion script for the specified shell
  get-config   Gets the contents of the config file
  help         Help about any command
  set-username Set your GitHub username

Flags:
  -h, --help   help for gog

Use "gog [command] --help" for more information about a command.
```

```sh
go build -o gog
sudo mv gog /usr/local/bin
gog
```

```sh
gog -u {new_username}
```
