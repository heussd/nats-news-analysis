import asyncio
import json
from timeit import default_timer

from html_sanitizer import Sanitizer
from nats.aio.msg import Msg

import RSSFullText
from Keywords import Keywords
from NATS import NATS
from model.NatsOutput import NatsOutput

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

    j = json.dumps({
        "service": "keyword-matcher-python",
        "match": match,
        "regex-id": id,
        "domain": fulltext.domain,
        "fulltext-length": len(text),
        "retrieval-duration-ms": int(1000*(retrieval_stop-retrieval_start)),
        "keyword-matching-duration-ms": int(1000*(matching_stop-matching_start)),
        "message": "Analysis complete"
    })
    print(j)


async def run():
    await nats.connect()
    await nats.subscribe(callback=callback)


if __name__ == '__main__':
    print("Starting NATS-News-Keyword-Matcher...")
    loop = asyncio.get_event_loop()
    loop.run_until_complete(run())
    loop.run_forever()
    loop.close()
