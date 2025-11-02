import os
from dotenv import load_dotenv
from fastembed import TextEmbedding

load_dotenv()

NATS_URL = os.getenv("NATS_URL", "nats://localhost:4222")
NATS_STREAM_NAME = os.getenv("NATS_STREAM_NAME", "news")
NATS_CONSUMER_NAME = os.getenv("NATS_CONSUMER_NAME", "indexer")

AI_SEARCH_ENDPOINT = os.getenv("AI_SEARCH_ENDPOINT", None)
AI_SEARCH_API_VERSION = os.getenv("AI_SEARCH_API_VERSION", None)
AI_SEARCH_API_KEY = os.getenv("AI_SEARCH_API_KEY", None)

EMBEDDING_MODEL_NAME = os.getenv("EMBEDDING_MODEL_NAME", None)

if not all([AI_SEARCH_ENDPOINT, AI_SEARCH_API_VERSION, AI_SEARCH_API_KEY, EMBEDDING_MODEL_NAME]):
    raise ValueError("Missing required environment variable")


model = TextEmbedding(
    model_name=EMBEDDING_MODEL_NAME,
    cache_folder="/.cache/",
    local_files_only=True
)
