FROM golang:1.13.5-buster AS build-env

WORKDIR /go-owasm
COPY go-owasm/. /go-owasm

WORKDIR /chain
COPY chain/ /chain

COPY chain/docker-config/run.sh .

RUN make install

CMD bandd start --rpc.laddr tcp://0.0.0.0:26657

# FROM alpine

# COPY chain/docker-config/ .

# RUN addgroup banduser && \
#     adduser -S -G banduser banduser -h "$BAND"

# USER banduser

# # Copy over binaries from the build-env
# COPY --from=build-env /go/bin/bandd /usr/bin/bandd
# COPY --from=build-env /go/bin/bandcli /usr/bin/bandcli
# COPY --from=build-env /go/bin/bandoracled2 /usr/bin/bandoracled2

# CMD ["bandd"]
