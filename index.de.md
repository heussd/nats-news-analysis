---
Title: Persönliche Nachrichtenanalyse
summary: Findet automatisch relevante Nachrichten aus dem Web.
---

Fragt systematisch Online-Nachrichtenartikel ab, reichert sie an, sucht nach Schlagworten und sendet Treffer zu [GetPocket.com](https://getpocket.com/). Alle Komponenten sind lose mit [NATS.io](https://nats.io/) work queues gekoppelt, was es auch erlaubt, Single-Core-CPU-intensive Komponenten einfach zu skalieren.


![](architecture.drawio.svg)

## Beteiligte Services

Alle Services sind durch [docker-compose](docker-compose.yml) orchestriert und skaliert.

### Eigene Services

<!--PYSPELL-BEGIN-IGNORE-->

- [rss-article-url-feeder](rss-article-url-feeder) - Speist Nachrichtenartikel aus RSS-Feeds ein.
- [keyword-matcher-python](keyword-matcher-python) - Gleicht gegen eine Keyword-Liste ab.
- [keyword-matcher-go](keyword-matcher-go) - Gleicht gegen eine Keyword-Liste ab. (Go Reimplementierung)
- [pocket-integration](pocket-integration) - Speist Treffer in getpocket.com ein.
- [docker.io/heussd/fivefilters-full-text-rss](https://hub.docker.com/r/heussd/fivefilters-full-text-rss) - Ruft den Volltext von Nachrichten ab.


### Drittanbieter Services

- [docker.io/nats](https://hub.docker.com/_/nats) - Event Queue, key-value store und Deduplikation.
- [getpocket.com API](https://getpocket.com/developer/) - Online Service zum "Später lesen".

<!--PYSPELL-END-IGNORE-->

## Message queue zum Skalieren

Anstatt die ganze Anwendung mit einer Single-Core-CPU-intensiven Schlagwortsuche zu blockieren, oder gar eine multithreading Schlagwortsuche zu implementieren, kommt die `scale`-Option von Docker compose zum Einsatz, um eine Single-Core-CPU Schlagwortsuche parallel auszuführen, zusammengehalten von einer Message Queue. Das erlaubt es, einzelne Komponenten und deren Pflege sehr einfach zu halten.


### Schlagwortsuche, skaliert

![](docker-container.png)


### Ein Core pro Schlagwortsuche

![](cpu-cores.png)
