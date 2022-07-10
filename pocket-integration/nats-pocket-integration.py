import asyncio
import nats
from nats.errors import ConnectionClosedError, TimeoutError, NoServersError

import Config
import Pocket


async def main():
    nc = await nats.connect("nats://" + Config.NATS_SERVER)

    async def add_to_pocket(msg):
        payload = msg.data.decode()
        print(f"Received a message on '{msg.subject} {msg.reply}': {payload}")
        Pocket.add_to_pocket(payload, "JUPP")
        await nc.publish(msg.reply, b'Item added to Pocket')

    await nc.subscribe(Config.NATS_SUBJECT, Config.NATS_QUEUE, add_to_pocket)

    while True:
        await asyncio.sleep(10)


if __name__ == '__main__':
    print("Starting NATS-Pocket-Integration...")
    asyncio.run(main())
