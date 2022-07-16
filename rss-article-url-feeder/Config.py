import os

NATS_SUBJECT = os.getenv('NATS_SUBJECT', "article-url")
NATS_QUEUE = os.getenv('NATS_QUEUE', "article-urls")
NATS_SERVER = os.getenv('NATS_SERVER', "localhost:4222")
RELOAD_EVERY_S = os.getenv('RELOAD_EVERY_S', 3*60*60)
URLS = os.getenv('URLS', "urls.txt")


print("-------------- Config variables -----------------")

for name, value in globals().copy().items():
    if "__" in name or "os" == name or "read_docker_secret" == name:
        continue
    print(name, "=", value)

print("-------------------------------------------------")


