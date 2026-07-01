# mira

**mira** is a tool installer that reads `mira.json` manifests from GitHub repositories. It downloads the correct binary for your platform, verifies its SHA256 checksum, and installs it to `~/.local/bin/`.

## Install

```bash
go install github.com/coderianx/mira/cmd/mira@latest
```

Or download a prebuilt binary from the [releases page](https://github.com/coderianx/mira/releases).

Make sure `~/.local/bin` is in your `PATH`:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

## Usage

```bash
# Install from the main branch
mira install github.com/user/reponame

# Install a specific release
mira install github.com/user/reponame@v1.0.0

# List installed tools
mira list

# Show details about an installed tool
mira info github.com/user/reponame

# Reinstall/update a tool
mira update github.com/user/reponame

# Uninstall a tool
mira uninstall github.com/user/reponame
```

## Manifest format

The target repository must contain a `mira.json` at its root:

```json
{
  "name": "my-tool",
  "version": "v1.0.0",
  "description": "A useful CLI tool",
  "author": "you",
  "repo": "github.com/you/my-tool",
  "bin": "my-tool",
  "platforms": {
    "linux/amd64": {
      "url": "https://raw.githubusercontent.com/you/my-tool/main/my-tool-linux-amd64",
      "sha256": "abc123..."
    },
    "darwin/arm64": {
      "url": "https://raw.githubusercontent.com/you/my-tool/main/my-tool-darwin-arm64",
      "sha256": "def456..."
    }
  }
}
```

The binary can be stored directly in the repo (served via `raw.githubusercontent.com`) or hosted elsewhere — the `url` field is fetched directly.

## How it works

1. Fetches `mira.json` from the repository (`main` branch or specified tag)
2. Matches your current OS/arch against the `platforms` map
3. Downloads the binary from the matching URL
4. Verifies the SHA256 checksum
5. Installs it to `~/.local/bin/<bin>` and records it in `~/.local/share/mira/state.json`

## Commands

| Command | Description |
|---|---|
| `install <repo>` | Download and install a tool |
| `update <repo>` | Re-download and reinstall |
| `uninstall <repo>` | Remove an installed tool |
| `list` | Show all installed tools |
| `info <repo>` | Show details about an installed tool |
| `version` | Print version info |

## Ideas

See [ideas.md](ideas.md) for planned features.

## License

MIT
