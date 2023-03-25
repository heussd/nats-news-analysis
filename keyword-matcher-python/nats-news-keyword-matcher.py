import asyncio
from time import process_time
from timeit import default_timer

import logzero
from html_sanitizer import Sanitizer
from nats.aio.msg import Msg

import Config
import RSSFullText
from Keywords import Keywords
from NATS import NATS
from logzero import logger

from model.NatsOutput import NatsOutput

logzero.json()
nats = NATS()

sanitizer = Sanitizer()
fulltext_service = RSSFullText.FullTextRss()

keywords = Keywords()
keywords.init()


def prepare_and_clean_string(fulltext):
    return " ".join([
        fulltext.title,
        fulltext.excerpt,
        sanitizer.sanitize(fulltext.content)
    ])


async def callback(message: Msg):
    url = message.data.decode()
    retrieval_start = default_timer()
    fulltext = fulltext_service.retrieve_full_text(url)
    retrieval_stop = default_timer()

    text = prepare_and_clean_string(fulltext)

    matching_start = default_timer()
    (match, id) = keywords.match(text)
    matching_stop = default_timer()

    if match:
        await nats.publish(NatsOutput(
            Url=url,
            RegExId=id
        ))

    await message.ack()

    logger.info({
        "service": "keyword-matcher-python",
        "match": match,
        "regex-id": id,
        "domain": fulltext.domain,
        "fulltext-length": len(text),
        "retrieval-duration-ms": int((retrieval_stop-retrieval_start)*1000),
        "keyword-matching-duration-ms": int((matching_stop-matching_start)*1000),
        "msg": "Analysis complete"
    })

async def listen():
    await nats.connect()
    await nats.subscribe(callback=callback)

    while True:
        await asyncio.sleep(Config.RELOAD_EVERY_S)


if __name__ == '__main__':
    logger.info("Starting NATS-News-Keyword-Matcher...")
    asyncio.run(listen())
