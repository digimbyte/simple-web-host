# Simple Web Host

A zero-config static file server in a single binary. Drop it next to your `index.html`, double-click, done.

- Opens your default browser automatically
- Closes itself when all browser tabs are closed
- No install, no dependencies, no configuration

## Downloads

| Platform | Architecture | Binary |
|----------|-------------|--------|
| Windows | x64 (amd64) | [simple-web-host-windows-amd64.exe](builds/simple-web-host-windows-amd64.exe) |
| Windows | ARM64 | [simple-web-host-windows-arm64.exe](builds/simple-web-host-windows-arm64.exe) |
| Linux | x64 (amd64) | [simple-web-host-linux-amd64](builds/simple-web-host-linux-amd64) |
| Linux | ARM64 | [simple-web-host-linux-arm64](builds/simple-web-host-linux-arm64) |

## Usage

1. Place the binary in the folder containing your `index.html`
2. Run it (double-click on Windows, `./simple-web-host-linux-amd64` on Linux)
3. Your browser opens to `http://localhost:8080`
4. Close the tab — the server shuts down automatically

If port 8080 is busy, it tries the next available port up to 8099.

## How it works

The server injects a tiny heartbeat script into any `.html` response. The script pings the server every 3 seconds. When no pings arrive for 10 seconds (all tabs closed), the process exits.

## Building from source

Requires Go 1.21+. Run the included PowerShell script:

```powershell
.\build.ps1
```

Outputs all 4 binaries to `builds/`.
