from pathlib import Path
from typing import Any

from fastapi import FastAPI, HTTPException, Query
from fastapi.responses import HTMLResponse
from pydantic import BaseModel, Field
import requests
import uvicorn

import ai_search
from index import prepare
from model import SearchDoc

app = FastAPI(
    title="News Search REST API",
    description="REST wrapper for indexing and semantic search.",
    version="1.0.0",
)

BASE_DIR = Path(__file__).resolve().parent
UI_PATH = BASE_DIR / "static" / "index.html"


class SearchRequest(BaseModel):
    query: str = Field(min_length=1)
    top: int = Field(default=10, ge=1, le=200)
    baseUrl: str | None = None


@app.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}


@app.get("/", response_class=HTMLResponse)
def home() -> str:
    try:
        return UI_PATH.read_text(encoding="utf-8")
    except OSError as exc:
        raise HTTPException(status_code=500, detail=f"Unable to load UI file: {UI_PATH}") from exc


@app.get("/search")
def search_get(
    q: str = Query(min_length=1),
    n: int = Query(default=10, ge=1, le=200),
    u: str | None = None,
) -> dict[str, Any]:
    try:
        return ai_search.search(
            query=q,
            top=n,
            baseUrl=u,
        )
    except requests.RequestException as exc:
        raise HTTPException(status_code=502, detail=f"Search backend error: {exc}") from exc


@app.post("/search")
def search_post(payload: SearchRequest) -> dict[str, Any]:
    try:
        return ai_search.search(
            query=payload.query,
            top=payload.top,
            baseUrl=payload.baseUrl,
        )
    except requests.RequestException as exc:
        raise HTTPException(status_code=502, detail=f"Search backend error: {exc}") from exc



if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
