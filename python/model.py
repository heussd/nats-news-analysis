from datetime import datetime
import json
from typing import List, Union
from numpyencoder import NumpyEncoder
from pydantic import BaseModel


class SearchDoc(BaseModel):
    title: str = ""
    excerpt: str = ""
    author: str = ""
    language: str = ""
    url: Union[str | None] = None
    content: str = ""
    date: Union[str | None] = None
    vector: List[float] = []
