FROM golang:1.13.5-buster

WORKDIR /zoracle


RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    curl docker.io sqlite3 \
    && rm -rf /var/lib/apt/lists/*

COPY . /zoracle

RUN make install

CMD bandd start --rpc.laddr tcp://0.0.0.0:26657 & go run cmd/provider/main.go
