# Simkl Bridge für Windows

Eine eigenständige Windows-App im Infobereich für eigene Simkl-Listen in Sonarr und Radarr. Ohne Docker, Python oder Node.

**Stand: erste Vorabversion.** Die automatischen Prüfungen werden beim Paketbau ausgeführt. Die persönliche Einrichtung und die vollständige Sichtprüfung aller Windows-Skalierungen sind separate Schritte; die vorhandene Bridge wird nicht übernommen oder verändert.

[English](README.en.md) · [Einrichtung](docs/setup-simkl.de.md) · [Installation](docs/installation.de.md)

## Installieren

Den Windows-x64-Installer aus dem Paket beziehungsweise den GitHub-Artefakten herunterladen und ausführen. Er installiert für den angemeldeten Benutzer; Administratorrechte sind nicht vorgesehen. Die App erscheint anschließend als **Simkl Bridge** im Startmenü. Der erste Start öffnet die Einrichtung.

Eine selbst gebaute Vorabversion ist nicht digital signiert. Windows kann deshalb nach einer Bestätigung fragen. Prüfe Herkunft und Prüfsumme des Pakets.

## Einrichten

1. Eine eigene **AUTH-V2-App** bei Simkl registrieren: **TV, devices & command line**.
2. Unter **Anmeldung** die Client-ID eintragen, **Mit Simkl verbinden** wählen und den Code im Browser bestätigen. Die App fordert ausschließlich lesenden Zugriff an.
3. Unter **Listen** die eigenen Listen laden oder einen Simkl-Listenlink einfügen. Sonarr beziehungsweise Radarr als Ziel wählen.
4. Die eingerichtete Liste markieren und **Listenadresse kopieren** wählen.
5. In Sonarr **Benutzerdefinierte Liste**, in Radarr **Benutzerdefinierte Liste / Custom List** als Importquelle verwenden und die lokale Adresse einfügen. Beide Programme müssen auf demselben Rechner laufen.
6. Einstellungen speichern. Im Rechtsklick-Menü bei Bedarf **Mit Windows starten** einschalten.

Für Custom Lists verlangt Simkl **PRO oder VIP und AUTH V2**. Eine nicht verifizierte eigene App kann verwendet werden. Es werden keine Sonarr-/Radarr-Administratorschlüssel benötigt. [Ausführliche Anleitung](docs/setup-simkl.de.md)

**PRO geht auch kostenlos ohne bezahltes Abo:** Über Simkls Treueprogramm kannst du PRO freischalten. Die [deutsche Schritt-für-Schritt-Erklärung](docs/setup-simkl.de.md#kostenloses-pro-ohne-bezahltes-abo) und die [englische Erklärung](docs/setup-simkl.en.md#free-pro-without-a-paid-subscription) beschreiben Nutzungstage und monatliche Aktivierung.

## Bedienung

- Rechtsklick auf das Simkl-Symbol: starten, stoppen, neu starten, Autostart, Einstellungen, Protokoll und Über. Alle Unterseiten erscheinen im selben kompakten Menü; Zurück führt ins Hauptmenü.
- Grün: Bridge läuft. Gelb: Start oder Wiederanlauf. Rot: gestoppt oder verständlich benannter Fehler. Der Hinweis beim Darüberfahren zeigt denselben Zustand.
- **Bridge stoppen** lässt das Tray-Symbol bestehen. **Beenden** beendet auch den eigenen Bridge-Prozess.
- Nach Abstürzen höchstens drei Wiederanläufe nach 2, 5 und 15 Sekunden. Nach einer stabilen Minute beginnt die Zählung neu. Anmelde-/Einrichtungsfehler lösen keine Endlosschleife aus.
- Deutsch und Englisch lassen sich im Tray-Menü wählen. Ein zweiter Programmstart öffnet die vorhandene Bedienung.

## Listen und Datenschutz

Die App stellt nur eingerichtete Listen über `127.0.0.1` bereit. Sonarr erhält TVDB-Kennungen, Radarr TMDb-Kennungen. Episodische Anime und Anime-Filme werden passend zugeordnet. Fehlende Kennungen und unvollständige Antworten führen zu einem Fehler, nicht zu einer scheinbar leeren Erfolgsliste.

Unveränderte Inhalte werden zwischengespeichert, auch über einen Neustart hinweg. Bei eigenen Listen werden zuerst Aktivitäten und bei Bedarf die Listen-Zeitstempel geprüft. Automatische Listen werden bei Benutzung höchstens einmal täglich neu geladen. Fremde Listen werden über die Übersicht des jeweiligen Besitzers geprüft. Geänderte Custom Lists müssen gegebenenfalls vollständig neu gelesen werden; eine Änderungsschnittstelle für einzelne Einträge wird nicht vorgetäuscht. Es gibt keine ständige Hintergrundabfrage der Simkl-API. [Details](docs/refresh.md)

Einstellungen, Listencache und Protokoll liegen unter `%LOCALAPPDATA%\SimklBridge`. Die Anmeldung wird mit dem Windows-Benutzerschutz DPAPI gespeichert und an die Client-ID gebunden. Tokens gehören niemals in ein Repository. Aktualisierungen und die normale Deinstallation lassen diese Daten bestehen.

## Selbst bauen

Windows x64, Go **1.26.1** und Inno Setup **7.1.1** bereitstellen. Anschließend in PowerShell:

```powershell
./scripts/build.ps1 -ISCCExe 'Pfad-zum-Compiler/ISCC.exe'
```

Die Ausgabe liegt im benachbarten Ordner `Simkl-Bridge-Windows`: Programm, Installer, Windows-Archiv und bereinigtes Quellcodearchiv. Optional kann `-ProjectURL` mit der echten GitHub-Repository-Adresse gesetzt werden; dieselbe Adresse erscheint dann im Über-Dialog. Es wird keine fremde Repository-Adresse vorgegeben.

Der GitHub-Workflow lässt sich manuell starten und legt geprüfte Bauartefakte ab. Er veröffentlicht kein Release automatisch. Eine zusätzliche Race-Prüfung erfordert unter Windows eine geeignete C-Werkzeugkette; die funktionalen Gleichzeitigkeitstests laufen ohne sie.

## Quellen und Lizenz

Copyright (c) 2026 **Veyilla Nightwhisper**. MIT-Lizenz. Unter **Über** lassen sich die eigene Lizenz, Quellenhinweise und weitere Lizenztexte außerhalb der App öffnen.

Unabhängiges Projekt, kein offizielles Simkl-Produkt. [Simkl](https://simkl.com/) · [Offizielle API-Dokumentation](https://api.simkl.org/)

MIT-Lizenz; siehe [LICENSE](LICENSE). Angepasste Prozessüberwachung auf Grundlage von [Resurrector 1.2.0](https://github.com/Maki-Daisuke/resurrector/tree/v1.2.0). Symbol aus Dashboard Icons, Windows-Bibliothek aus dem Go-Projekt. Vollständige Hinweise und Lizenztexte liegen in [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md) und `licenses`.

