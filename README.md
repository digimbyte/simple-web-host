# Simple Web Host

A zero-config static file server in a single binary. Drop it next to your `index.html`, double-click, done.

- Opens your default browser automatically
- Closes itself when all browser tabs are closed
- No install, no dependencies, no configuration

## Downloads

| Platform | Architecture | Binary |
|----------|-------------|--------|
| Windows | x64 (amd64) | [simple-web-host-windows-amd64.exe](https://raw.githubusercontent.com/digimbyte/simple-web-host/master/builds/simple-web-host-windows-amd64.exe) |
| Windows | ARM64 | [simple-web-host-windows-arm64.exe](https://raw.githubusercontent.com/digimbyte/simple-web-host/master/builds/simple-web-host-windows-arm64.exe) |
| Linux | x64 (amd64) | [simple-web-host-linux-amd64](https://raw.githubusercontent.com/digimbyte/simple-web-host/master/builds/simple-web-host-linux-amd64) |
| Linux | ARM64 | [simple-web-host-linux-arm64](https://raw.githubusercontent.com/digimbyte/simple-web-host/master/builds/simple-web-host-linux-arm64) |

## Usage

1. Place the binary in the folder containing your `index.html` and static files
2. Run it (double-click on Windows, `./simple-web-host-linux-amd64` on Linux)
3. Your browser opens automatically
4. Close the tab — the server shuts down automatically

The server serves **everything** in the same directory as the binary — HTML, CSS, JS, images, markdown, JSON, etc.

### Options

```
-port <number>    Listen on a specific port
-stable           Persistent mode: no auto-close, no heartbeat injection
```

Examples:

```bash
# Default: auto-selects port, closes when all tabs close
./simple-web-host-windows-amd64.exe

# Use a specific port
./simple-web-host-windows-amd64.exe -port 3000

# Persistent server — stays running until you Ctrl+C
./simple-web-host-windows-amd64.exe -stable

# Combine both
./simple-web-host-windows-amd64.exe -stable -port 3000
```

If no `-port` is given and 8080 is busy, it tries the next available port up to 8099.

**Default mode** — the server injects a heartbeat into HTML responses and shuts down when all browser tabs are closed. Great for quick local viewing.

**Stable mode** (`-stable`) — pure static file server with zero injection. HTML is served unmodified. The server runs until manually stopped. Use this for development or when serving to other devices on the network.

## How it works

- Serves all files from the directory the binary is in
- Auto-opens your default browser on startup
- Injects a tiny heartbeat script into `.html` responses — pings the server every 3 seconds
- When no pings arrive for 10 seconds (all tabs closed), the process exits
- Multiple browser tabs are supported — the server stays alive until the last one closes

## Building from source

Requires Go 1.21+. Run the included PowerShell script:

```powershell
.\build.ps1
```

Outputs all 4 binaries to `builds/`.
