# Listen sparsam aktualisieren

Beim ersten Abruf wird eine eingerichtete Liste vollständig gelesen. Danach prüft die Bridge bei Bedarf zuerst `/sync/activities`. Nur ein geänderter Custom-Lists-Zeitstempel führt zu einer neuen Listenübersicht. Anschließend entscheidet `updated_at`, welche Inhalte erneut benötigt werden. Unveränderte Inhalte werden aus dem lokalen Cache geliefert. Parallele Abrufe teilen denselben Inhaltsdownload.

Eine gelöschte oder nicht mehr erreichbare Liste wird durch einen Abruf nach ihrer ID geprüft; ein alter Inhalt wird dann nicht als unverändert weitergereicht. Fehlende Pflichtzähler oder widersprüchliche Medienarten führen zu einer Fehlermeldung. Erst ein vollständiger erfolgreicher Abruf ersetzt den gespeicherten Stand.

Automatische Listen werden bei Benutzung höchstens einmal täglich neu geladen, solange ihre Metadaten unverändert bleiben. Fremde Listen werden anhand der Übersicht ihres Besitzers geprüft; ihr Besitzer verändert nicht den Aktivitätszeitstempel des angemeldeten Kontos. Fehlende Zeitstempel führen zu einer erneuten Prüfung statt zu einem dauerhaft unveränderten Cache.

Sonarr und Radarr erhalten immer vollständige Listen. Für Custom Lists gibt es keine zugesagte Änderungsschnittstelle für einzelne hinzugefügte Titel; eine geänderte Liste muss gegebenenfalls vollständig gelesen werden. Die Watchlist-Route `/sync/all-items?date_from=` ist für die normalen Statuslisten gedacht und wird nicht mit Custom Lists verwechselt.

Es gibt keine ständige Hintergrundabfrage der Simkl-API. Der Cache überlebt einen Neustart und ist an das erkannte Konto gebunden. Statusmeldungen zwischen den beiden lokalen Programmteilen benötigen keine zusätzliche Simkl-Anfrage.

[Offizieller Sync-Leitfaden](https://api.simkl.org/guides/sync) · [Offizieller Custom-Lists-Leitfaden](https://api.simkl.org/guides/custom-lists)

## English

Configured lists are read completely on first use. Later requests check Custom List activity and metadata before downloading changed contents. Unchanged data is cached across restarts, while concurrent requests share downloads. Missing lists, missing counters and inconsistent media types never become successful empty results.

Auto lists refresh on demand at most daily while metadata remains unchanged. Foreign lists use their owner's index. Sonarr and Radarr always receive complete lists; Custom Lists are separate from built-in watchlist deltas. There is no constant background API polling.
