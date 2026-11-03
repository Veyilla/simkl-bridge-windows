<p align="center">
  <img src="docs/images/simkl.png" width="88" height="88" alt="Simkl-Bridge-Symbol: weißes S auf Schwarz">
</p>

<h1 align="center">Simkl Bridge für Windows<br><sub>Deine Listen. Direkt verbunden.</sub></h1>

<p align="center">
  Importiere Simkl-Custom-Lists mit einer kompakten Windows-App in Sonarr und Radarr.<br>
  Docker, Python und Node werden nicht benötigt.
</p>

<p align="center">
  <a href="docs/installation.de.md"><img src="https://img.shields.io/badge/Windows-x64-0078D4?style=flat-square" alt="Windows x64"></a>
  <a href="https://github.com/Veyilla/simkl-bridge-windows/commits/main/"><img src="https://img.shields.io/badge/Version-0.1.0--preview.4-D9A441?style=flat-square" alt="Version 0.1.0-preview.4"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-3DA66D?style=flat-square" alt="Lizenz: MIT"></a>
</p>

<p align="center">
  <a href="README.md"><img src="https://github.githubassets.com/images/icons/emoji/unicode/1f1ec-1f1e7.png" width="20" height="20" alt=""> English</a> · <strong><a href="README.de.md"><img src="https://github.githubassets.com/images/icons/emoji/unicode/1f1e9-1f1ea.png" width="20" height="20" alt=""> Deutsch</a></strong>
</p>

<p align="center">
  <a href="docs/installation.de.md">Installation</a> · <a href="docs/setup-simkl.de.md">Einrichtung</a> · <a href="docs/setup-simkl.de.md#fehlerbehebung">Fehlerbehebung</a>
</p>

## Was die App macht

Simkl Bridge richtet sich an alle, die ihre eigenen Simkl-Listen mit möglichst wenig Einrichtungsaufwand in Sonarr und Radarr nutzen möchten. Das kompakte Menü neben der Windows-Uhr bündelt Anmeldung, Listenauswahl und die tägliche Bedienung. Die [Einrichtungsanleitung](docs/setup-simkl.de.md) führt durch die erste Liste; erfahrene Nutzer können direkt mit der Kurzanleitung unten beginnen.

- Ruft nach der Einrichtung deine konfigurierten Simkl-Custom-Lists automatisch ab, sobald Sonarr oder Radarr sie nach dem eigenen Importzeitplan anfordert, prüft auf Änderungen und gibt vollständige Listen an das Zielprogramm weiter. Lass die Bridge dafür laufen; bei späteren Aktualisierungen musst du nichts manuell exportieren oder kopieren.
- Bietet die tägliche Steuerung, Windows-Autostart, Einstellungen und Protokoll in einem Menü neben der Uhr.
- Überwacht den eigenen Bridge-Dienst und startet ihn nach einem Absturz automatisch neu, mit begrenzten Wiederholungsversuchen.
- Speichert Anmeldung und Listeneinstellungen. Es muss kein Konsolenfenster geöffnet bleiben.
- Bietet eine AMOLED-schwarze Oberfläche mit deutscher und englischer Beschriftung.

### Listen und unterstützte Medien

- Füge mehrere Listen über die Auswahl deines Kontos oder über einen zugänglichen Simkl-Listenlink beziehungsweise eine Listen-ID hinzu. Aus unterstützten Links erkennt die Bridge die ID automatisch. Zugängliche Listen anderer Nutzer lassen sich ebenfalls per Link hinzufügen; Simkls Zugriffsrechte gelten weiterhin.
- Serien gehen an Sonarr, Filme an Radarr. Bei Anime-Listen erhält Sonarr die Einträge, die keine Filme sind, und Radarr die Anime-Filme. Ordne dieselbe Anime-Liste beiden Zielen zu, wenn du beides übernehmen möchtest. Die benötigten TVDB-/TMDb-Zuordnungen müssen verfügbar sein; fehlende Zuordnungen führen zu einer Fehlermeldung. Doppelte Zielkennungen erscheinen pro Listenantwort nur einmal.
- Die Bridge liest Custom Lists. Erstelle oder bearbeite deren Inhalte auf Simkl. Eingebaute Kategorien wie Watching oder Plan to Watch werden über diese App nicht importiert. Überwachung, Suche und Downloads übernimmt Sonarr/Radarr.

**Aktuelle Version:** 0.1.0-preview.4. Dies ist eine Vorabversion; die Windows-Pakete sind nicht digital signiert. Die [Installationshinweise](docs/installation.de.md#hinweise-zur-vorabversion) erläutern die aktuellen Einschränkungen.

## Voraussetzungen

- Windows x64; Sonarr und/oder Radarr laufen auf demselben Rechner.
- Ein Simkl-Konto mit **PRO oder VIP**, einschließlich kostenlosem Loyalty-PRO.
- Eine eigene Simkl-App-Registrierung mit **AUTH V2**. Ein Client Secret oder Administratorschlüssel von Sonarr/Radarr wird nicht benötigt.

**Kostenloses PRO:** Ein bezahltes Abo ist nicht erforderlich. Die [Einrichtungsanleitung](docs/setup-simkl.de.md#kostenloses-pro-ohne-bezahltes-abo) erklärt die benötigten Nutzungstage und die monatliche Freischaltung.

**Free-Konten und AUTH V2:** AUTH V2 unterstützt auch Free-Konten. Das Tageskontingent beträgt 500 anrechenbare Anfragen für Free, 1.000 für PRO und 10.000 für VIP, gemeinsam für die verbundenen V2-Apps. [Offizielle Abfragelimits](https://api.simkl.org/resources/rate-limits)

Das Lesen der **Custom-List-Inhalte**, die diese Bridge importiert, verlangt zusätzlich PRO oder VIP. Free-Konten erhalten `premium_only` statt Listeneinträgen; seltenere Abrufe schalten diese nicht frei. [Offizielle Voraussetzungen für Custom Lists](https://api.simkl.org/guides/custom-lists)

## Installation

Ein Windows-Installer als **EXE-Datei** ist vorbereitet und soll demnächst unter [Releases](https://github.com/Veyilla/simkl-bridge-windows/releases) bereitgestellt werden. Ein Download-Release ist hier derzeit noch nicht veröffentlicht; ein Termin steht noch nicht fest. Der Dateiname folgt dem Muster `Simkl-Bridge-<Version>-Windows-x64-Setup.exe`.

Das Setup führt auf Deutsch oder Englisch durch die Lizenzvereinbarung und Installation, kopiert die App samt Anleitungen und Lizenztexten in deinen Benutzer-Programmordner und legt einen Startmenü-Eintrag an. Am Ende kannst du die App optional starten. Beim ersten Start verbindest du dein Simkl-Konto und wählst deine Listen; den Windows-Autostart aktivierst du anschließend bei Bedarf im App-Menü. Docker, Python und Node musst du nicht installieren.

Wenn du bereits einen Vorabversions-Installer hast, folge der [Installationsanleitung](docs/installation.de.md). Bis zum Download-Release kannst du die App alternativ [selbst bauen](#selbst-bauen).

## Schnelle Einrichtung

Die App startet zunächst auf Deutsch. Du kannst die Sprache unten im Hauptmenü wechseln. Wenn gerade die Einstellungen geöffnet sind, gelangst du mit **Zurück** dorthin.

1. Öffne Simkl Bridge. Trage unter **Anmeldung** die Client-ID deiner AUTH-V2-App ein.
2. Wähle **Mit Simkl verbinden** und bestätige den angezeigten Code im Browser.
3. Lade unter **Listen** deine Listen oder füge einen Simkl-Listenlink ein. Ordne die Liste Sonarr oder Radarr zu.
4. Markiere die eingerichtete Liste, wähle **Adresse kopieren** und anschließend **Speichern und zurück**. Starte die Bridge, falls sie gestoppt ist.
5. Richte die kopierte Adresse wie unter [Sonarr und Radarr verbinden](#sonarr-und-radarr-verbinden) beschrieben als Importliste ein.

Wiederhole die Listenschritte für jedes Ziel. Registrierung, genaue Menüpfade und Erklärungen stehen in der [Schritt-für-Schritt-Anleitung](docs/setup-simkl.de.md).

## Sonarr und Radarr verbinden

Die Anmeldung bei Simkl allein reicht nicht: Lege für jede gewünschte Liste auch eine Importliste im Zielprogramm an. Sonarr und Radarr müssen auf demselben Windows-Rechner wie die Bridge laufen.

| Zielprogramm | Liste in der Bridge zuordnen | Menüpfad im Zielprogramm |
|---|---|---|
| Sonarr für Serien | **Listen → Auswahl → Sonarr** | **Einstellungen → Importlisten → + → Benutzerdefinierte Liste** |
| Radarr für Filme | **Listen → Auswahl → Radarr** | **Einstellungen → Importlisten → + → Benutzerdefinierte Liste** |

Bei englischer Oberfläche lautet der Menüpfad in beiden Programmen **Settings → Import Lists → + → Custom List**.

1. Markiere in der Bridge unter **Listen** den eingerichteten Eintrag und wähle **Adresse kopieren**. Verwende die Bridge-Adresse, nicht den Link zur Simkl-Webseite.
2. Wähle **Speichern und zurück** und stelle sicher, dass **Bridge läuft** angezeigt wird. Während du die Bridge-Einstellungen bearbeitest, pausiert der Dienst.
3. Öffne im Zielprogramm den oben genannten Menüpfad. Gib der Importliste einen erkennbaren Namen und füge die kopierte Adresse in das URL-Feld ein.
4. Aktiviere das automatische Hinzufügen, wenn neue Titel automatisch in die Bibliothek übernommen werden sollen. Wähle einen gültigen **Root-Ordner** für Serien beziehungsweise Filme und ein **Qualitätsprofil**. Lege außerdem fest, ob neue Titel überwacht und fehlende Episoden beziehungsweise Filme sofort gesucht werden sollen.
5. Prüfe die Aufräumoptionen der Mediathek: Sie bestimmen, was mit Titeln geschieht, die später aus Listen verschwinden.
6. Wähle **Prüfen**, dann **Speichern**. Starte im Zielprogramm die Listenaktualisierung oder warte auf den nächsten geplanten Abruf.

Wiederhole die Einrichtung für jede Liste und jedes Zielprogramm. Danach laufen die Listenabrufe automatisch, solange die Bridge läuft. Nach einer Änderung des Bridge-Ports musst du die neuen Adressen erneut in Sonarr/Radarr eintragen. Bei Problemen hilft die [Fehlerbehebung](#fehlerbehebung).

## Bedienung und Listenverwaltung

### Listen verwalten

Markiere unter **Einstellungen → Listen** einen eingerichteten Eintrag und wähle **Zuordnung entfernen**, dann **Speichern und zurück**, um diese Adresse nicht mehr bereitzustellen. Entfernt wird die Bridge-Zuordnung. Die Simkl-Liste und Titel in Sonarr/Radarr werden dadurch nicht gelöscht; verwalte die zugehörige Importliste auch im Zielprogramm.

### Bedienung im Menü

Klicke mit der rechten Maustaste auf das S-Symbol neben der Windows-Uhr. **Bridge stoppen** beendet den Hintergrunddienst; das Symbol bleibt verfügbar. **Beenden** schließt die App und ihren Dienst. Mit **Mit Windows starten** aktivierst du den Autostart.

Einstellungen, Protokoll und Über sind Seiten im selben Tray-Menü; mit **Zurück** gelangst du wieder zur Hauptseite. Ein erneuter App-Start öffnet die bereits laufende Instanz.

Statustext und Hinweis beim Darüberfahren erklären den Zustand: Grün bedeutet laufend, Gelb Start oder Neustart und Rot gestoppt oder Fehler. Unter **Über** öffnest du die Projektseite, die Lizenz und die Quellenhinweise.

Mit **Neu starten** startest du den Dienst erneut. Nach einem unerwarteten Ende versucht die App höchstens drei automatische Neustarts. Danach stoppt sie und zeigt den Fehler an. Unter **Einstellungen** kannst du dein Konto erneut verbinden, Listen verwalten und den lokalen Port ändern. Speichere und verlasse die Einstellungen, damit der Dienst fortgesetzt wird.

Die automatischen Neustartversuche warten 2, 5 und 15 Sekunden. Nach einer Minute stabilem Betrieb wird der Versuchszähler zurückgesetzt. Bei Anmelde- und Konfigurationsfehlern stoppt der Dienst, damit du die Ursache beheben kannst; absichtliches Stoppen löst keinen automatischen Neustart aus.

<p>
  <img src="docs/images/tray-menu.png" width="368" alt="Simkl-Bridge-Menü mit Status, Steuerung, Einstellungen und Sprachauswahl">
</p>

*Menüvorschau; der gezeigte Verbindungsstatus ist ein Beispiel.*

## Listen aktualisieren und Daten schützen

Die Bridge prüft bei einem Listenabruf auf Änderungen und verwendet unveränderte Inhalte erneut. Sonarr und Radarr erhalten immer vollständige Listen. API-Fehler werden gemeldet und nicht als leere Listen ausgegeben. [So funktioniert die Aktualisierung](docs/refresh.de.md)

Nur eingerichtete Listen sind über `127.0.0.1` erreichbar. Einstellungen, Listencache und Protokoll liegen in `%LOCALAPPDATA%\SimklBridge`. Anmeldetokens werden mit dem Windows-Benutzerschutz DPAPI geschützt. Aktualisierungen und die normale Deinstallation erhalten diese Daten.

Die Bridge fordert ausschließlich lesenden Zugriff an und erneuert das gespeicherte Zugriffstoken automatisch, wenn die API es verlangt und das Refresh-Token noch gültig ist. Kann die Anmeldung nicht mehr erneuert werden, verbinde das Konto unter **Einstellungen → Anmeldung** erneut. Änderungen werden nicht zu Simkl zurückgeschrieben.

## Fehlerbehebung

Die [Hilfetabelle](docs/setup-simkl.de.md#fehlerbehebung) erklärt Anmeldefehler, leere Listen, PRO-Status und Portkonflikte. Wähle im Menü **Protokoll öffnen**. Das Protokoll enthält Zeitstempel, erkannte Dienst- und Listenfehler sowie Neustartzustände. Wiederholte unveränderte Listenfehler werden nicht bei jedem Abruf erneut protokolliert; es ist kein vollständiger Mitschnitt aller Aktionen. Die Protokollmeldungen sind derzeit auch bei englischer Oberfläche auf Deutsch.

Bleibt das Problem bestehen, [erstelle ein Issue](https://github.com/Veyilla/simkl-bridge-windows/issues/new) mit diesen Angaben:

- App-Version, Windows-Version und betroffenes Zielprogramm: Sonarr oder Radarr.
- Schritte zum Nachstellen, erwartetes Ergebnis und tatsächliche Fehlermeldung.
- Zeitpunkt des Fehlers und passende Zeilen aus `%LOCALAPPDATA%\SimklBridge\bridge.log`.
- Bei Darstellungsfehlern: Screenshot, Bildschirmauflösung und Windows-Skalierung in Prozent.

Prüfe den Auszug vor dem Veröffentlichen und entferne Tokens, persönliche Listenadressen und andere private Angaben. Gibt es keinen passenden Protokolleintrag, erwähne das im Issue und ergänze die Fehlermeldung aus der App oder aus Sonarr/Radarr.

## Selbst bauen

Installiere unter Windows Go **1.26.1** und Inno Setup **7.1.1**. Führe anschließend im Repository-Ordner aus:

```powershell
./scripts/build.ps1 -ISCCExe 'C:/Pfad/zu/ISCC.exe' -ProjectURL 'https://github.com/Veyilla/simkl-bridge-windows'
```

Das Skript führt die Prüfungen aus und erstellt App, Installer, Windows-Archiv und Quellcodearchiv im benachbarten Ordner `Simkl-Bridge-Windows`. Eine zusätzliche Race-Prüfung benötigt unter Windows eine C-Werkzeugkette.

## Lizenz und Quellen

Copyright © 2026 **Veyilla Nightwhisper**. Veröffentlicht unter der [MIT-Lizenz](LICENSE).

Dies ist eine unabhängige App, kein offizielles Simkl-Produkt. Sie verwendet angepasste Bausteine aus [Resurrector](https://github.com/Maki-Daisuke/resurrector/tree/v1.2.0) von [Daisuke (yet another) Maki](https://github.com/Maki-Daisuke), das Simkl-Bild aus Dashboard Icons und die Go-Windows-Bibliothek. Die ursprünglichen Hinweise und Lizenztexte stehen in den [Quellenhinweisen](THIRD_PARTY_NOTICES.de.md) und im [Lizenzordner](licenses).
