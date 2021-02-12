FROM ubuntu:18.04

# built binaries
COPY build/bandcli bin/bandcli
COPY build/bandd bin/bandd

# go-owasm built dynamic lib
COPY build/libgo_owasm.so lib/libgo_owasm.so

COPY docker-config/validator1/ validator1/
COPY docker-config/validator2/ validator2/
COPY docker-config/validator3/ validator3/
COPY docker-config/validator4/ validator4/

COPY docker-config/query-node/ query-node/

# generated genesis
COPY docker-config/genesis.json genesis.json
COPY docker-config/init.sh run.sh

CMD bandcli rest-server --laddr tcp://0.0.0.0:1317 --node tcp://172.18.0.15:26657 --chain-id bandchain --trust-node
