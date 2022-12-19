import asyncio
import json
import time

from nats.aio.msg import Msg

import Pocket
from NATS import NATS


async def add_to_pocket(msg: Msg):
    url, matchingtext = "", "UNDEFINED"
    try:
        payload = json.loads(msg.data.decode())
        matchingtext = payload["MatchingText"]
        url = payload["Url"]
    except json.decoder.JSONDecodeError:
        url = msg.data.decode()

    print(f"Received a message on '{msg.subject} {msg.reply}': {url} {matchingtext}")
    Pocket.add_to_pocket(url, matchingtext)
    await msg.ack()


async def listen():
    nats = NATS()
    await nats.connect()
    await nats.subscribe(callback=add_to_pocket)

    while True:
        await asyncio.sleep(10)


if __name__ == '__main__':
    print("Allow NATS-Server to come up...")
    time.sleep(10)
    print("Starting NATS-Pocket-Integration...")
    asyncio.run(listen())
