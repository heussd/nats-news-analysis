from typing import Optional

from pydantic import BaseModel


class RSSFullTextResponse(BaseModel):

    title: Optional[str]
    excerpt: Optional[str]
    date: Optional[str]

    author: Optional[str]
    language: Optional[str]
    url: str
    effective_url: str
    domain: str
    word_count: int

    og_url: Optional[str]
    og_title: Optional[str]
    og_description: Optional[str]
    og_image: Optional[str]
    og_type: Optional[str]

    twitter_card: Optional[str]
    twitter_site: Optional[str]
    twitter_creator: Optional[str]
    twitter_image: Optional[str]
    twitter_title: Optional[str]
    twitter_description: Optional[str]

    content: str

