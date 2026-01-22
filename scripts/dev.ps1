$ErrorActionPreference = "Stop"

$root = Resolve-Path (Join-Path $PSScriptRoot "..")
$global:serverProcess = $null
$global:restartPending = $false
$global:lastRestart = Get-Date

function Start-Server {
  Write-Host "Starting server..."
  $global:serverProcess = Start-Process -FilePath "go" -ArgumentList "run ./cmd/server" -WorkingDirectory $root -NoNewWindow -PassThru
}

function Stop-Server {
  if ($global:serverProcess -and -not $global:serverProcess.HasExited) {
    Write-Host "Stopping server..."
    $global:serverProcess.Kill()
    $global:serverProcess.WaitForExit()
  }
}

function Restart-Server {
  $global:restartPending = $false
  Stop-Server
  Start-Server
  $global:lastRestart = Get-Date
}

Start-Server

$watcher = New-Object System.IO.FileSystemWatcher $root
$watcher.IncludeSubdirectories = $true
$watcher.NotifyFilter = [IO.NotifyFilters]"LastWrite, FileName, DirectoryName"
$watcher.EnableRaisingEvents = $true

$action = {
  $path = $Event.SourceEventArgs.FullPath
  if ($path -match "\\\.git\\") { return }
  if ($path -match "\\bin\\") { return }
  if ($path -match "\\tmp\\") { return }
  if ($path -notmatch "\.(go|html|css|sql)$") { return }

  $now = Get-Date
  if (($now - $global:lastRestart).TotalMilliseconds -lt 400) { return }
  $global:restartPending = $true
}

$subs = @()
$subs += Register-ObjectEvent -InputObject $watcher -EventName Changed -Action $action
$subs += Register-ObjectEvent -InputObject $watcher -EventName Created -Action $action
$subs += Register-ObjectEvent -InputObject $watcher -EventName Deleted -Action $action
$subs += Register-ObjectEvent -InputObject $watcher -EventName Renamed -Action $action

try {
  while ($true) {
    if ($global:restartPending) {
      Restart-Server
    }
    Start-Sleep -Milliseconds 200
  }
} finally {
  Stop-Server
  foreach ($sub in $subs) {
    Unregister-Event -SourceIdentifier $sub.Name
  }
  $watcher.Dispose()
}
