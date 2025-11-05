# Install uv
FROM python:3.12-slim AS base
COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

# Change the working directory to the `app` directory
WORKDIR /app

# Install dependencies
COPY uv.lock pyproject.toml ./
RUN uv sync --locked --no-install-project --no-editable --compile-bytecode

ENV FASTEMBED_CACHE_PATH=/.cache/
ARG FASTEMBED_CACHE_PATH=/.cache/
ENV HF_HOME=/.cache/
ARG HF_HOME=/.cache/

ENV EMBEDDING_MODEL_NAME="sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2"
ARG EMBEDDING_MODEL_NAME="sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2"
RUN uv run python -c "import os; from fastembed import TextEmbedding; TextEmbedding(os.getenv('EMBEDDING_MODEL_NAME'), cache_folder='/.cache/')"

ENV HF_HUB_OFFLINE=1 \
    TRANSFORMERS_OFFLINE=1 \
    F_DATASETS_OFFLINE=1 \
    ONNX_RUNTIME_OFFLINE=1 \
    ORT_DISABLE_ALL_LOGS=1


FROM base as devcontainer

RUN apt update && apt install -y --no-install-recommends \
    golang \
    git \
    make && \
    rm -rf /var/lib/apt/lists/*

RUN go install github.com/nats-io/natscli/nats@latest

ENV PATH="/root/go/bin:${PATH}"


FROM base as build

COPY *.py ./

# Sync the project
RUN --mount=type=cache,target=/root/.cache/uv \
    uv sync --locked --no-editable --compile-bytecode


ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1


FROM build as search-mcp-server

CMD ["uv", "run", "search-mcp-server.py"]


FROM build as nats-indexer

CMD ["uv", "run", "nats-indexer.py"]

