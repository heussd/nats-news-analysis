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
        jsonmsgs.append(SearchDoc(
            title=json_data["title"],
            author=json_data["author"],
            content=json_data["content"],
            excerpt=json_data["excerpt"],
            date=json_data["date"],
            language=json_data["language"],
            url=json_data["url"],
        ))
    prepared_search_docs = prepare(jsonmsgs)
    if not add(prepared_search_docs):
        raise Exception("Failed to add documents to AI search index")


async def run():
    nc = await nats.connect(config.NATS_URL)
    js = nc.jetstream()

    psub = await js.pull_subscribe("*",
                                   stream=config.NATS_STREAM_NAME,
                                   durable=config.NATS_CONSUMER_NAME,
                                   config=nats.js.api.ConsumerConfig(
                                       deliver_policy=nats.js.api.DeliverPolicy.NEW,
                                   ))

    while True:
        try:
            msgs = await psub.fetch(10, timeout=5)
            print(f"Fetched {len(msgs)} messages")
            await process(msgs)
            for msg in msgs:
                await msg.ack()
        except asyncio.TimeoutError:
            print("No new messages, waiting 10 secs...")
            await asyncio.sleep(10)
            continue



if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    loop.run_until_complete(run())
    loop.run_forever()
    loop.close()
