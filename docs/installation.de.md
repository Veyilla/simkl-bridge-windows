# Installation unter Windows

1. Den Installer `Simkl-Bridge-…-Windows-x64-Setup.exe` starten.
2. Deutsch oder Englisch wählen und die Installation für den angemeldeten Benutzer abschließen.
3. Simkl Bridge über das Startmenü öffnen. Die Einrichtung führt zur Simkl-Anmeldung und den Listen.

Es werden weder Docker noch Python oder Node benötigt. Der Installer ersetzt keine vorhandene Python-Bridge und übernimmt keine persönlichen Einstellungen daraus.

Das Setup verwendet einen schwarzen Hintergrund und weiße Schrift, auch im Lizenzfeld. Titelleiste, Schaltflächen und Scrollbalken verwenden den eingebauten Dunkelstil und können dunkelgrau sein. Windows-Sicherheitsabfragen werden von Windows gestaltet. Ein aktiviertes Windows-Kontrastdesign hat Vorrang.

## Aktualisieren

Den neuen Installer ausführen und denselben Installationsordner verwenden. Einstellungen und Anmeldung liegen außerhalb des Programmordners und bleiben erhalten. Die App sollte vorher über **Beenden** geschlossen werden.

## Autostart und Deinstallation

**Mit Windows starten** im Tray-Menü schaltet ausschließlich den Autostart dieser App für den aktuellen Benutzer. Die normale Windows-Deinstallation entfernt die Programmdateien und den zugehörigen eigenen Autostart-Eintrag. Persönliche Einstellungen, Anmeldung und Cache bleiben erhalten; die App enthält keine automatische Löschung dieser Daten.

## Hinweise zur Vorabversion

Das Paket ist nicht digital signiert. Prüfe seine Herkunft und die mitgelieferte Prüfsumme. Die vollständige persönliche Einrichtung und die Darstellung auf allen Bildschirmkonfigurationen müssen vor einer allgemeinen Freigabe praktisch geprüft werden.

Wird „Einstellungen prüfen“ angezeigt, kontrolliere die AUTH-V2-Client-ID, PRO/VIP, Listenzuordnung und den lokalen Port. Port 8787 kann bereits von einer anderen Bridge verwendet werden; wähle dann einen anderen freien Port und kopiere die Listenadressen erneut.
