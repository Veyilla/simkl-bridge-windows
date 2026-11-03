# So funktioniert die Aktualisierung

<a href="refresh.md"><img src="https://github.githubassets.com/images/icons/emoji/unicode/1f1ec-1f1e7.png" width="20" height="20" alt=""> English</a> · <strong><a href="refresh.de.md"><img src="https://github.githubassets.com/images/icons/emoji/unicode/1f1e9-1f1ea.png" width="20" height="20" alt=""> Deutsch</a></strong>

[README](../README.de.md) · [Einrichtungsanleitung](setup-simkl.de.md)

## Im Alltag

Sonarr oder Radarr ruft eine Importliste nach seinem eigenen Zeitplan ab. Beim ersten Abruf wird die entsprechende Simkl-Custom-List vollständig gelesen. Später verwendet die Bridge den gespeicherten Inhalt erneut, wenn keine relevante Änderung erkannt wurde. Es gibt keine ständige Hintergrundabfrage bei Simkl.

Die Bridge liefert immer eine **vollständige Liste**, nicht nur neu hinzugefügte Titel. Dieses Format benötigen die Zielprogramme. Die Änderungsprüfung verringert unnötige Downloads von Simkl; sie ändert nicht das an Sonarr oder Radarr übergebene Listenformat.

## Wann erneut heruntergeladen wird

| Liste | So wird sie geprüft |
|---|---|
| Deine eigene normale Custom List | Aktivitäten der Custom Lists prüfen, bei Bedarf die Listenübersicht des Besitzers erneuern und den Änderungszeitstempel vergleichen. Erst danach wird gegebenenfalls der Inhalt neu geladen. |
| Liste eines anderen Nutzers | Die Übersicht des Besitzers prüfen. Änderungen an seiner Liste verändern nicht den Aktivitätszeitstempel deines Kontos. |
| Automatische Liste | Bei Bedarf prüfen, einschließlich einer täglichen Inhaltsaktualisierung bei unveränderten Metadaten. Fehlende Metadaten können einen weiteren Inhaltsabruf erfordern. |

Ein kurzer lokaler Zwischenspeicher vermeidet wiederholte Prüfungen innerhalb derselben Minute. Gleichzeitige Anfragen teilen sich den Inhaltsdownload; der gespeicherte Cache wird auch nach einem Neustart verwendet. Er ist an das erkannte Simkl-Konto gebunden.

## Fehler und entfernte Listen

Fehlende Listen, fehlende Pflichtzähler, widersprüchliche Medienarten und unvollständige API-Antworten führen zu einem Fehler. Sie werden nicht als erfolgreiche leere Liste ausgegeben. Erst ein vollständiger erfolgreicher Download ersetzt den gespeicherten Stand. Eine tatsächlich leere Simkl-Liste wird weiterhin leer zurückgegeben.

## Custom Lists und eingebaute Watchlists

Custom Lists sind etwas anderes als die eingebauten Statuslisten von Simkl. Die Route `/sync/all-items?date_from=` liefert Änderungen dieser Statuslisten; sie ist keine zugesagte Änderungsschnittstelle für einzelne Einträge einer Custom List. Eine geänderte Custom List muss gegebenenfalls vollständig neu gelesen werden.

Lokale Statusmeldungen zwischen App und Hintergrunddienst verursachen keine zusätzlichen Simkl-API-Aufrufe.

[Offizieller Sync-Leitfaden](https://api.simkl.org/guides/sync) · [Offizieller Custom-Lists-Leitfaden](https://api.simkl.org/guides/custom-lists)
