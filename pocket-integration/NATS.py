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

    async def subscribe(self, callback):
        await self.js.subscribe(Config.NATS_SUBJECT,
                                Config.NATS_QUEUE,
                                cb=callback)

    async def close(self):
        await self.nc.close()
