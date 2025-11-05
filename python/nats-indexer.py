import json
import nats
import asyncio

from model import SearchDoc
from index import prepare, add
import config


async def process(msg):
    print(f"Received a message")

    json_data = json.loads(msg.data.decode())
    prepared_search_docs = prepare(SearchDoc(
        title=json_data["title"],
        author=json_data["author"],
        content=json_data["content"],
        excerpt=json_data["excerpt"],
        date=json_data["date"],
        language=json_data["language"],
        url=json_data["url"],
    ))

    add(prepared_search_docs)
    await msg.ack()


async def run():
    nc = await nats.connect(config.NATS_URL)
    js = nc.jetstream()

    await js.subscribe(
        config.NATS_STREAM_NAME,
        config=nats.js.api.ConsumerConfig(durable_name=config.NATS_CONSUMER_NAME),
        cb=process,
    )


if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    loop.run_until_complete(run())
    loop.run_forever()
    loop.close()
