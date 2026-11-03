param([string]$GoExe='go',[string]$ISCCExe='ISCC',[string]$OutputRoot,[string]$Version='0.1.0-preview.4',[string]$ProjectURL='')
$ErrorActionPreference='Stop'
$source=[IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..'))
if(-not $OutputRoot){$OutputRoot=Join-Path $source '..\Simkl-Bridge-Windows'}
$output=[IO.Path]::GetFullPath($OutputRoot)
if($output.StartsWith($source+[IO.Path]::DirectorySeparatorChar,[StringComparison]::OrdinalIgnoreCase)-or$output-eq$source){throw 'Bauausgabe muss außerhalb des öffentlichen Quellcodeordners liegen.'}
if($Version-notmatch '^\d+\.\d+\.\d+(?:-[A-Za-z0-9.]+)?$'){throw 'Ungültige Versionsnummer.'}
if($ProjectURL-and$ProjectURL-notmatch '^https://github\.com/[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+/?$'){throw 'Bitte die echte GitHub-Repository-Adresse verwenden.'}
Push-Location $source
try{
    & (Join-Path $PSScriptRoot 'check-release.tests.ps1')
    & (Join-Path $PSScriptRoot 'check-release.ps1') -SourceRoot $source
    & $GoExe test ./... -count=1
    if($LASTEXITCODE-ne 0){throw 'Programmtests fehlgeschlagen.'}
    $package=Join-Path $output 'Programm'
    New-Item -ItemType Directory -Force -Path $package | Out-Null
    & $GoExe run ./cmd/icon (Join-Path $source 'packaging\windows\simkl.ico')
    if($LASTEXITCODE-ne 0){throw 'Symbol konnte nicht erstellt werden.'}
    & $GoExe build -trimpath -ldflags "-H=windowsgui -s -w -X main.version=$Version -X main.projectURL=$ProjectURL" -o (Join-Path $package 'simkl-bridge.exe') ./cmd/simkl-bridge
    if($LASTEXITCODE-ne 0){throw 'Programm konnte nicht erstellt werden.'}
    & $GoExe run ./cmd/icon (Join-Path $source 'packaging\windows\simkl.ico') (Join-Path $package 'simkl-bridge.exe')
    if($LASTEXITCODE-ne 0){throw 'Programmsymbol konnte nicht eingebettet werden.'}
    foreach($name in @('README.md','README.de.md','LICENSE','THIRD_PARTY_NOTICES.md','THIRD_PARTY_NOTICES.de.md','docs','licenses')){Copy-Item -LiteralPath (Join-Path $source $name) -Destination $package -Recurse -Force}
    Copy-Item -LiteralPath (Join-Path $source 'packaging\windows\simkl.ico') -Destination (Join-Path $package 'simkl.ico') -Force
    function Make-Archive([string]$Root,[string]$Destination,[bool]$SourcePackage){
        Add-Type -AssemblyName System.IO.Compression
        $stream=[IO.File]::Open($Destination,[IO.FileMode]::Create)
        $zip=[IO.Compression.ZipArchive]::new($stream,[IO.Compression.ZipArchiveMode]::Create)
        try{foreach($file in Get-ChildItem -LiteralPath $Root -Recurse -File -Force){$relative=[IO.Path]::GetRelativePath($Root,$file.FullName).Replace('\','/');if($SourcePackage-and$relative.StartsWith('.git/')){continue};[IO.Compression.ZipFileExtensions]::CreateEntryFromFile($zip,$file.FullName,$relative)|Out-Null}}finally{$zip.Dispose();$stream.Dispose()}
    }
    $windowsZip=Join-Path $output "Simkl-Bridge-$Version-Windows-x64.zip"
    Make-Archive $package $windowsZip $false
    & (Join-Path $PSScriptRoot 'check-release.ps1') -SourceRoot $source -Archive $windowsZip -ArchiveKind windows
    $sourceZip=Join-Path $output "Simkl-Bridge-$Version-Quellcode.zip"
    Make-Archive $source $sourceZip $true
    & (Join-Path $PSScriptRoot 'check-release.ps1') -SourceRoot $source -Archive $sourceZip -ArchiveKind source
    & $ISCCExe "/DAppVersion=$Version" "/DPackageDir=$package" "/DOutputDir=$output" (Join-Path $source 'packaging\windows\setup.iss')
    if($LASTEXITCODE-ne 0){throw 'Installer konnte nicht erstellt werden.'}
    Get-ChildItem -LiteralPath $output -File | Where-Object {$_.Extension-in@('.exe','.zip')} | Get-FileHash -Algorithm SHA256 | Select-Object @{n='Datei';e={Split-Path $_.Path -Leaf}},Hash
}finally{Pop-Location}

