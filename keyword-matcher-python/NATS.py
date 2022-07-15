import nats
from nats.js.api import StreamConfig, RetentionPolicy

import Config


class NATS:

    HEADER_MESSAGE_ID = "Nats-Msg-Id"
    DUPLICATION_WINDOW_1_MONTH = 3600000000000 * 24 * 30

    def __init__(self):
        self.nc = None
        self.js = None

    async def connect(self):
        self.nc = await nats.connect("nats://" + Config.NATS_SERVER)
        self.js = self.nc.jetstream()

        await self.js.add_stream(
            name=Config.NATS_QUEUE_OUTPUT,
            subjects=[Config.NATS_SUBJECT_OUTPUT],
            config=StreamConfig(
                retention=RetentionPolicy.WORK_QUEUE,
                duplicate_window=NATS.DUPLICATION_WINDOW_1_MONTH
            )
        )

    async def publish(self, message):
        ack = await self.js.publish(
            Config.NATS_SUBJECT_OUTPUT,
            message.encode(),
            headers={
                NATS.HEADER_MESSAGE_ID: "news-keyword-matcher" + message
            }
        )
        print(ack)

    async def subscribe(self, callback):
        await self.js.subscribe(Config.NATS_SUBJECT_INPUT,
                                Config.NATS_QUEUE_INPUT,
                                cb=callback)

    async def close(self):
        await self.nc.close()
