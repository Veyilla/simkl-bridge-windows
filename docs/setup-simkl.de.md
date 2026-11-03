# Simkl Bridge einrichten

<a href="setup-simkl.en.md"><img src="https://github.githubassets.com/images/icons/emoji/unicode/1f1ec-1f1e7.png" width="20" height="20" alt=""> English</a> · <strong><a href="setup-simkl.de.md"><img src="https://github.githubassets.com/images/icons/emoji/unicode/1f1e9-1f1ea.png" width="20" height="20" alt=""> Deutsch</a></strong>

[README](../README.de.md) · [Installation](installation.de.md) · [Fehlerbehebung](#fehlerbehebung)

Führe diese Schritte einmal für deine erste Liste aus. Unter **Einstellungen** kannst du später weitere Listen hinzufügen oder den Port ändern. Du hast bereits eine gültige AUTH-V2-Client-ID? Dann beginne bei [Konto verbinden](#2-konto-verbinden).

Die App startet unabhängig von der Installationssprache zunächst auf Deutsch. Wenn du Englisch bevorzugst, wähle **Zurück**, dann unten im Hauptmenü **English** und öffne die Einstellungen erneut. Die folgenden Schritte verwenden die deutschen Beschriftungen.

## Bevor du anfängst

Sonarr und/oder Radarr müssen auf demselben Windows-Rechner installiert sein. Außerdem brauchst du ein Simkl-Konto mit **PRO oder VIP** und eine Custom List mit den gewünschten Titeln. Kostenloses Loyalty-PRO zählt ebenfalls.

Eine **Custom List** ist eine Sammlung, die du auf Simkl selbst benennst. Sie ist etwas anderes als die eingebauten Kategorien **Watching** und **Plan to Watch**. Ein Titel in einer solchen Kategorie erscheint nicht automatisch in deiner Custom List. [Simkls Erklärung zu Custom Lists](https://api.simkl.org/guides/custom-lists)

## Kostenloses PRO ohne bezahltes Abo

1. Nutze angemeldet die Simkl-Webseite oder die Simkl-Lists-App an mindestens **20 Tagen im Monat**.
2. Öffne die [Treueprogramm-Seite](https://simkl.com/vip/loyalty/) und prüfe deinen Fortschritt.
3. Sobald verfügbar, wähle **Claim Free PRO**. Damit wird PRO für einen Monat freigeschaltet.
4. Wiederhole die Freischaltung monatlich; sie verlängert sich nicht automatisch.

Für diesen Weg ist kein bezahltes Abo nötig. Die Bridge sammelt keine Nutzungstage automatisch. [Aktuelle Simkl-Bedingungen](https://docs.simkl.org/how-to-use-simkl/account-and-billing/special-plans/offers/loyalty-program)

## 1. Eine Simkl-App mit AUTH V2 registrieren

Die Registrierung gibt deiner Bridge eine **Client-ID**. Sie erstellt kein weiteres Simkl-Konto.

1. Melde dich bei Simkl an und öffne [Einstellungen → Developer](https://simkl.com/settings/developer/).
2. Erstelle eine neue App mit **AUTH V2**.
3. Wähle **TV, devices & command line**. Diese Option gilt auch für Skripte und Windows-Apps mit Anmeldung per Code.
4. Gib beispielsweise **List Bridge for Simkl** als Namen ein. Optionale Felder für Beschreibung, Homepage und Kontakt dürfen leer bleiben.
5. Prüfe die Angaben und erstelle die App.
6. Kopiere die **Client ID**. Du fügst sie im nächsten Schritt in die Bridge ein.

Eine Registrierung reicht für Film-, Serien- und Anime-Listen. Für diesen App-Typ brauchst du **kein Client Secret und keine Redirect-URL**. Eine nicht verifizierte eigene App kann verwendet werden. [Simkl-Dokumentation zur Anmeldung](https://api.simkl.org/api-reference/auth)

## 2. Konto verbinden

1. Öffne Simkl Bridge. Falls das Hauptmenü angezeigt wird, wähle **Einstellungen → Anmeldung**.
2. Füge die Client-ID in **Client-ID deiner AUTH-V2-App** ein.
3. Wähle **Mit Simkl verbinden**. Im Menü erscheint ein kurzer Anmeldecode.
4. Wähle **Browser öffnen**. Prüfe, ob du dort im gewünschten Simkl-Konto angemeldet bist. Gib den Code ein, falls du dazu aufgefordert wirst, und bestätige den lesenden Zugriff.
5. Warte, bis die Bridge die erfolgreiche Anmeldung meldet. Lass die Einstellungsseite währenddessen geöffnet.

Die Bridge speichert die Anmeldung geschützt in deinem Windows-Benutzerkonto. Du musst keine Access-Tokens oder Refresh-Tokens kopieren. Die Client-ID bezeichnet die App; deine Bestätigung im Browser verbindet dein Konto.

## 3. Liste auswählen und Adresse kopieren

1. Öffne **Listen** und wähle **Meine Listen laden**.
2. Markiere eine Liste. Wähle **Auswahl → Sonarr** für Serien oder **Auswahl → Radarr** für Filme. Alternativ kannst du einen Simkl-Listenlink oder eine Listen-ID einfügen und den passenden Knopf **Link → …** verwenden.
3. Markiere den hinzugefügten Eintrag in der unteren Liste und wähle **Adresse kopieren**.
4. Wähle **Speichern und zurück**. Falls die Bridge gestoppt ist, wähle **Bridge starten**.

Die Adresse liegt nun in der Zwischenablage. Du verwendest sie im nächsten Abschnitt. Speichere und verlasse die Einstellungsseite vor dem Test in Sonarr oder Radarr: Während der Bearbeitung pausiert der Bridge-Dienst. Wiederhole die Schritte für jede Liste und jedes Ziel.

## 4. Importliste in Sonarr oder Radarr hinzufügen

Im jeweiligen Zielprogramm:

1. Öffne **Einstellungen → Importlisten** und wähle **+**.
2. Wähle **Benutzerdefinierte Liste** beziehungsweise **Custom List**, falls deine Oberfläche englisch ist.
3. Gib der Importliste einen erkennbaren Namen und füge die kopierte Bridge-Adresse in das URL-Feld ein.
4. Wähle den gewünschten Root-Ordner und das Qualitätsprofil. Entscheide, ob neue Titel automatisch überwacht und gesucht werden sollen.
5. Prüfe vor dem ersten Import die Aufräumoptionen der Mediathek. Sie bestimmen den Umgang mit Titeln, die aus Listen verschwinden; die Bridge ändert diese Einstellung nicht.
6. Wähle **Prüfen**, dann **Speichern**. Starte die Listenaktualisierung im Zielprogramm oder warte auf dessen geplanten Abruf.

Sonarr und Radarr haben eigene Importintervalle. Die Bridge liefert die Listen; sie bestimmt weder die Downloadqualität noch startet sie Importe über einen Administratorschlüssel. Beide Zielprogramme müssen die Bridge auf demselben Rechner erreichen können.

## 5. Mit Windows starten

Klicke mit der rechten Maustaste auf das S-Symbol neben der Uhr und aktiviere **Mit Windows starten**. Mit **Bridge stoppen** stoppst du nur den Dienst; mit **Beenden** schließt du die gesamte App. Beim nächsten Start werden Anmeldung und Listen wiederverwendet.

## Warum AUTH V1 hier nicht funktioniert

Custom Lists verlangen **AUTH V2**. Eine V1-Client-ID lässt sich nicht umwandeln: Registriere eine separate V2-App und verbinde dich mit ihrer Client-ID neu. [Anleitung für den Umstieg](https://api.simkl.org/guides/migrating-v1-to-v2)

Bei der dokumentierten Prüfung vom **3. November 2026** nannte Simkl **ungefähr April 2027** für die Abschaltung von AUTH V1, ohne festen Termin. Die [offizielle AUTH-V1-Ankündigung](https://api.simkl.org/api-reference/auth-v1) ist für spätere Änderungen maßgeblich.

## Fehlerbehebung

| Was du siehst | Was du prüfen kannst |
|---|---|
| `invalid_client` oder ein Hinweis auf die Client-ID | Verwende die vollständige Client-ID einer **AUTH-V2**-Registrierung. |
| Anmeldung erforderlich | Öffne **Einstellungen → Anmeldung** und verbinde dich erneut. Bestätige im Browser das gewünschte Konto. |
| PRO/VIP erforderlich | Prüfe den aktuellen Status des bestätigten Kontos. Kostenloses PRO muss manuell freigeschaltet und verlängert werden. |
| Leere Liste | Öffne die Custom List auf Simkl und prüfe ihren Inhalt. Eine eingebaute Statuskategorie ist eine andere Liste. |
| Prüfen schlägt in Sonarr/Radarr fehl | Verlasse die Bridge-Einstellungen mit **Speichern und zurück** und starte die Bridge. Kopiere die Adresse bei geändertem Port erneut. |
| Port bereits belegt | Eine andere App oder ältere Bridge kann Port 8787 verwenden. Wähle unter **Einstellungen → Allgemein** einen unbenutzten Port, speichere und kopiere die neuen Adressen. |
| Falsches Ziel oder fehlende Kennungen | Ordne Serien Sonarr und Filme Radarr zu. Prüfe das **Protokoll** auf fehlende TVDB-/TMDb-Kennungen. |

Unter [Installation](installation.de.md) stehen Hinweise zu Updates und Deinstallation. [So funktioniert die Aktualisierung](refresh.de.md) erklärt den Listencache.

Lösen diese Prüfungen das Problem nicht, folge den Hinweisen zum [Melden eines Fehlers](../README.de.md#fehlerbehebung). Ergänze im Issue einen passenden Protokollauszug und die Schritte zum Nachstellen; entferne vorher private Angaben. Das Protokoll verwendet unabhängig von der Oberflächensprache derzeit deutsche Meldungen.
