# Install, update and uninstall

<strong><a href="installation.en.md"><img src="https://github.githubassets.com/images/icons/emoji/unicode/1f1ec-1f1e7.png" width="20" height="20" alt=""> English</a></strong> · <a href="installation.de.md"><img src="https://github.githubassets.com/images/icons/emoji/unicode/1f1e9-1f1ea.png" width="20" height="20" alt=""> Deutsch</a>

[README](../README.md) · [Setup guide](setup-simkl.en.md)

## Install a preview

This repository does not have a downloadable release yet. These steps apply if you already have a preview installer. Otherwise, see [Building from source](../README.md#building-from-source). When releases are available, download them from this repository's [Releases page](https://github.com/Veyilla/simkl-bridge-windows/releases).

1. Close a running copy of Simkl Bridge using **Quit** in its tray menu.
2. Start the file named `Simkl-Bridge-…-Windows-x64-Setup.exe`.
3. Select English or German and complete the installation for your Windows user account.
4. Open **Simkl Bridge** from the Start menu. The first launch opens the sign-in and list settings.
5. Continue with the [step-by-step setup guide](setup-simkl.en.md).

The app needs no Docker, Python or Node runtime. Installation is per user; administrator rights are not intended to be necessary. The installer does not migrate the configuration of other bridge apps.

## Update

1. In Simkl Bridge, choose **Quit**.
2. Run the new installer and use the same installation folder.
3. Open Simkl Bridge again. Your saved sign-in and list assignments are reused.

If you change the bridge port later, copy each updated list address into Sonarr or Radarr.

## Windows startup

Right-click the S icon next to the Windows clock and enable **Start with Windows**. Disable the same option if you prefer to launch the app manually. This setting applies to the current Windows user.

## Uninstall

Quit the app, then uninstall **Simkl Bridge** through Windows Settings. Standard uninstall removes its program files and its own startup entry. Saved sign-in, settings and cached lists remain in `%LOCALAPPDATA%\SimklBridge`.

## Preview details

- The installer is unsigned. Check the package's source and any supplied checksum before running it.
- The setup uses black backgrounds and white text. Its built-in dark controls, title bar and scrollbars may be dark gray. Windows security prompts use Windows' own appearance; a high-contrast theme takes precedence.
- Not every physical display arrangement or long-term scenario has been verified. Report reproducible problems through [Issues](https://github.com/Veyilla/simkl-bridge-windows/issues), without tokens or other private information.

For sign-in errors, empty lists or a port already in use, see [Troubleshooting](setup-simkl.en.md#troubleshooting).
