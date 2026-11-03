$ErrorActionPreference = 'Stop'
$root = Join-Path ([IO.Path]::GetTempPath()) ('bridge-pruefung-' + [guid]::NewGuid())
New-Item -ItemType Directory -Path $root | Out-Null
foreach ($name in @('README.md','LICENSE','THIRD_PARTY_NOTICES.md','go.mod')) { Set-Content -LiteralPath (Join-Path $root $name) -Value 'Beispiel' }
$checker = Join-Path $PSScriptRoot 'check-release.ps1'
function Must-Fail([string]$Name,[string]$Content) {
    $fixture = Join-Path ([IO.Path]::GetTempPath()) ('bridge-negativ-' + [guid]::NewGuid())
    New-Item -ItemType Directory -Path $fixture | Out-Null
    foreach ($requiredFile in @('README.md','LICENSE','THIRD_PARTY_NOTICES.md','go.mod')) { Set-Content -LiteralPath (Join-Path $fixture $requiredFile) -Value 'Beispiel' }
    $path = Join-Path $fixture $Name
    New-Item -ItemType Directory -Force -Path (Split-Path $path) | Out-Null
    Set-Content -LiteralPath $path -Value $Content
    $rejected = $false
    try { & $checker -SourceRoot $fixture | Out-Null } catch { $rejected = $true }
    if (-not $rejected) { throw "Unerlaubter Paketinhalt akzeptiert: $Name" }
}
& $checker -SourceRoot $root | Out-Null
Must-Fail '.env' 'private=example'
Must-Fail 'docs/example.md' ('simkl_at_' + ('x' * 40))
Must-Fail 'docs/example.md' ('C:' + '\Users\Example\private.txt')
Must-Fail 'arbeitsplan.md' 'Interne Notizen'
Must-Fail 'docs/internal-notes.md' 'Interne Notizen'
Must-Fail 'docs/setup-simkl.de.md' ('Generated' + ' by Example tool')
$archive=Join-Path ([IO.Path]::GetTempPath()) ('bridge-negativ-'+[guid]::NewGuid()+'.zip')
$zip=[IO.Compression.ZipFile]::Open($archive,[IO.Compression.ZipArchiveMode]::Create)
try{$entry=$zip.CreateEntry('docs/setup-simkl.de.md');$writer=[IO.StreamWriter]::new($entry.Open());$writer.Write('Generated'+' by Example tool');$writer.Dispose()}finally{$zip.Dispose()}
$rejected=$false
try{& $checker -SourceRoot $root -Archive $archive|Out-Null}catch{$rejected=$true}
if(-not $rejected){throw 'Unerwünschter Zusatzvermerk im tatsächlichen Archiv wurde akzeptiert.'}
Write-Output 'Paketprüfung: sauberes Paket akzeptiert; Zugangsdaten, Rechnerpfade und fremde Dateien abgewiesen.'
