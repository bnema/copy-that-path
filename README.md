# yank-that

Single CLI for yanking paths, file content, and GitHub file URLs to your clipboard.

## Install

```bash
go install github.com/bnema/yank-that/cmd/yt@latest
```

## Usage

### `yt path` (shell alias: `ytp`) - yank resolved absolute path

```bash
yt path              # copies current directory
yt path .            # copies current directory
yt path file.txt     # copies absolute path to file.txt
yt path ../foo       # copies resolved absolute path
```

### `yt content` (shell alias: `ytc`) - yank file content

```bash
yt content file.txt     # copies file content to clipboard
yt content ../notes.md  # copies content from resolved path
```

### `yt github` (shell alias: `ytg`) - yank GitHub file URL

```bash
yt github README.md        # copies https://.../blob/<ref>/README.md
yt github internal/app/app.go:42
```

### Install shell aliases

```bash
yt install-aliases         # installs aliases for bash, zsh, and fish
yt install-aliases bash    # installs only for bash
yt install-aliases zsh     # installs only for zsh
yt install-aliases fish    # installs only for fish
```

## Requirements

Requires `wl-copy` (Wayland) or `xclip` (X11) to be installed.
