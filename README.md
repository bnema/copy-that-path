# yank-that

`yank-that` is a small CLI that copies useful values to your clipboard:
- resolved file paths
- file contents
- GitHub file URLs (with optional line anchors)

## Install

```bash
go install github.com/bnema/yank-that/cmd/yt@latest
```

## Commands

- `yt path [target]` (`ytp` shell alias): copy resolved absolute path
- `yt content <file>` (`ytc` shell alias): copy file content
- `yt github <file[:line]>` (`ytg` shell alias): copy GitHub blob URL

## Examples

```bash
# path
yt path
yt path .
yt path file.txt
yt path ../foo

# content
yt content file.txt
yt content ../notes.md

# github URL
yt github README.md
yt github internal/app/app.go:42
```

## Shell aliases

Install aliases for `ytp`, `ytc`, and `ytg`:

```bash
yt install-aliases
```

Install for one shell only:

```bash
yt install-aliases bash
yt install-aliases zsh
yt install-aliases fish
```

Then reload your shell config (for example, `source ~/.bashrc`, `source ~/.zshrc`, or restart fish).

## Requirements

Install at least one clipboard backend:
- `wl-copy` (Wayland)
- `xclip` (X11)
