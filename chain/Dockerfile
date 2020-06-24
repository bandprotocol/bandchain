FROM golang:1.13.5-buster

WORKDIR /go-owasm
COPY go-owasm/. /go-owasm

WORKDIR /chain
COPY chain/ /chain

COPY chain/docker-config/run.sh .

RUN make install && make faucet

CMD bandd start --rpc.laddr tcp://0.0.0.0:26657
