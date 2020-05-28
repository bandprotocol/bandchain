FROM golang:1.13.5-buster

WORKDIR /oracle

COPY . /oracle

RUN make install

CMD bandd start --rpc.laddr tcp://0.0.0.0:26657
