# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go desktop GUI application for flashing firmware onto SAVART-EV electric vehicle components: **VCU, BMS, HMI, and Keyless** modules. It uses **Fyne v2** for the GUI and communicates with both a local PostgreSQL database and remote REST APIs for serial number generation, firmware versioning, and authentication.

## Build & Run

```powershell
# Run directly
go run main.go

# Build executable
go build -o produksi.exe .

# Sync dependencies
go mod tidy
```

No Makefile or test suite is present. There is no `go test` coverage to run.

## Architecture & Data Flow

### Entry Point

`main.go` → `modul.Loginapp()` is the sole entry point. Login checks `login.txt` for a stored token before showing the login UI.

### Module Responsibilities (`modul/`)

| File | Responsibility |
|---|---|
| `gui.go` | Fyne UI: login screen, main form with board-type-specific fields, dialogs |
| `app.go` | Orchestration: `ProsesFlash()` → `Ceking_firmware()` + `Programflash()` |
| `flash.go` | OpenOCD invocation per board type; MCU serial number reading |
| `database.go` | PostgreSQL connection + REST API calls to `part.savart-ev.com` and `team.savart-ev.com` |
| `download.go` | Firmware download, CRC verification, version file management |
| `sn.go` | Board-specific serial number generation via REST API + binary encoding |
| `tooldir.go` | File/folder utilities, RAR extraction, login file parsing, version file I/O |
| `exsel.go` | Excel production log creation/appending |
| `crc.go` | CRC32 firmware verification |
| `time.go` | NTP time fetch from `pool.ntp.org` |

### Core Data Structure

The `Bus` struct (defined in `download.go`) is the central data carrier passed through the flash pipeline:

```go
type Bus struct {
    Bord          string  // "vcu", "bms", "hmi", "keyless"
    Id            uint32  // Hardware model ID
    Model         uint32
    Modelver      string
    Versifirmapp  uint32  // Application firmware version
    Versifirmboot uint32  // Bootloader firmware version
    Num           SN      // Serial number fields
}
```

### Processing Pipeline

1. **Login**: `IpaLogin()` → POST `https://team.savart-ev.com/api/login` (Basic Auth) → token saved to `login.txt`
2. **Firmware check**: `Ceking_firmware()` → `CekVersion()` → GET `https://part.savart-ev.com/api/[board]/versi` → download if version mismatch → verify CRC
3. **Flash**: `Programflash()` routes to board-specific `ExecuteOpenOCD[Board]()` in `flash.go`
   - Reads existing MCU serial number via JTAG (`ReadSNH7()` for STM32H7, `ReadSN()` for others)
   - Flashes bootloader + application using `.\openocd\bin\openocd.exe` with STLink-v2
4. **Serial Number**: `Sn[Board]()` in `sn.go` → POST to `https://part.savart-ev.com/api/[board]/generate-sn` → binary-encodes SN fields → programs into MCU
5. **Database update**: `Update[Board]()` in `database.go` → REST API + local PostgreSQL insert

### Configuration Sources

- `battrymenu.xlsx` — BMS model/voltage/parallel/cell options (read at runtime by `gui.go`)
- `menu.xlsx` — Dropdown options for type, LCD, VCU/HMI/Keyless MCU variants
- `./bin/[board]/[model]/versi.txt` — Local firmware version cache
- `login.txt` — Stored auth token (format: `{token}`)

### External Runtime Dependencies

- **OpenOCD** at `.\openocd\bin\openocd.exe` — must be present; called via `exec.Command`
- **PostgreSQL** at `localhost:5432`, database `server_produksi`, user `postgres`
- **Internet access** to `part.savart-ev.com` and `team.savart-ev.com`

### Firmware Storage Layout

```
./bin/[board]/[model]/
    bootloader.bin
    application.bin
    versi.txt
```

Downloaded as RAR archives, extracted by `ExtractRar()` / `Ekstrak_move()` in `tooldir.go`.

## Key Conventions

- Board type strings are lowercase: `"vcu"`, `"bms"`, `"hmi"`, `"keyless"`.
- Indonesian is used in variable names, comments, and some UI strings alongside English — this is intentional.
- Production logs are written to `./data_produksi/[board]/` as Excel files via `exsel.go`.
- STM32H7 boards use `ReadSNH7()` with dual-bank flash handling; other boards use `ReadSN()`.
