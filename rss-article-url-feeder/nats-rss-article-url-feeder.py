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

                print("Publishing", link)
                await nats.publish(link)

    print("Waiting", Config.RELOAD_EVERY_S, "seconds to reload...")
    await asyncio.sleep(Config.RELOAD_EVERY_S)


if __name__ == '__main__':
    print("Starting NATS-RSS-Article-URL-Feeder...")

    while True:
        asyncio.run(feed_urls())


