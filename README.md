# bigcrock

[![Go Version](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/Ars-byte/bigcrock?color=brightgreen)](https://github.com/Ars-byte/bigcrock/releases)

> **Language / Idioma**: **English** | [Español](README.es.md)

**bigcrock** is a fast, lightweight, Ranger-inspired terminal directory tree analyzer and disk usage ranker built in Go using [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

It scans your directories and ranks them by size, showing at a glance which folders and hidden directories (like `.local`, `.cache`, etc.) consume the most space.

---

## Features

- **Ranger-Style Miller Columns**:
  - **Left column**: Parent directory context.
  - **Center column**: Active directory entries ranked by disk space usage with size badges and percentage progress bars.
  - **Right column**: Deep inspector with top subdirectory breakdown, largest files, metadata, and preview.
- **Top Space Consumers Ranking**: Instantly displays folders ranked from heaviest to lightest (`#1`, `#2`, `#3`...).
- **Hidden Directories Included**: Easily inspect hidden folders such as `.local`, `.config`, `.cache` (toggleable with `.` or `H`).
- **Interactive Tree Navigation**: Drill down into any directory or subdirectory (`Enter` / `l` / `→`) and navigate back (`←` / `h` / `Backspace`).
- **Non-Blocking Background Scanning**: Asynchronous concurrent size calculation keeps the UI smooth and responsive even across massive file hierarchies.
- **Intelligent In-Memory Cache**: Instant navigation when returning to previously scanned directories.
- **Dual Language Support**: Switch between English and Spanish seamlessly in-app using `Tab` or `t`.
- **Multiple Sort Modes**: Sort by size descending, size ascending, alphabetical, or items count (`s`).
- **About Screen**: Clean, borderless project view (`?` or `a`).
- **Modular Architecture**: Built with clean package separation modeled after the Crayoz project.

---

## Keybindings

| Key | Action |
| --- | --- |
| `↑` / `k` | Move cursor up |
| `↓` / `j` | Move cursor down |
| `Enter` / `→` / `l` | Enter selected directory |
| `←` / `h` / `Backspace` | Go up to parent directory |
| `Home` / `g` | Jump to top (#1 heaviest item) |
| `End` / `G` | Jump to bottom item |
| `PgUp` / `PgDn` | Scroll page up / page down |
| `s` | Cycle sort mode (Size ↓ / Size ↑ / Name A-Z / Name Z-A / Items Count) |
| `.` / `H` | Toggle hidden files and directories |
| `r` | Force refresh and recalculate sizes |
| `Tab` / `t` | **Toggle language (English ↔ Español)** |
| `?` / `a` | **Show About screen** |
| `Any key` | Return from About screen |
| `q` / `Ctrl+C` | Quit application |

---

## Installation & Build

### Requirements
- [Go](https://go.dev/) 1.22 or higher.

### Build from source
```bash
git clone https://github.com/Ars-byte/bigcrock.git
cd bigcrock
go mod download
go build -ldflags="-s -w" -o bigcrock .
./bigcrock
```

### Usage
```bash
# Analyze current directory
./bigcrock

# Analyze home directory
./bigcrock ~

# Analyze specific directory
./bigcrock /var/log
```

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
