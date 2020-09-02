FROM golang:1.13.5-buster as builder

WORKDIR /go-owasm
COPY go-owasm/. /go-owasm

WORKDIR /chain
COPY chain/ /chain

COPY chain/docker-config/run.sh .

RUN make install && make faucet

FROM debian:buster

ENV WORKDIR=/chain
WORKDIR ${WORKDIR}

COPY --from=builder /go-owasm /go-owasm
COPY --from=builder /chain/docker-config ${WORKDIR}/docker-config
COPY --from=builder /chain/scripts ${WORKDIR}/scripts
COPY --from=builder /go/bin/ /usr/local/bin/

RUN cp ${WORKDIR}/docker-config/run.sh ${WORKDIR}/run.sh && chmod +x ${WORKDIR}/run.sh

CMD bandd start --rpc.laddr tcp://0.0.0.0:26657
