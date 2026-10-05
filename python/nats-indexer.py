import json
from typing import List
import nats
import asyncio

from model import SearchDoc
from index import prepare, add
import config


async def process(msgs: List[nats.aio.msg.Msg]):
    jsonmsgs = []
    for msg in msgs:
        json_data = json.loads(msg.data.decode())
        jsonmsgs.append(SearchDoc(**json_data))
    prepared_search_docs = prepare(jsonmsgs)
    add(prepared_search_docs)


async def run():
    nc = await nats.connect(config.NATS_URL)
    js = nc.jetstream()

    psub = await js.pull_subscribe("news.*",
                                   stream=config.NATS_STREAM_NAME,
                                   durable=config.NATS_CONSUMER_NAME,
                                   config=nats.js.api.ConsumerConfig(
                                       deliver_policy=nats.js.api.DeliverPolicy.NEW,
                                   ))

    while True:
        try:
            msgs = await psub.fetch(10, timeout=5)
            print(f"Fetched {len(msgs)} messages")
            try:
                await process(msgs)
            except Exception as exc:
                print(f"Failed to process batch; acknowledging messages anyway: {exc}")

            for msg in msgs:
                try:
                    await msg.ack()
                except Exception as exc:
                    print(f"Failed to acknowledge message: {exc}")
        except asyncio.TimeoutError:
            print("No new messages, waiting 10 secs...")
            await asyncio.sleep(10)
            continue



if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    loop.run_until_complete(run())
    loop.run_forever()
    loop.close()
