import asyncio

import Config
import RSS
from NATS import NATS


async def feed_urls():

    nats = NATS()
    await nats.connect()

    with open(Config.URLS) as f:
        feedurls = f.readlines()

        for feedurl in feedurls:
            feedurl = feedurl.strip()

            for link in RSS.retrieve_article_links(feedurl):
                link = link.strip()

                if link == "":
                    continue

                # Workaround for https://github.com/nats-io/nats-server/issues/3272
                # Use KV storage for remembering what we already put on the queue.
                if not await nats.has_KV(link):
                    print("Publishing", link)
                    await nats.publish(link)
                    await nats.put_KV(link, "1")

    print("Waiting", Config.RELOAD_EVERY_S, "seconds to reload...")
    await asyncio.sleep(Config.RELOAD_EVERY_S)


if __name__ == '__main__':
    print("Starting NATS-RSS-Article-URL-Feeder...")

    while True:
        asyncio.run(feed_urls())


