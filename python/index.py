from datetime import datetime
from typing import List, Union

from model import SearchDoc

import config
from ai_search import add
import base64


def prepare(searchDoc: Union[SearchDoc, List[SearchDoc]]):
    if isinstance(searchDoc, SearchDoc):
        documents = [searchDoc]
    else:
        documents = searchDoc

    embeddings_list = list(config.model.embed([f"{doc.title} {doc.excerpt} {doc.content}" for doc in documents]))

    jsonDocs = []

    for idx, searchDoc in enumerate(documents):
        if not searchDoc.url:
            continue

        jsonDocs.append(
            {
            "@search.action": "mergeOrUpload",
            "id": base64.urlsafe_b64encode(searchDoc.url.encode()).decode().rstrip("="),
            "title": searchDoc.title,
            "excerpt": searchDoc.excerpt,
            "author": searchDoc.author,
            "language": searchDoc.language,
            "url": searchDoc.url,
            "baseUrl": "/".join(searchDoc.url.split("/")[:3]) if searchDoc.url else "",
            "content": searchDoc.content,
            "vector": embeddings_list[idx],
            "date": searchDoc.date or str(datetime.now().astimezone().isoformat()),
            }
        )

    return jsonDocs


if __name__ == "__main__":
    docs = [
        SearchDoc(
            title="Example Document",
            excerpt="This is an example excerpt.",
            author="John Doe",
            language="en",
            content="Gday this is an important Australian message.",
            url="https://example.com/doc1",
        ),
        SearchDoc(
            title="Example Document",
            excerpt="This is an example excerpt.",
            author="John Doe",
            language="de",
            content="Guten Tag das ist eine wichtige deutsche Nachricht.",
            url="https://example.com/doc2",
        ),
    ]

    docs = prepare(docs)
    add(docs)
