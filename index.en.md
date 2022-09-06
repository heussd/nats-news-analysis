---
Title: Personal News Analysis
summary: Automatically find relevant news from the Web.
---

Systematically retrieves online news articles, enriches them, scans them for keywords and sends hits to [GetPocket.com](https://getpocket.com/). All analysis components are loosely-coupled with [NATS.io](https://nats.io/) work queues, which also allows scaling single-core-CPU-intensive components easily.


![](architecture.drawio.svg)


## Involved services

All services are orchestrated and scaled with [docker-compose](docker-compose.yml).

### Own services

- [ghcr.io/heussd/nats-rss-article-url-feeder](https://github.com/heussd/nats-rss-article-url-feeder/pkgs/container/nats-rss-article-url-feeder) - Feeds news articles from RSS feeds.
- [ghcr.io/heussd/nats-news-keyword-matcher](https://github.com/heussd/nats-news-keyword-matcher/pkgs/container/nats-news-keyword-matcher) - Matches against keywords list.
- [ghcr.io/heussd/nats-pocket-integration](https://github.com/heussd/nats-pocket-integration/pkgs/container/nats-pocket-integration) - Feeds matches into getpocket.com.
- [docker.io/heussd/fivefilters-full-text-rss](https://hub.docker.com/r/heussd/fivefilters-full-text-rss) - Retrieves full text of web pages.


### Third party services

- [docker.io/nats](https://hub.docker.com/_/nats) - Event queue, key-value store and deduplication.
- [getpocket.com API](https://getpocket.com/developer/) - "Read it later" online service.

## Message queue for scaling

Instead of blocking the application with a single core keyword matching operation, or even trying to build a complex multi-threading keyword matching, we are using the `scale` option of docker compose to run multiple single-core keyword matching components in parallel, wired together with the message queue. This allows us to keep individual components super straight-forward and easy to maintain.


### Keyword matching containers, scaled up

![](docker-container.png)


### One core per keyword matching

![](cpu-cores.png)
