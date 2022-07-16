import asyncio

from nats.aio.msg import Msg

import Pocket
from NATS import NATS


async def add_to_pocket(msg: Msg):
    payload = msg.data.decode()
    print(f"Received a message on '{msg.subject} {msg.reply}': {payload}")
    Pocket.add_to_pocket(payload, "JUPP")
    await msg.ack()


async def listen():
    nats = NATS()
    await nats.connect()
    await nats.subscribe(callback=add_to_pocket)

    while True:
        await asyncio.sleep(10)


if __name__ == '__main__':
    print("Starting NATS-Pocket-Integration...")
    asyncio.run(listen())
