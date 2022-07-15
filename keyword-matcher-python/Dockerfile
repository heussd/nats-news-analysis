FROM  python:slim-buster

#RUN   apk add --no-cache --update libxml2-dev libxslt-dev g++ gcc
RUN   apt-get update && \
      apt-get -y install --no-install-recommends \
      libxml2-dev libxslt-dev g++ gcc zlib1g-dev \
      && rm -rf /var/lib/apt/lists/*

WORKDIR /
COPY    requirements.txt ./
RUN   pip install -r requirements.txt
COPY model ./model/
COPY  *.py ./


CMD [ "python3", "-u", "nats-news-keyword-matcher.py"]
