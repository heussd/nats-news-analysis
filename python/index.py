from datetime import datetime
from typing import List, Union

from model import SearchDoc

import config
from ai_search import add


def prepare(searchDoc: Union[SearchDoc, List[SearchDoc]]):
    if isinstance(searchDoc, SearchDoc):
        documents = [searchDoc]
    else:
        documents = searchDoc

    embeddings_list = list(config.model.embed([doc.content for doc in documents]))

    jsonDocs = [
        {
            "@search.action": "mergeOrUpload",
            "id": searchDoc.id,
            "title": searchDoc.title,
            "excerpt": searchDoc.excerpt,
            "author": searchDoc.author,
            "language": searchDoc.language,
            "url": searchDoc.url or "",
            "baseUrl": "/".join(searchDoc.url.split("/")[:3]) if searchDoc.url else "",
            "content": searchDoc.content,
            "vector": embeddings_list[documents.index(searchDoc)],
            "date": searchDoc.date or str(datetime.now().astimezone().isoformat()),
        }
        for searchDoc in documents
    ]

    return jsonDocs


if __name__ == "__main__":
    docs = [
        SearchDoc(
            title="Example Document",
            excerpt="This is an example excerpt.",
            author="John Doe",
            language="en",
            content="Gday this is an important Australian message.",
        ),
        SearchDoc(
            title="Example Document",
            excerpt="This is an example excerpt.",
            author="John Doe",
            language="de",
            content="Guten Tag das ist eine wichtige deutsche Nachricht.",
        ),
    ]

    docs = prepare(docs)
    add(docs)
