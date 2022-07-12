import asyncio

import nats
from nats.js.api import StreamConfig, RetentionPolicy

import Config
import NATS
import RSS


async def feed_urls():
    nc = await nats.connect("nats://" + Config.NATS_SERVER)
    js = nc.jetstream()

    await js.add_stream(
        name=Config.NATS_SUBJECT+"s",
        subjects=[Config.NATS_SUBJECT],
        config=StreamConfig(
            retention=RetentionPolicy.WORK_QUEUE,
            duplicate_window=NATS.DUPLICATION_WINDOW_1_MONTH
        )
    )

    with open(Config.URLS) as f:
        feedurls = f.readlines()

        for feedurl in feedurls:
            feedurl = feedurl.strip()

            for link in RSS.retrieve_article_links(feedurl):
                link = link.strip()

                print("Publishing", link)

                ack = await js.publish(
                    Config.NATS_SUBJECT,
                    link.encode(),
                    headers={
                        NATS.HEADER_MESSAGE_ID: link
                    }
                )
                print(ack)

    print("Waiting", Config.RELOAD_EVERY_S, "seconds to reload...")
    await asyncio.sleep(Config.RELOAD_EVERY_S)


if __name__ == '__main__':
    print("Starting NATS-RSS-Article-URL-Feeder...")

    while True:
        asyncio.run(feed_urls())


