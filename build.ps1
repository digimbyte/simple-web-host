$out = "$PSScriptRoot\builds"
$pkg = $PSScriptRoot
$env:CGO_ENABLED = "0"

Write-Host "=== Simple Web Host Build ==="
Write-Host "Package: $pkg"
Write-Host "Output:  $out"
Write-Host ""

Write-Host "Creating output directory..."
New-Item -ItemType Directory -Force -Path $out | Out-Null

# Generate Windows resource files (icon + manifest)
Write-Host ""
Write-Host "Generating Windows resources..." -ForegroundColor Cyan
$winresCheck = Get-Command go-winres -ErrorAction SilentlyContinue
if (-not $winresCheck) {
    Write-Host "  Installing go-winres..." -ForegroundColor Yellow
    go install github.com/tc-hib/go-winres@latest
}
Push-Location $pkg
go-winres make
Pop-Location
if ($LASTEXITCODE -eq 0) {
    Write-Host "  OK  resources generated" -ForegroundColor Green
} else {
    Write-Host "  WARN  go-winres failed, Windows builds will have no icon" -ForegroundColor Yellow
}

$targets = @(
    @{ GOOS = "windows"; GOARCH = "amd64"; File = "simple-web-host-windows-amd64.exe" },
    @{ GOOS = "windows"; GOARCH = "arm64"; File = "simple-web-host-windows-arm64.exe" },
    @{ GOOS = "linux";   GOARCH = "amd64"; File = "simple-web-host-linux-amd64" },
    @{ GOOS = "linux";   GOARCH = "arm64"; File = "simple-web-host-linux-arm64" }
)

$step = 0
$total = $targets.Count

foreach ($t in $targets) {
    $step++
    $label = "$($t.GOOS)/$($t.GOARCH)"
    $outFile = "$out\$($t.File)"

    Write-Host ""
    Write-Host "[$step/$total] Building $label..." -ForegroundColor Cyan

    $env:GOOS = $t.GOOS
    $env:GOARCH = $t.GOARCH
    go build -ldflags "-s -w" -o "$outFile" "$pkg"

    if ($LASTEXITCODE -eq 0) {
        $size = (Get-Item $outFile).Length / 1MB
        Write-Host "  OK  $($t.File) ($([math]::Round($size,2)) MB)" -ForegroundColor Green
    } else {
        Write-Host "  FAIL  $label failed (exit $LASTEXITCODE)" -ForegroundColor Red
    }
}

Remove-Item Env:\GOOS, Env:\GOARCH, Env:\CGO_ENABLED

Write-Host ""
Write-Host "=== Done ==="
Write-Host ""
Read-Host "Press Enter to close"
