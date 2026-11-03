# Simkl Bridge for Windows

A standalone Windows tray app for using Simkl Custom Lists in Sonarr and Radarr. No Docker, Python or Node runtime is required.

**Status: first preview build.** Automated checks run during packaging. Personal setup and complete visual checks across Windows scaling settings remain separate. The existing personal bridge is not migrated or modified.

[Deutsch](README.md) · [Setup guide](docs/setup-simkl.en.md)

Run the Windows x64 installer, then launch Simkl Bridge from the Start menu. Installation is per user. Preview builds are unsigned.

1. Register an AUTH V2 app using **TV, devices & command line**.
2. Paste its client ID under **Sign-in**, choose **Connect to Simkl**, then approve the displayed code in your browser. Only read access is requested.
3. Load your lists or paste a Simkl list link under **Lists**. Select Sonarr or Radarr as the destination.
4. Select a configured list and copy its local address into a Custom List import in the destination application. Both applications must run on the same computer.
5. Save settings and optionally enable Windows startup from the tray menu.

Custom Lists require Simkl PRO or VIP and AUTH V2. An unverified personal app can be used. Sonarr/Radarr administrator API keys are not required.

**PRO is also available for free without a paid subscription:** See the [English instructions](docs/setup-simkl.en.md#free-pro-without-a-paid-subscription) or [German instructions](docs/setup-simkl.de.md#kostenloses-pro-ohne-bezahltes-abo) for the Loyalty Program's activity days and monthly activation.

The tray provides start, stop, restart, startup, settings, log and About. Green means running, amber means starting, and red means stopped or an actionable failure. Hover for status. Stop keeps the tray running; Quit also terminates the owned worker. Crashes are retried at most three times after 2, 5 and 15 seconds, resetting after a stable minute.

Only configured lists are exposed on `127.0.0.1`. Sonarr receives TVDB IDs and Radarr receives TMDb IDs. Unchanged contents are cached across restarts. Custom-list activity and metadata control refreshes; auto lists refresh at most daily on demand. Changed lists may require a complete download. Errors never become successful empty lists.

Local settings, list cache and log live under `%LOCALAPPDATA%\SimklBridge`. Authentication is protected with Windows DPAPI and bound to the client ID. Updates and standard uninstall preserve user data.

Build on Windows with Go 1.26.1 and Inno Setup 7.1.1 using `./scripts/build.ps1 -ISCCExe 'path/to/ISCC.exe'`. Optionally set `-ProjectURL` to the actual repository address. Output includes the installer and checked source archive. The GitHub workflow creates artifacts without automatically publishing a release. Race checking additionally requires a Windows C toolchain.

Independent project; not an official Simkl product. [Simkl](https://simkl.com/) · [API documentation](https://api.simkl.org/). MIT license; see LICENSE and THIRD_PARTY_NOTICES.md. Process supervision adapts Resurrector 1.2.0, with its original license included.

Copyright (c) 2026 **Veyilla Nightwhisper**. **About** opens the license, third-party notices and other license texts outside the app.
