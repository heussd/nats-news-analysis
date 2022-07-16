import nats
from nats.js.api import StreamConfig, RetentionPolicy

import Config


class NATS:

    HEADER_MESSAGE_ID = "Nats-Msg-Id"
    DUPLICATION_WINDOW_1_MONTH = 3600000000000 * 24 * 30

    def __init__(self):
        self.js = None

    async def connect(self):
        nc = await nats.connect("nats://" + Config.NATS_SERVER)
        self.js = nc.jetstream()

        await self.js.add_stream(
            name=Config.NATS_QUEUE,
            subjects=[Config.NATS_SUBJECT],
            config=StreamConfig(
                retention=RetentionPolicy.WORK_QUEUE,
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
