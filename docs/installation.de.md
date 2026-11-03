# Installieren, aktualisieren und deinstallieren

[English](installation.en.md) · **Deutsch**

[README](../README.de.md) · [Einrichtungsanleitung](setup-simkl.de.md)

## Vorabversion installieren

In diesem Repository gibt es noch kein Release zum Herunterladen. Diese Schritte gelten, wenn du bereits einen Vorabversions-Installer hast. Andernfalls kannst du die App [selbst bauen](../README.de.md#selbst-bauen). Sobald Releases verfügbar sind, findest du sie auf der [Releases-Seite](https://github.com/Veyilla/simkl-bridge-windows/releases) dieses Repositorys.

1. Schließe eine laufende Simkl Bridge über **Beenden** im Menü neben der Uhr.
2. Starte die Datei `Simkl-Bridge-…-Windows-x64-Setup.exe`.
3. Wähle Deutsch oder Englisch und schließe die Installation für dein Windows-Benutzerkonto ab.
4. Öffne **Simkl Bridge** über das Startmenü. Beim ersten Start erscheinen Anmeldung und Listeneinstellungen.
5. Folge der [Schritt-für-Schritt-Anleitung](setup-simkl.de.md).

Docker, Python und Node werden nicht benötigt. Die Installation gilt für den angemeldeten Benutzer; Administratorrechte sollen nicht erforderlich sein. Der Installer übernimmt keine Konfiguration anderer Bridge-Apps.

## Aktualisieren

1. Wähle in Simkl Bridge **Beenden**.
2. Führe den neuen Installer aus und verwende denselben Installationsordner.
3. Öffne Simkl Bridge erneut. Die gespeicherte Anmeldung und Listenzuordnung werden wiederverwendet.

Wenn du später den Bridge-Port änderst, kopiere die aktualisierten Listenadressen in Sonarr oder Radarr.

## Mit Windows starten

Klicke mit der rechten Maustaste auf das S-Symbol neben der Windows-Uhr und aktiviere **Mit Windows starten**. Deaktiviere dieselbe Option, wenn du die App lieber manuell öffnest. Die Einstellung gilt für den aktuellen Windows-Benutzer.

## Deinstallieren

Beende die App und deinstalliere **Simkl Bridge** über die Windows-Einstellungen. Die normale Deinstallation entfernt die Programmdateien und den eigenen Autostart-Eintrag. Anmeldung, Einstellungen und Listencache bleiben in `%LOCALAPPDATA%\SimklBridge` erhalten.

## Hinweise zur Vorabversion

- Der Installer ist nicht digital signiert. Prüfe vor dem Start die Herkunft des Pakets und gegebenenfalls die mitgelieferte Prüfsumme.
- Das Setup verwendet schwarze Flächen und weiße Schrift. Die eingebauten dunklen Schaltflächen, die Titelleiste und die Scrollbalken können dunkelgrau sein. Windows-Sicherheitsabfragen verwenden die Windows-Darstellung; ein Kontrastdesign hat Vorrang.
- Noch nicht jede physische Monitoranordnung und jeder Langzeitfall ist geprüft. Melde nachvollziehbare Probleme unter [Issues](https://github.com/Veyilla/simkl-bridge-windows/issues), ohne Tokens oder andere private Angaben.

Hinweise zu Anmeldefehlern, leeren Listen und belegten Ports stehen unter [Fehlerbehebung](setup-simkl.de.md#fehlerbehebung).
