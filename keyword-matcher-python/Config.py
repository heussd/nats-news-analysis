import os

NATS_SUBJECT_INPUT = os.getenv('NATS_SUBJECT_INPUT', "article-url")
NATS_SUBJECT_OUTPUT = os.getenv('NATS_SUBJECT_OUTPUT', "match-url")
NATS_QUEUE_INPUT = os.getenv('NATS_QUEUE_INPUT', "article-urls")
NATS_QUEUE_OUTPUT = os.getenv('NATS_QUEUE_OUTPUT', "match-urls")
NATS_SERVER = os.getenv('NATS_SERVER', "localhost:4222")
FULLTEXTRSS_SERVER = os.getenv('FULLTEXTRSS_SERVER', "http://localhost:80")
KEYWORDS_FILE = os.getenv('KEYWORDS_FILE', "keywords.txt")
RELOAD_EVERY_S = os.getenv('RELOAD_EVERY_S', 10)
URLS = os.getenv('URLS', "urls.txt")


print("-------------- Config variables -----------------")

for name, value in globals().copy().items():
    if "__" in name or "os" == name or "read_docker_secret" == name:
        continue
    print(name, "=", value)

print("-------------------------------------------------")


