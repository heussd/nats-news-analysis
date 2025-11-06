FROM python:3.12-slim AS base
COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

RUN apt update && apt install -y --no-install-recommends \
    npm \
    git \
    make