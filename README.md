<p align="center">
  <img src="docs/images/simkl.png" width="88" height="88" alt="Simkl Bridge icon: a white S on black">
</p>

<h1 align="center">Simkl Bridge for Windows<br><sub>Your lists. Directly connected.</sub></h1>

<p align="center">
  Import Simkl Custom Lists into Sonarr and Radarr with a compact Windows tray app.<br>
  No Docker, Python or Node runtime required.
</p>

<p align="center">
  <a href="docs/installation.en.md"><img src="https://img.shields.io/badge/Windows-x64-0078D4?style=flat-square" alt="Windows x64"></a>
  <a href="https://github.com/Veyilla/simkl-bridge-windows/commits/main/"><img src="https://img.shields.io/badge/Version-0.1.0--preview.4-D9A441?style=flat-square" alt="Version 0.1.0-preview.4"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-3DA66D?style=flat-square" alt="License: MIT"></a>
</p>

<p align="center">
  <strong><a href="README.md"><img src="https://github.githubassets.com/images/icons/emoji/unicode/1f1ec-1f1e7.png" width="20" height="20" alt=""> English</a></strong> · <a href="README.de.md"><img src="https://github.githubassets.com/images/icons/emoji/unicode/1f1e9-1f1ea.png" width="20" height="20" alt=""> Deutsch</a>
</p>

<p align="center">
  <a href="docs/installation.en.md">Installation</a> · <a href="docs/setup-simkl.en.md">Setup guide</a> · <a href="docs/setup-simkl.en.md#troubleshooting">Troubleshooting</a>
</p>

## What it does

Simkl Bridge is for anyone who wants to use their own Simkl lists in Sonarr and Radarr with as little setup as possible. A compact tray menu brings account connection, list selection and everyday controls together. The [setup guide](docs/setup-simkl.en.md) walks you through your first list; experienced users can start with the quick setup below.

- After setup, automatically retrieves your configured Simkl Custom Lists when Sonarr or Radarr requests them on its import schedule, checks for changes and passes complete lists to the destination app. Keep the bridge running; no manual export or copying is needed for subsequent updates.
- Provides start, stop, restart, Windows startup, settings and logs from one tray menu.
- Keeps sign-in and list settings between sessions. No console window needs to stay open.
- Supports an AMOLED-black interface with German and English labels.

**Current version:** 0.1.0-preview.4. This is a preview release, and the Windows packages are unsigned. See [installation notes](docs/installation.en.md#preview-details) for the current limitations.

## Requirements

- Windows x64, with Sonarr and/or Radarr running on the same computer.
- A Simkl account with **PRO or VIP**, including free Loyalty PRO.
- A personal Simkl **AUTH V2** app registration. You do not need a client secret or Sonarr/Radarr administrator API keys.

**Free PRO:** a paid subscription is not required. The [setup guide](docs/setup-simkl.en.md#free-pro-without-a-paid-subscription) explains Simkl's activity requirement and monthly activation.

## Installation

No downloadable release has been published in this repository yet. If you already have a preview installer, follow the [installation guide](docs/installation.en.md). Otherwise, use the [build instructions](#building-from-source). Future downloads will be listed under [Releases](https://github.com/Veyilla/simkl-bridge-windows/releases).

## Quick setup

The app initially opens in German. If settings are shown, choose **Zurück**, then **English** at the bottom of the main menu. Reopen **Settings** to follow the English labels below.

1. Open Simkl Bridge. Under **Sign-in**, enter the Client ID of your AUTH V2 app.
2. Choose **Connect to Simkl** and approve the displayed code in your browser.
3. Under **Lists**, load your lists or paste a Simkl list link. Assign the list to Sonarr or Radarr.
4. Select the configured list, choose **Copy address**, then **Save and back**. Start the bridge if it is stopped.
5. Add the copied address as an import list, following [Connect Sonarr and Radarr](#connect-sonarr-and-radarr).

Repeat the list steps for each destination. For registration, exact menu paths and explanations, see the [step-by-step setup guide](docs/setup-simkl.en.md).

## Connect Sonarr and Radarr

Signing in to Simkl is only part of setup: create an import list in the destination app for each list you want to use. Sonarr and Radarr must run on the same Windows computer as the bridge.

| Destination | Assign the list in the bridge | Menu path in the destination app |
|---|---|---|
| Sonarr for series | **Lists → Selection → Sonarr** | **Settings → Import Lists → + → Custom List** |
| Radarr for movies | **Lists → Selection → Radarr** | **Settings → Import Lists → + → Custom List** |

With a German interface, the path in both apps is **Einstellungen → Importlisten → + → Benutzerdefinierte Liste**.

1. Under **Lists** in the bridge, select the configured entry and choose **Copy address**. Use the bridge address, rather than the Simkl website link.
2. Choose **Save and back** and make sure the status says **Bridge running**. The service pauses while you edit the bridge settings.
3. Open the destination app's menu path shown above. Give the import list a recognizable name and paste the copied address into its URL field.
4. Enable automatic adding if you want new titles added to the library automatically. Select a valid **Root Folder** for series or movies and a **Quality Profile**. Also choose whether new titles should be monitored and missing episodes or movies searched for immediately.
5. Review the library-cleaning options: they determine what happens to titles that later disappear from lists.
6. Choose **Test**, then **Save**. Run the destination app's list refresh action or wait for its next scheduled request.

Repeat setup for each list and destination app. List requests then run automatically while the bridge is running. If you change the bridge port, update the addresses in Sonarr/Radarr. See [Troubleshooting](#troubleshooting) if something fails.

## Everyday use

Right-click the S icon next to the Windows clock. **Stop bridge** stops the background service but keeps the tray icon available. **Quit** closes the app and its service. Enable **Start with Windows** for automatic startup.

The status text and hover tooltip explain the current state: green means running, amber means starting or restarting, and red means stopped or an error. **About** opens the project page, license and third-party notices.

Choose **Restart bridge** to restart the service. After an unexpected exit, the app attempts up to three automatic restarts, then stops and shows the error. **Settings** lets you reconnect your account, manage lists and change the local port. Save and leave settings to resume the service.

<p>
  <img src="docs/images/tray-menu.en.png" width="368" alt="Simkl Bridge tray menu with status, controls, settings and language selection">
</p>

*Menu preview; the connection status shown is an example.*

## Refreshing lists and protecting your data

The bridge checks for changes when a list is requested and reuses unchanged contents. Sonarr and Radarr always receive complete lists. API failures are reported instead of being returned as empty lists. [How refreshing works](docs/refresh.md)

Only configured lists are served on `127.0.0.1`. Settings, cached lists and logs stay in `%LOCALAPPDATA%\SimklBridge`. Sign-in tokens are protected with Windows DPAPI. Updates and standard uninstall preserve these data files.

## Troubleshooting

See the [troubleshooting table](docs/setup-simkl.en.md#troubleshooting) for sign-in errors, empty lists, PRO status and port conflicts. Choose **Open log** in the tray menu. The log records timestamps, recognized service and list errors, and restart states. Repeated unchanged list-error states are not logged on every request; it is not a complete record of every action. Log messages currently use German even with the English interface.

If the problem persists, [open an issue](https://github.com/Veyilla/simkl-bridge-windows/issues/new) and include:

- The app version, Windows version and whether Sonarr or Radarr is affected.
- The steps to reproduce the problem, expected result and actual error message.
- The time of the failure and the relevant lines from `%LOCALAPPDATA%\SimklBridge\bridge.log`.
- For display problems, a screenshot, screen resolution and Windows scaling percentage.

Check the excerpt before posting and remove tokens, personal list addresses and other private details. If the log contains no matching entry, mention that in the issue and include the error shown in the app or Sonarr/Radarr.

## Building from source

On Windows, install Go **1.26.1** and Inno Setup **7.1.1**, then run from the repository folder:

```powershell
./scripts/build.ps1 -ISCCExe 'C:/Path/to/ISCC.exe' -ProjectURL 'https://github.com/Veyilla/simkl-bridge-windows'
```

The script runs the checks and creates the app, installer, Windows archive and source archive in the adjacent `Simkl-Bridge-Windows` folder. Race checking additionally requires a Windows C toolchain.

## License and credits

Copyright © 2026 **Veyilla Nightwhisper**. Licensed under the [MIT License](LICENSE).

This is an independent app, not an official Simkl product. It uses adapted components from [Resurrector](https://github.com/Maki-Daisuke/resurrector/tree/v1.2.0) by [Daisuke (yet another) Maki](https://github.com/Maki-Daisuke), the Dashboard Icons Simkl image and the Go Windows library. Their original notices and license texts are included in [Third-party notices](THIRD_PARTY_NOTICES.md) and the [licenses folder](licenses).
