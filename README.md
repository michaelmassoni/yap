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
| `yap remove <pkg>` | `yay -Rs <pkg>` | Remove packages (keep config) |
| `yap purge <pkg>` | `yay -Rns <pkg>` | Remove packages and config |
| `yap search <query>` | `yay -Ss <query>` | Search for packages |
| `yap info <pkg>` | `yay -Si <pkg>` | Show package info |
| `yap clean` | `yay -Sc` | Clean package cache |
| `yap autoremove` | `yay -Yc` | Remove unneeded dependencies |
| `yap list` | `yay -Q` | List all installed packages |
| `yap list explicit`| `yay -Qe` | List explicitly installed packages |
| `yap list native` | `yay -Qn` | List native repository packages |
| `yap list aur` | `yay -Qm` | List AUR packages |

## Installation

### AUR (Arch Linux)
`yap` is available on the AUR as `yap-bin`.

```bash
yay -S yap-bin
```

### Manual Build

1.  Clone the repository:
    ```bash
    git clone https://github.com/michaelmassoni/yap.git
    cd yap
    ```

2.  Build and install:
    ```bash
    go build -o yap main.go
    sudo mv yap /usr/local/bin/
    ```
