# yap

**yap** (Yet Another Package-frontend) is a simple, human-readable wrapper for [yay](https://github.com/Jguer/yay), the popular AUR helper for Arch Linux.

It provides a user-friendly verb-based syntax instead of `pacman`/`yay`'s traditional flag-based system.

## Features

- **Intuitive Syntax**: use `yap install` instead of `yay -S`.
- **Transparent**: All commands are proxied to `yay`, preserving interactive prompts, colors, and functionality.
- **Safe**: No complex internal logic—just simple argument mapping.

## Usage

| Command | Equivalent `yay` Command | Description |
| :--- | :--- | :--- |
| `yap update` | `yay -Syu` | Update all packages (Repo + AUR) |
| `yap install <pkg>` | `yay -S <pkg>` | Install packages |
| `yap remove <pkg>` | `yay -Rs <pkg>` | Remove packages and unneeded dependencies (modified config files are kept as `.pacsave`) |
| `yap purge <pkg>` | `yay -Rns <pkg>` | Like `remove`, but without keeping `.pacsave` backups (files in `~` are never touched) |
| `yap search <query>` | `yay -Ss <query>` | Search for packages |
| `yap info <pkg>` | `yay -Si <pkg>` | Show package info from the repositories/AUR |
| `yap info --local <pkg>` | `yay -Qi <pkg>` | Show info for an installed package |
| `yap clean` | `yay -Sc` | Clean package cache |
| `yap autoremove` | `yay -Yc` | Remove unneeded dependencies |
| `yap list` / `yap list all` | `yay -Q` | List all installed packages |
| `yap list explicit`| `yay -Qe` | List explicitly installed packages |
| `yap list native` | `yay -Qn` | List native repository packages |
| `yap list aur` | `yay -Qm` | List AUR packages |
| `yap version` | | Show the yap version |

Extra options such as `--needed` or `--noconfirm` are passed straight through to `yay`. Flags that would switch `yay` to a different operation (e.g. `yap info -R pkg`) are rejected, so each command only does what its name says. Arguments after `--` are passed through untouched.

### Dry run

Add `--dry-run` to any command to print the `yay` command without running it:

```console
$ yap install firefox --dry-run
yay -S firefox
```

### Shell completion

```bash
# fish
yap completion fish > ~/.config/fish/completions/yap.fish
# bash
yap completion bash > ~/.local/share/bash-completion/completions/yap
# zsh (any directory in your $fpath)
yap completion zsh > "${fpath[1]}/_yap"
```

## Installation

### AUR (Arch Linux)
`yap` is available on the AUR as `yap-bin`.

```bash
yay -S yap-bin
```

> **Note:** `yap-bin` conflicts with the unrelated [`yap`](https://aur.archlinux.org/packages/yap) package (YAP Prolog), since both install `/usr/bin/yap`.

### Manual Build

1.  Clone the repository:
    ```bash
    git clone https://github.com/michaelmassoni/yap.git
    cd yap
    ```

2.  Build and install:
    ```bash
    make build
    sudo install -Dm755 yap /usr/local/bin/yap
    ```

### Verifying release binaries

Release binaries are built by GitHub Actions from the tagged source. Each release includes a `SHA256SUMS` file and a signed build provenance attestation:

```bash
sha256sum -c --ignore-missing SHA256SUMS
gh attestation verify yap_x86_64 --repo michaelmassoni/yap
```

## Development

```bash
make test        # go vet + unit tests
make build       # build ./yap for this machine
make build-all   # cross-compile release binaries into bin/
```

### Releasing

Push a `v*` tag and the [release workflow](.github/workflows/release.yml) tests, builds all architectures, generates checksums and provenance, and publishes the GitHub release:

```bash
git tag v1.1.0 && git push origin v1.1.0
```

Then update the [`yap-bin`](https://aur.archlinux.org/packages/yap-bin) PKGBUILD with the new version and the checksums from `SHA256SUMS`.
