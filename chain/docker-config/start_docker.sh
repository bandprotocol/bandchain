#!/bin/bash

DIR=`dirname "$0"`

# remove old genesis
rm -rf ~/.band*

make install

# initial new node
bandd init node-validator --chain-id bandchain

# add data sources to genesis

chmod +x $DIR/add_os_ds.sh
$DIR/add_os_ds.sh

# create acccounts
echo "lock nasty suffer dirt dream fine fall deal curtain plate husband sound tower mom crew crawl guard rack snake before fragile course bacon range" \
    | bandcli keys add validator1 --recover --keyring-backend test

echo "loyal damage diet label ability huge dad dash mom design method busy notable cash vast nerve congress drip chunk cheese blur stem dawn fatigue" \
    | bandcli keys add validator2 --recover --keyring-backend test

echo "whip desk enemy only canal swear help walnut cannon great arm onion oval doctor twice dish comfort team meat junior blind city mask aware" \
    | bandcli keys add validator3 --recover --keyring-backend test

echo "unfair beyond material banner okay genre camera dumb grit balcony permit room intact code degree execute twin flip half salt script cause demand recipe" \
    | bandcli keys add validator4 --recover --keyring-backend test

echo "smile stem oven genius cave resource better lunar nasty moon company ridge brass rather supply used horn three panic put venue analyst leader comic" \
    | bandcli keys add requester --recover --keyring-backend test

# add accounts to genesis
bandd add-genesis-account validator1 10000000000000uband --keyring-backend test
bandd add-genesis-account validator2 10000000000000uband --keyring-backend test
bandd add-genesis-account validator3 10000000000000uband --keyring-backend test
bandd add-genesis-account validator4 10000000000000uband --keyring-backend test
bandd add-genesis-account requester 100000000000000uband --keyring-backend test

# genesis configurations
bandcli config chain-id bandchain
bandcli config output json
bandcli config indent true
bandcli config trust-node true

# create copy of config.toml
cp ~/.bandd/config/config.toml ~/.bandd/config/config.toml.temp
cp -r ~/.bandd/files docker-config/

# modify moniker
sed 's/node-validator/🙎‍♀️Alice \& Co./g' ~/.bandd/config/config.toml.temp > ~/.bandd/config/config.toml

# register initial validators
bandd gentx \
    --amount 100000000uband \
    --node-id 11392b605378063b1c505c0ab123f04bd710d7d7 \
    --pubkey bandvalconspub1addwnpepq06h7wvh5n5pmrejr6t3pyn7ytpwd5c0kmv0wjdfujs847em8dusjl96sxg \
    --name validator1 \
    --details "Alice's Adventures in Wonderland (commonly shortened to Alice in Wonderland) is an 1865 novel written by English author Charles Lutwidge Dodgson under the pseudonym Lewis Carroll." \
    --website "https://www.alice.org/" \
    --ip 172.18.0.11 \
    --keyring-backend test

# modify moniker
sed 's/node-validator/Bobby.fish 🐡/g' ~/.bandd/config/config.toml.temp > ~/.bandd/config/config.toml

bandd gentx \
    --amount 100000000uband \
    --node-id 0851086afcd835d5a6fb0ffbf96fcdf74fec742e \
    --pubkey bandvalconspub1addwnpepqfey4c5ul6m5juz36z0dlk8gyg6jcnyrvxm4werkgkmcerx8fn5g2gj9q6w \
    --name validator2 \
    --details "Fish is best known for his appearances with Ring of Honor (ROH) from 2013 to 2017, where he wrestled as one-half of the tag team reDRagon and held the ROH World Tag Team Championship three times and the ROH World Television Championship once." \
    --website "https://www.wwe.com/superstars/bobby-fish" \
    --ip 172.18.0.12 \
    --keyring-backend test

# modify moniker
sed 's/node-validator/Carol/g' ~/.bandd/config/config.toml.temp > ~/.bandd/config/config.toml

bandd gentx \
    --amount 100000000uband \
    --node-id 7b58b086dd915a79836eb8bfa956aeb9488d13b0 \
    --pubkey bandvalconspub1addwnpepqwj5l74gfj8j77v8st0gh932s3uyu2yys7n50qf6pptjgwnqu2arxkkn82m \
    --name validator3 \
    --details "Carol Susan Jane Danvers is a fictional superhero appearing in American comic books published by Marvel Comics. Created by writer Roy Thomas and artist Gene Colan." \
    --website "https://www.marvel.com/characters/captain-marvel-carol-danvers" \
    --ip 172.18.0.13 \
    --keyring-backend test

# modify moniker
sed 's/node-validator/Eve 🦹🏿‍♂️the evil with a really long moniker name/g' ~/.bandd/config/config.toml.temp > ~/.bandd/config/config.toml

bandd gentx \
    --amount 100000000uband \
    --node-id 63808bd64f2ec19acb2a494c8ce8467c595f6fba \
    --pubkey bandvalconspub1addwnpepq0grwz83v8g4s06fusnq5s4jkzxnhgvx67qr5g7v8tx39ur5m8tk7rg2nxj \
    --name validator4 \
   --details "Evil is an American supernatural drama television series created by Robert King and Michelle King that premiered on September 26, 2019, on CBS. The series is produced by CBS Television Studios and King Size Productions." \
    --website "https://www.imdb.com/title/tt9055008/" \
    --ip 172.18.0.14 \
    --keyring-backend test

# remove temp test
rm -rf ~/.bandd/config/config.toml.temp

# collect genesis transactions
bandd collect-gentxs

# copy genesis to the proper location!
cp ~/.bandd/config/genesis.json $DIR/genesis.json

cd ..

docker-compose up -d --build

sleep 10

for v in {1..4}
do
    rm -rf ~/.yoda
    yoda config chain-id bandchain
    yoda config node tcp://172.18.0.1$v:26657
    yoda config chain-rest-server http://172.18.0.20:1317
    yoda config validator $(bandcli keys show validator$v -a --bech val --keyring-backend test)

    # activate validator
    echo "y" | bandcli tx oracle activate --from validator$v --keyring-backend test

    # wait for activation transaction success
    sleep 2

    for i in $(eval echo {1..5})
    do
    # add reporter key
    yoda keys add reporter$i

    # send band tokens to reporter
    echo "y" | bandcli tx send validator$v $(yoda keys show reporter$i) 1000000uband --keyring-backend test

    # wait for sending band tokens transaction success
    sleep 2

    # add reporter to bandchain
    echo "y" | bandcli tx oracle add-reporter $(yoda keys show reporter$i) --from validator$v --keyring-backend test

    # wait for addding reporter transaction success
    sleep 2
    done

    docker create --network bandchain_bandchain --name bandchain_oracle${v} band-validator:latest yoda r
    docker cp ~/.yoda bandchain_oracle${v}:/root/.yoda
    docker start bandchain_oracle${v}
done

# Create faucet container
rm -rf ~/.faucet
faucet config chain-id bandchain
faucet config node tcp://172.18.0.15:26657
faucet config port 5005
for i in $(eval echo {1..5})
do
    # add worker key
    faucet keys add worker$i

    # send band tokens to worker
    echo "y" | bandcli tx send requester $(faucet keys show worker$i) 1000000000000uband --keyring-backend test

    # wait for addding reporter transaction success
    sleep 2
done

docker create --network bandchain_bandchain --name bandchain_faucet --ip 172.18.0.17 band-validator:latest faucet r
docker cp ~/.faucet bandchain_faucet:/root/.faucet
docker start bandchain_faucet
