import os


def read_docker_secret(secret):
    secret_path = "/run/secrets/" + secret

    if not os.path.isfile(secret_path):
        print("Docker secret not found: " + str(secret_path) + " - trying to fall back to local file")
        secret_path = "./" + secret + ".txt"

    file = open(secret_path, 'r')
    secret_string = file.read()

    assert secret_string is not None
    return secret_string.strip()


POCKET_ACCESS_TOKEN = read_docker_secret("POCKET_ACCESS_TOKEN")
POCKET_CONSUMER_KEY = read_docker_secret("POCKET_CONSUMER_KEY")

NATS_SUBJECT = os.getenv('NATS_SUBJECT', "4pocket")
NATS_SERVER = os.getenv('NATS_SERVER', "localhost:4222")
NATS_QUEUE = os.getenv('NATS_QUEUE', "NATS-RPLY-22")


print("-------------- Config variables -----------------")

for name, value in globals().copy().items():
    if "__" in name or "os" == name or "read_docker_secret" == name:
        continue
    print(name, "=", value)

print("-------------------------------------------------")


