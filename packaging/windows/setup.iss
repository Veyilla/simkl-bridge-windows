#ifndef AppVersion
  #define AppVersion "0.1.0-preview.4"
#endif
#ifndef PackageDir
  #error PackageDir muss auf den geprüften Windows-Paketordner zeigen.
#endif
#ifndef OutputDir
  #error OutputDir muss gesetzt sein.
#endif
[Setup]
AppId={{D1D856AF-1F95-4DD0-AD53-014BBC3BED0D}
AppName=Simkl Bridge
AppVersion={#AppVersion}
AppCopyright=Copyright (c) 2026 Veyilla Nightwhisper
DefaultDirName={localappdata}\Programs\Simkl Bridge
DefaultGroupName=Simkl Bridge
PrivilegesRequired=lowest
ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible
OutputDir={#OutputDir}
OutputBaseFilename=Simkl-Bridge-{#AppVersion}-Windows-x64-Setup
SetupIconFile=simkl.ico
UninstallDisplayIcon={app}\simkl.ico
Compression=lzma2
SolidCompression=yes
WizardStyle=modern dark includetitlebar hidebevels
WizardBackColor=#000000
WizardSmallImageFile=
WizardImageFile=
CloseApplications=yes
RestartApplications=no
LicenseFile=..\..\LICENSE
DisableProgramGroupPage=yes
[Languages]
Name: "german"; MessagesFile: "compiler:Languages\German.isl"
Name: "english"; MessagesFile: "compiler:Default.isl"
[Files]
Source: "{#PackageDir}\simkl-bridge.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#PackageDir}\simkl.ico"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#PackageDir}\README.de.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#PackageDir}\README.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#PackageDir}\LICENSE"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#PackageDir}\THIRD_PARTY_NOTICES.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#PackageDir}\THIRD_PARTY_NOTICES.de.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#PackageDir}\docs\*"; DestDir: "{app}\docs"; Flags: ignoreversion recursesubdirs
Source: "{#PackageDir}\licenses\*"; DestDir: "{app}\licenses"; Flags: ignoreversion recursesubdirs
[Icons]
Name: "{autoprograms}\Simkl Bridge"; Filename: "{app}\simkl-bridge.exe"; IconFilename: "{app}\simkl.ico"
[CustomMessages]
german.LaunchBridge=Simkl Bridge starten
english.LaunchBridge=Start Simkl Bridge
[Run]
Filename: "{app}\simkl-bridge.exe"; Description: "{cm:LaunchBridge}"; Flags: nowait postinstall skipifsilent unchecked
[Code]
procedure InitializeWizard;
begin
  if not HighContrastActive then begin
    WizardForm.LicenseMemo.UseRichEdit := False;
    WizardForm.LicenseMemo.StyleElements := WizardForm.LicenseMemo.StyleElements - [seClient, seFont];
    WizardForm.LicenseMemo.Color := clBlack;
    WizardForm.LicenseMemo.Font.Color := clWhite;
  end;
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
var Existing, Expected: String;
begin
  if CurUninstallStep = usUninstall then begin
    Expected := '"' + ExpandConstant('{app}\simkl-bridge.exe') + '" --autostart';
    if RegQueryStringValue(HKCU, 'Software\Microsoft\Windows\CurrentVersion\Run', 'SimklBridge', Existing) then
      if Existing = Expected then
        RegDeleteValue(HKCU, 'Software\Microsoft\Windows\CurrentVersion\Run', 'SimklBridge');
  end;
end;

