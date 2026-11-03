param([Parameter(Mandatory)][string]$SourceRoot,[string]$Archive,[ValidateSet('source','windows')][string]$ArchiveKind='source')
$ErrorActionPreference='Stop'
function Test-Content([string]$Name,[byte[]]$Bytes,[string]$Kind) {
    $name=$Name.Replace('\','/')
    if ($name -match '(^|/)\.\.?(/|$)' -or $name.StartsWith('/')) { throw "Ungültiger Paketpfad: $name" }
    if ($name -match '(^|/)(settings\.json|authorization\.dat|lists-cache\.json|bridge\.log|\.env)(/|$)') { throw "Private Laufzeitdatei im Paket: $name" }
    $rootFiles=@('README.md','README.de.md','README.en.md','LICENSE','THIRD_PARTY_NOTICES.md','THIRD_PARTY_NOTICES.de.md','go.mod','go.sum','.gitignore')
    $allowed=$rootFiles -contains $name
    if ($Kind -eq 'source') { $allowed=$allowed -or $name -match '^(assets/.+\.(png|go)|cmd/.+\.go|internal/.+\.(go|json)|docs/.+\.(md|png)|licenses/.+\.txt|packaging/.+\.(iss|ico)|scripts/.+\.ps1|\.github/workflows/.+\.yml)$' }
    else { $allowed=$allowed -or $name -match '^(simkl-bridge\.exe|simkl\.ico|docs/.+\.(md|png)|licenses/.+\.txt)$' }
    if($name.StartsWith('docs/')){$allowed=$name -in @('docs/setup-simkl.de.md','docs/setup-simkl.en.md','docs/refresh.md','docs/refresh.de.md','docs/installation.de.md','docs/installation.en.md','docs/images/simkl.png','docs/images/tray-menu.png','docs/images/tray-menu.en.png')}
    if (-not $allowed) { throw "Nicht freigegebene Datei: $name" }
    if ($name -match '\.(png|ico)$') { return }
    $text=[Text.Encoding]::UTF8.GetString($Bytes)
    if($name -notmatch '\.exe$'-and -not $name.StartsWith('licenses/')-and$text -match '(?i)\b(?:generated|assisted)\s+(?:by|with)\b'){throw "Unerwünschter Zusatzvermerk in: $name"}
    if ($text -match 'simkl_(at|rt|cs)_[A-Za-z0-9_-]{20,}') { throw "Zugangsdaten im Paket: $name" }
    if ($text -match '(?i)[A-Z]:\\Users\\[^\\\s"''<>]+') { throw "Persönlicher Rechnerpfad im Paket: $name" }
    if($name -notmatch '\.exe$'){foreach($match in [regex]::Matches($text,'(?i)(?<![a-z0-9])[a-f0-9]{64}(?![a-z0-9])')) { if ($match.Value -notmatch '^([a-fA-F0-9])\1{63}$') { throw "Ungeklärter privater Kennwert in: $name" } }}
}
$resolved=(Resolve-Path -LiteralPath $SourceRoot).Path
foreach($required in @('README.md','LICENSE','THIRD_PARTY_NOTICES.md','go.mod')) { if(-not(Test-Path -LiteralPath (Join-Path $resolved $required))){throw "Erforderliche Datei fehlt: $required"} }
foreach($file in Get-ChildItem -LiteralPath $resolved -Recurse -File -Force) {
    $relative=[IO.Path]::GetRelativePath($resolved,$file.FullName).Replace('\','/')
    if($relative.StartsWith('.git/')){continue}
    Test-Content $relative ([IO.File]::ReadAllBytes($file.FullName)) 'source'
}
if($Archive){
    Add-Type -AssemblyName System.IO.Compression
    $zip=[IO.Compression.ZipFile]::OpenRead((Resolve-Path -LiteralPath $Archive).Path)
    try { foreach($entry in $zip.Entries){if($entry.FullName.EndsWith('/')){continue};$stream=$entry.Open();$buffer=[IO.MemoryStream]::new();try{$stream.CopyTo($buffer);Test-Content $entry.FullName $buffer.ToArray() $ArchiveKind}finally{$stream.Dispose();$buffer.Dispose()}} }finally{$zip.Dispose()}
}
Write-Output 'Paketinhalt geprüft.'
