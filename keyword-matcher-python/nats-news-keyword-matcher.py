import asyncio

from nats.aio.msg import Msg

import Config
import Keywords
import RSSFullText
from NATS import NATS

nats = NATS()


async def callback(message: Msg):
    print(message)
    news = RSSFullText.retrieve_full_text(message.data.decode())
    match = Keywords.match(news)

    if match is not None:
        print("Found relevant match", match)
        await nats.publish(match.news.url)

    await message.ack()


async def listen():
    await nats.connect()
    await nats.subscribe(callback=callback)

    while True:
        await asyncio.sleep(Config.RELOAD_EVERY_S)


if __name__ == '__main__':
    print("Starting NATS-News-Keyword-Matcher...")
    asyncio.run(listen())


