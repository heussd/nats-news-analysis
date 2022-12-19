from dataclasses import dataclass

from model.RSSFullTextResponse import RSSFullTextResponse


@dataclass
class Match:
    kw_id: str
    kw_pattern: str
    news: RSSFullTextResponse
