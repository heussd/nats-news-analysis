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

- [ghcr.io/heussd/nats-rss-article-url-feeder](https://github.com/heussd/nats-rss-article-url-feeder/pkgs/container/nats-rss-article-url-feeder) - Speist Nachrichtenartikel aus RSS-Feeds ein.
- [ghcr.io/heussd/nats-news-keyword-matcher](https://github.com/heussd/nats-news-keyword-matcher/pkgs/container/nats-news-keyword-matcher) - Gleicht gegen eine Keyword-Liste ab.
- [ghcr.io/heussd/nats-pocket-integration](https://github.com/heussd/nats-pocket-integration/pkgs/container/nats-pocket-integration) - Speist Treffer in getpocket.com ein.
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
