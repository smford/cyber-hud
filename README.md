# CYBER-HUD // Neural Mainframe Terminal HUD & Cyberdeck

[![Go Version](https://img.shields.io/badge/go-1.24%2B-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![Build Status](https://img.shields.io/badge/build-passing-05FFA1?style=flat-square&logo=githubactions)](https://github.com/smford/cyber-hud)
[![Architecture](https://img.shields.io/badge/arch-Bubble%20Tea%20%2F%20Lip%20Gloss-FCEE0A?style=flat-square)](https://github.com/charmbracelet/bubbletea)
[![Binary](https://img.shields.io/badge/runtime-zero--dependency%20%7C%20static-00F0FF?style=flat-square)](#sre-grade-architecture--design-principles)
[![License](https://img.shields.io/badge/license-MIT-FF003C?style=flat-square)](LICENSE)

A high-performance, SRE-grade terminal Head-Up Display (HUD) and netrunner cyberdeck environment written in Go. Inspired by [`smford.github.io`](https://stephenford.org), `cyber-hud` renders an interactive multi-pane cyberpunk workstation directly within your terminal's alternate screen buffer.

Built with Charm's **Bubble Tea**, **Lip Gloss**, and **Bubbles**, it delivers hardware telemetry, an interactive netrunner command shell, hardware engineering build dossiers with ASCII schematics, a 27-node firmware repository deck, skill telemetry, and an authentic VT100 Matrix digital rain screensaver.

---

## Table of Contents

- [Architectural & SRE Principles](#sre-grade-architecture--design-principles)
- [Key Features](#key-features)
- [HUD Workspace Views](#hud-workspace-views)
- [Quick Start](#quick-start)
- [CLI Reference & Flags](#cli-reference--flags)
- [Controls & Keybinding Matrix](#controls--keybinding-matrix)
- [Cyberdeck Virtual Shell (CLI)](#cyberdeck-virtual-shell-cli)
- [Configuration Specification](#configuration-specification-configyaml--configjson)
- [Operations & Deployment Runbook](#operations--deployment-runbook)
- [Testing & Quality Assurance](#testing--quality-assurance)
- [License & Credits](#license--credits)

---

## SRE-Grade Architecture & Design Principles

`cyber-hud` was designed from the ground up under senior developer and site reliability engineering (SRE) paradigms to ensure zero friction in critical terminal workflows:

```
                      ┌────────────────────────────────────────┐
                      │             OS / Terminal              │
                      │  (SIGWINCH, SIGINT, Mouse, Keystrokes) │
                      └───────────────────┬────────────────────┘
                                          │
                                          ▼
                      ┌────────────────────────────────────────┐
                      │          Bubble Tea Event Loop         │
                      │        (Elm Architecture / MVT)        │
                      └───────────────────┬────────────────────┘
                                          │
                  ┌───────────────────────┼───────────────────────┐
                  ▼                       ▼                       ▼
        ┌───────────────────┐   ┌───────────────────┐   ┌───────────────────┐
        │  Model State &    │   │  Virtual Terminal │   │  Overlay Engines  │
        │  Multi-Tab Router │   │  Engine & VFS     │   │  (Rain / Dossier) │
        └─────────┬─────────┘   └─────────┬─────────┘   └─────────┬─────────┘
                  │                       │                       │
                  └───────────────────────┼───────────────────────┘
                                          │
                                          ▼
                      ┌────────────────────────────────────────┐
                      │         Lip Gloss Style Pipeline       │
                      │     (Dynamic Geometry Clamping)        │
                      └───────────────────┬────────────────────┘
                                          │
                                          ▼
                      ┌────────────────────────────────────────┐
                      │   Terminal Alternate Screen Buffer     │
                      │  (Zero Scrollback State Pollution)     │
                      └────────────────────────────────────────┘
```

1. **Deterministic State Transitions (Elm Architecture / MVT)**:
   - State mutations follow a strict Model-View-Update cycle. Sub-models (`Overview`, `BuildLogs`, `FirmwareDeck`, `OperatorStack`, `DossierModal`, `MatrixRain`) are decoupled and isolated.
2. **Zero Terminal Pollution & Clean Lifecycle Recovery**:
   - Runs exclusively within the terminal alternate screen buffer (`tea.WithAltScreen`).
   - Intercepts signals (`SIGINT`, `SIGTERM`, `SIGWINCH`). On shutdown or crash, terminal raw flags, cursor visibility, and previous scrollback history are cleanly restored.
3. **Strict Dynamic Viewport Reflow**:
   - Subscribes to `tea.WindowSizeMsg`, dynamically recalculating cards, split viewports, and grids.
   - Enforces strict line-count clamping against terminal height (`m.height`), preventing viewport drift or unintended terminal scrolling.
4. **Dual Modality & Operator Parity**:
   - 100% of workflows are executable via keyboard.
   - Full mouse tracking (`tea.WithMouseCellMotion`) allows point-and-click tab navigation, quick chip triggering, and wheel scrolling.
   - Instant mouse toggle (`F2` or `Ctrl+M`) disables mouse reporting on the fly, freeing native mouse selection for immediate clipboard copy-paste in macOS/X11/Wayland.
5. **Zero External Runtime Dependencies**:
   - Compiles to a self-contained, statically linkable binary (`CGO_ENABLED=0`).
   - Embedded fallback configuration ensures continuous operation even if external files are unreadable or missing.

---

## Key Features

- **Multi-Tab Workstation Layout**:
  - `[01] OVERVIEW_HUD`: Telemetry metrics, real-time command terminal, quick-launch chips, and live dispatches.
  - `[02] BUILD_LOGS`: Hardware and electronics build logs with category filters and interactive modal dossiers.
  - `[03] FIRMWARE_DECK`: Interactive registry of 27 repositories with GitHub star/fork metrics, tech badges, and browser launching.
  - `[04] OPERATOR_STACK`: Skill bars with dynamic percentage rendering and workshop bench inventories.
- **Embedded ASCII Schematic & BOM Dossiers**:
  - Full-screen modal inspector displaying circuit ASCII schematics, Bill of Materials (BOM), and GPIO pinout assignments.
- **Matrix Digital Rain Screensaver**:
  - Authentic VT100 Katakana digital rain stream with frame throttling; exit anytime via `Esc` or `Q`.
- **Dual Cyberware Color Themes**:
  - `CYBERPUNK 2077`: Night City Yellow (`#FCEE0A`), Neon Cyan (`#00F0FF`), Arasaka Red (`#FF003C`).
  - `ALTERED CARBON`: Psychasec Ice Blue (`#00E5FF`), Needlecast Violet (`#BD00FF`), Emerald Green (`#05FFA1`).
  - Toggle anytime via `F1`, `T`, or clicking the header theme badge.
- **Virtual File System (VFS) & Netrunner CLI**:
  - Embedded virtual shell supporting `ls`, `cat`, `mcat`, `which`, `projects`, `cidr`, and custom YAML commands.

---

## HUD Workspace Views

### 1. `[01] OVERVIEW_HUD`
Main telemetry and command bridge. Displays kernel, cyberware status, and active security ICE levels alongside an interactive terminal emulator with persistent command history.

### 2. `[02] BUILD_LOGS`
Hardware engineering build logs with category filtering (`ALL LOGS`, `Hardware`, `Optics`, `Fabrication`). Press `Enter` or click any log to summon the schematic dossier modal.

### 3. `[03] FIRMWARE_DECK`
Filterable portfolio containing 27 project repositories across categories (`Embedded`, `RaspberryPi`, `Networking`, `3DPrinting`, `Terminal`, `Web`, `Fabrication`). Press `Enter` or click to open repository and demo URLs directly in your default browser.

### 4. `[04] OPERATOR_STACK`
Real-time capability meters (C++ / ESP32, Go / Systems, Linux / SRE, CAD / Fusion360, Timber Fabrication) and physical workshop inventories (Laser Lab, 3D Print Farm, Timber Mill, Electronics Bench).

---

## Quick Start

### Prerequisites
- **Go 1.24+** (Go 1.21+ compatible)
- A modern terminal emulator with 24-bit TrueColor support (e.g., Ghostty, iTerm2, WezTerm, Alacritty, Kitty).
- Optional: Nerd Font installed for full glyph fidelity.

### Build from Source

```bash
# Clone the repository
git clone https://github.com/smford/cyber-hud.git
cd cyber-hud

# Build optimized binary
go build -ldflags="-s -w" -o cyber-hud ./cmd/cyber-hud

# Run cyber-hud
./cyber-hud
```

---

## CLI Reference & Flags

`cyber-hud` provides command-line flags for automation, headless environments, and configuration management:

```
Usage: ./cyber-hud [OPTIONS]

Options:
  -c string
        Path to YAML or JSON config file (default: embedded/local config.yaml)
  -F
        Request full-screen terminal window mode (macOS native automation)
  --skip-boot
        Skip BIOS POST boot animation and jump directly to HUD
  --dump-config string
        Export default configuration to specified YAML or JSON file and exit
  -v
        Display cyber-hud version, target OS, and architecture
```

### Examples

```bash
# Skip boot animation directly into HUD
./cyber-hud --skip-boot

# Launch with macOS native fullscreen space
./cyber-hud -F

# Load a dedicated incident response configuration
./cyber-hud -c /etc/cyber-hud/incident-response.yaml

# Export default configuration to inspect or modify
./cyber-hud --dump-config config.yaml
```

---

## Controls & Keybinding Matrix

`cyber-hud` delivers complete operator parity: every interaction can be triggered via keyboard, mouse, or virtual terminal CLI.

### Global Navigation

| Key / Action | Scope | Description |
|---|---|---|
| `1`, `2`, `3`, `4` | Global | Switch to Tab 1 (`OVERVIEW`), 2 (`LOGS`), 3 (`FIRMWARE`), or 4 (`OPERATOR`) |
| `Tab` / `Shift+Tab` | Global | Cycle sequentially between tabs |
| `F1` or `T` | Global | Toggle theme (`CYBERPUNK 2077` ⟷ `ALTERED CARBON`) |
| `F2` or `Ctrl+M` | Global | **Toggle Mouse Mode**: Instantly toggle mouse cell tracking on/off |
| `Esc` or `Q` | Modals / Screensaver | Dismiss schematic dossier or exit Matrix digital rain |
| `Ctrl+C` | Global | Graceful emergency exit and terminal restoration |

### List & Viewport Navigation

| Key / Action | Scope | Description |
|---|---|---|
| `j` / `k` or `↑` / `↓` | Lists & Viewports | Scroll items or scroll terminal history |
| `h` / `l` or `←` / `→` | Category Bars | Switch active category filter in Logs and Firmware tabs |
| `PgUp` / `PgDn` | Viewports | Scroll half-page up or down |
| `Enter` / `Space` | Lists | Inspect schematic dossier or launch repository in browser |

### Mouse Interactions (When Mouse Mode is Enabled)

- **Click Tab Bar**: Instantly switch HUD workspace.
- **Click Theme Chip**: Toggles color scheme.
- **Click Quick Chips**: Executes terminal command chips (`matrix-rain`, `cidr`, `projects`, etc.).
- **Mouse Wheel**: Smooth vertical viewport and list scrolling.

> **SRE Tip**: Need to select and copy text from the terminal? Hit **`F2`** or **`Ctrl+M`** to release mouse capture. Select your text natively, then press **`F2`** to re-engage mouse mode.

---

## Cyberdeck Virtual Shell (CLI)

The embedded terminal on `[01] OVERVIEW_HUD` features a custom netrunner shell:

| Command | Arguments | Description |
|---|---|---|
| `help` | — | Display available commands and quick-start tips |
| `projects` | — | List all 27 repository nodes with star counts and URLs |
| `matrix-rain` | — | Launch full-terminal VT100 Katakana digital rain screensaver |
| `cidr` | `<ip/mask>` or `<ip>+<count>` | Decompose subnet, network, broadcast, wildcard mask, and host ranges |
| `ls` | — | List files in the embedded Virtual File System (VFS) |
| `cat` | `<filename>` | Display file content (e.g. `cat notes.md`, `cat hardware.ino`) |
| `mcat` | `<filename>` | Stream file content with matrix visual headers |
| `open` | `<repo_or_alias>` | Open target repository or custom command in the host default browser |
| `which` | `<binary>` | Query the simulated cyberdeck PATH registry |
| `delim` | — | Render zero-allocation terminal visual barrier |
| `clear` | — | Flush terminal scrollback buffer |
| `reboot` / `boot` | — | Re-run hardware BIOS boot sequence |

---

## Configuration Specification (`config.yaml` / `config.json`)

`cyber-hud` can be entirely customized via YAML or JSON. You can export the default configuration at any time:

```bash
./cyber-hud --dump-config config.yaml
```

### Schema Overview

```yaml
operator:
  handle: "smford"
  host: "mainframe"
  title: "smford.github.io"
  subtitle: "> WOODWORKING SOFTWARE ENGINEER // HARDWARE PROTOTYPER"
  bio: "Bridging low-level firmware engineering and physical timber fabrication."
  sec_net: "ACTIVE"

telemetry:
  kernel: "6.1.0-DEBIAN"
  cyberware: "ESP32 / PI / CAD"
  security_ice: "UNRESTRICTED"
  repositories: "27 DEPLOYED"

settings:
  default_theme: "cyberpunk"   # "cyberpunk" | "altered-carbon"
  mouse_enabled: true          # Set true to enable mouse cell motion by default
  skip_boot: false             # Set true to bypass BIOS POST on startup

quick_chips:
  - label: "help"
    command: "help"
  - label: "matrix-rain"
    command: "matrix-rain"
  - label: "cidr"
    command: "cidr 192.168.1.10+5"
  - label: "projects"
    command: "projects"

custom_commands:
  blackhole:
    description: "Chandra-X black hole raymarcher"
    output: "[CHANDRA-X ACCRETION SIMULATOR] Event Horizon: 29.53 km"
    url: "https://stephenford.org/blackhole/"

repositories:
  - name: "esp32-asyncwebserver-fileupload-example"
    category: "Embedded"
    stars: 104
    forks: 30
    desc: "Asynchronous web server for ESP32 with drag-and-drop SPIFFS uploads."
    tech: ["C++", "ESP32", "SPIFFS"]
    url: "https://github.com/smford/esp32-asyncwebserver-fileupload-example"

build_logs:
  - id: "log-042"
    category: "Hardware"
    title: "EEH ESP32 RFID Tool Sentry"
    date: "2026-04-12"
    status: "DEPLOYED_ACTIVE"
    summary: "Safety interlock for 3kW workshop tools."
    bom:
      - item: "ESP32 NodeMCU"
        qty: "1"
        role: "Microcontroller"
    pinout:
      - gpio: "GPIO 25"
        device: "Relay Pin"
        note: "Active Low"
    schematic: |
      [ 230V AC Mains ] ---> [ 30A Contactor ] ---> [ Table Saw ]
                                   ^
      [ ESP32 GPIO 25 ] ---> [ Optocoupler ] ---> [ 5V Relay ]
```

---

## Operations & Deployment Runbook

### Cross-Compilation

Compile static binaries for Linux jump-boxes, Raspberry Pi field nodes, or Windows workstations:

```bash
# Linux x86_64 (Standard server / Bastion)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o cyber-hud-linux-amd64 ./cmd/cyber-hud

# Linux ARM64 (Raspberry Pi 4 / 5, Graviton)
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o cyber-hud-linux-arm64 ./cmd/cyber-hud

# macOS Universal (Apple Silicon & Intel)
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o cyber-hud-darwin-arm64 ./cmd/cyber-hud
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o cyber-hud-darwin-amd64 ./cmd/cyber-hud
lipo -create -output cyber-hud-macos cyber-hud-darwin-arm64 cyber-hud-darwin-amd64

# Windows x86_64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o cyber-hud-windows.exe ./cmd/cyber-hud
```

### SSH Bastion & Headless Execution

`cyber-hud` runs smoothly over SSH connections. For optimal rendering:
1. Ensure `TERM` supports color (`export TERM=xterm-256color`).
2. If running within `tmux`, enable mouse and TrueColor support in `~/.tmux.conf`:
   ```tmux
   set -g mouse on
   set -g default-terminal "tmux-256color"
   set -ga terminal-overrides ",xterm-256color:Tc"
   ```

---

## Testing & Quality Assurance

The codebase includes automated unit tests covering configuration roundtripping (YAML & JSON), default fallback validation, terminal CLI execution, and viewport rendering:

```bash
# Run all unit tests
go test -v ./...

# Run tests with race condition detector
go test -race ./...

# Run static analysis
go vet ./...
```

---

## License & Credits

- **License**: MIT License. See [LICENSE](LICENSE) for details.
- **Author**: Stephen Ford ([@smford](https://github.com/smford)) // [stephenford.org](https://stephenford.org).
- **Core TUI Framework**: Powered by [Charm](https://charm.sh) ([`bubbletea`](https://github.com/charmbracelet/bubbletea), [`lipgloss`](https://github.com/charmbracelet/lipgloss), [`bubbles`](https://github.com/charmbracelet/bubbles)).
