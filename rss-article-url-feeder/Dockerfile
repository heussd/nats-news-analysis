FROM  python:slim-buster
RUN   apt-get update && \
      apt-get -y install --no-install-recommends \
      build-essential libxml2-dev libxslt1-dev zlib1g-dev python3-pip python-dev\
      && rm -rf /var/lib/apt/lists/*
WORKDIR /
COPY    requirements.txt ./
RUN   pip install -r requirements.txt
COPY  *.py ./

CMD [ "python3", "-u", "main.py"]
