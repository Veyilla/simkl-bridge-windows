# Simkl in wenigen Schritten verbinden

Diese Anleitung gehört zur Windows-Vorabversion. Anmeldung, Listen und Einstellungen befinden sich direkt im Tray-Menü.

## Was du brauchst

Ein Simkl-Konto mit PRO oder VIP fuer den Abruf eigener Custom Lists. Die normalen Kategorien wie „Watching“ und „Plan to Watch“ sind etwas anderes als selbst benannte Listen. [Simkl: Custom Lists](https://api.simkl.org/guides/custom-lists)

## Kostenloses PRO ohne bezahltes Abo

Auch das kostenlose PRO aus dem Simkl-Treueprogramm zählt als PRO-Zugang:

1. Nutze angemeldet die Simkl-Webseite oder die Simkl-Lists-App an mindestens **20 Tagen im Monat**.
2. Öffne die [Treueprogramm-Seite](https://simkl.com/vip/loyalty/) und prüfe deinen Fortschritt.
3. Sobald verfügbar, klicke **Claim Free PRO**. Das schaltet PRO für einen Monat frei.
4. Wiederhole die Freischaltung monatlich; sie verlängert sich nicht automatisch.

Dafür ist kein bezahltes Abo nötig. Die Bridge sammelt keine Nutzungstage automatisch. Maßgeblich sind die [aktuellen Simkl-Bedingungen](https://docs.simkl.org/how-to-use-simkl/account-and-billing/special-plans/offers/loyalty-program).

## 1. Eine V2-App bei Simkl anlegen

1. Melde dich bei Simkl an und oeffne [Einstellungen → Developer](https://simkl.com/settings/developer/).
2. Erstelle eine neue App mit **AUTH V2**.
3. Waehle fuer diese Bridge **TV, devices & command line**. Die App verwendet die dortige Anmeldung per Code, auch auf einem Windows-PC.
4. Als Namen kannst du **List Bridge for Simkl** verwenden. Die optionalen Angaben duerfen leer bleiben.
5. Pruefe die Angaben und erstelle die App. Kopiere ihre **Client ID**.

Fuer diesen App-Typ brauchst du weder Client Secret noch Redirect-URL. Eine Registrierung reicht fuer Film-, Serien- und Anime-Listen. [Simkl: Anmeldung einrichten](https://api.simkl.org/api-reference/auth)

## 2. Mit der Bridge verbinden

Unter Anmeldung traegst du die Client ID ein und waehlst „Mit Simkl verbinden“. Oeffne den angebotenen Link in deinem Browser und bestaetige den lesenden Zugriff. Die Bridge uebernimmt den technischen Austausch und speichert ihre Anmeldung lokal. Du musst keine Tokens kopieren oder Befehle eingeben.

Die Client ID bezeichnet die registrierte App. Erst deine Bestaetigung verbindet dein Konto. Access Token und Refresh Token bleiben privat. [Simkl: AUTH V2](https://api.simkl.org/api-reference/auth-v2)

## 3. Listen zuordnen

Waehle die gewuenschte Liste und das Ziel Sonarr oder Radarr. Die Bridge stellt die dazugehoerige lokale Listenadresse bereit, die du im Zielprogramm unter „Benutzerdefinierte Liste“ eintraegst. Alternativ kannst du direkt einen Simkl-Listenlink oder eine Listen-ID eingeben.

Vor dem Speichern pruefst du die Aufraeumoptionen des Zielprogramms: Sie koennen bestimmen, was mit Titeln passiert, die nicht mehr in einer Liste stehen. Die Bridge aendert diese Einstellungen nicht selbst.

## Warum AUTH V1 hier nicht funktioniert

Gemeint ist das alte **Anmeldesystem AUTH V1**, nicht jede bisherige API-Funktion. Die Custom-Lists-Schnittstelle verlangt schon jetzt AUTH V2. Eine alte V1-Client-ID laesst sich nicht umschalten; du registrierst eine neue V2-App und meldest dich damit neu an. [Custom-Lists-Voraussetzungen](https://api.simkl.org/guides/custom-lists), [Umstieg auf V2](https://api.simkl.org/guides/migrating-v1-to-v2)

**Gepruefter Stand: 3. November 2026.** Simkl erwartet die Abschaltung von AUTH V1 ungefaehr im April 2027. Ein verbindlicher genauer Termin steht noch nicht fest. Die verlinkte offizielle Seite ist fuer spaetere Aenderungen massgeblich. [Offizielle Ankuendigung](https://api.simkl.org/api-reference/auth-v1)

## Wenn es nicht klappt

- **Falsche Client ID / invalid_client:** Pruefe, ob sie wirklich zu einer AUTH-V2-App gehoert.
- **V1 funktioniert sonst noch:** Das aendert nichts an der V2-Pflicht fuer Custom Lists.
- **PRO/VIP erforderlich:** Die API prueft das Konto, mit dem du den Zugriff bestaetigt hast.
- **Liste leer:** Pruefe ihre Inhalte und ob es wirklich die gewuenschte eigene Liste ist.

Mit Zurück gelangst du aus den Unterseiten wieder ins Hauptmenü. Die Anmeldung im Browser bleibt davon getrennt. Persönliche Daten sind im Quellcode und im Installationspaket nicht enthalten.

