FROM  python:slim-buster
WORKDIR /
COPY    requirements.txt ./
RUN   pip install -r requirements.txt
COPY  *.py ./

CMD [ "python3", "-u", "nats-pocket-integration.py"]
