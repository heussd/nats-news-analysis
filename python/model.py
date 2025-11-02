from datetime import datetime
import json
from typing import List, Union
from numpyencoder import NumpyEncoder
from pydantic import BaseModel
from uuid import uuid4
from pydantic import Field


class SearchDoc(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid4()))
    title: str = ""
    excerpt: str = ""
    author: str = ""
    language: str = ""
    url: Union[str | None] = None
    content: str = ""
    date: Union[str | None] = None
    vector: List[float] = []
