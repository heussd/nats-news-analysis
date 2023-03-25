import json

import nats
from nats.js.api import StreamConfig, RetentionPolicy

import Config
from model.NatsOutput import NatsOutput


class NATS:

    HEADER_MESSAGE_ID = "Nats-Msg-Id"
    DUPLICATION_WINDOW_1_MONTH = 3600000000000 * 24 * 30

    def __init__(self):
        self.nc = None
        self.js = None

    async def connect(self):
        self.nc = await nats.connect("nats://" + Config.NATS_SERVER)
        self.js = self.nc.jetstream()

        # This stopped working with nats-py 2.2.0 :/
        # await self.js.add_stream(
        #     name=Config.NATS_QUEUE_OUTPUT,
        #     subjects=[Config.NATS_SUBJECT_OUTPUT],
        #     config=StreamConfig(
        #         retention=RetentionPolicy.WORK_QUEUE,
        #         duplicate_window=NATS.DUPLICATION_WINDOW_1_MONTH
        #     )
        # )

    async def publish(self, message: NatsOutput):
        ack = await self.js.publish(
            Config.NATS_SUBJECT_OUTPUT,
            json.dumps({
                "Url": message.Url,
                "RegExId": message.RegExId
            }).encode(),
            headers={
                NATS.HEADER_MESSAGE_ID: "news-keyword-matcher" + message.Url
            }
        )
        print(ack)

    async def subscribe(self, callback):
        await self.js.subscribe(Config.NATS_SUBJECT_INPUT,
                                Config.NATS_QUEUE_INPUT,
                                cb=callback)

    async def close(self):
        await self.nc.close()
