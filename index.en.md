---
Title: Personal News Analysis
summary: Automatically find relevant news from the Web.
---

Systematically retrieves online news articles, enriches them, scans them for keywords and sends hits to [GetPocket.com](https://getpocket.com/). All analysis components are loosely-coupled with [NATS.io](https://nats.io/) work queues, which also allows scaling single-core-CPU-intensive components easily.


![](architecture.drawio.svg)


## Message queue for scaling

Instead of blocking the application with a single core keyword matching operation, or even trying to build a complex multi-threading keyword matching, we are using the `scale` option of docker compose to run multiple single-core keyword matching components in parallel, wired together with the message queue. This allows us to keep individual components super straight-forward and easy to maintain.


### Keyword matching containers, scaled up

![](docker-container.png)


### One core per keyword matching

![](cpu-cores.png)
