import nats
from nats.js.api import StreamConfig, RetentionPolicy, StorageType

import Config


class NATS:

    HEADER_MESSAGE_ID = "Nats-Msg-Id"
    DUPLICATION_WINDOW_1_MONTH = 3600000000000 * 24 * 30

    def __init__(self):
        self.js = None
        self.kv = None

    async def connect(self):
        nc = await nats.connect("nats://" + Config.NATS_SERVER)
        self.js = nc.jetstream()
        self.kv = await self.js.create_key_value(bucket='article-urls-proposed')

        await self.js.add_stream(
            name=Config.NATS_QUEUE,
            subjects=[Config.NATS_SUBJECT],
            config=StreamConfig(
                retention=RetentionPolicy.WORK_QUEUE,
                storage=StorageType.FILE,
                duplicate_window=NATS.DUPLICATION_WINDOW_1_MONTH
            )
        )

    async def publish(self, message):
        ack = await self.js.publish(
            Config.NATS_SUBJECT,
            message.encode(),
            headers={
                NATS.HEADER_MESSAGE_ID: "rss-article-url-feed-" + message
            }
        )
        print(ack)

    async def has_KV(self, key):
        try:
            await self.kv.get(key)
            return True
        except nats.js.errors.NotFoundError:
            return False

    async def put_KV(self, key, value: str):
        await self.kv.put(key, value.encode())
