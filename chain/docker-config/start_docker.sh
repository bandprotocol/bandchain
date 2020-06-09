#!/bin/bash

DIR=`dirname "$0"`

# remove old genesis
rm -rf ~/.band*

mkdir -p pkg/owasm/res

# Build genesis oracle scripts
cd ../owasm/chaintests

for f in *; do
    if [ -d "$f" ]; then
        RUSTFLAGS='-C link-arg=-s' cargo build --target wasm32-unknown-unknown --release --package $f
        cp ../target/wasm32-unknown-unknown/release/$f.wasm ../../chain/pkg/owasm/res
    fi
done

cd ../../chain

make install

# initial new node
bandd init node-validator --chain-id bandchain

# add data sources to genesis
bandd add-data-source \
    "Coingecko script" \
    "The script that queries crypto price from https://coingecko.com" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./datasources/coingecko_price.py
bandd add-data-source \
	"Crypto compare script" \
	"The script that queries crypto price from https://cryptocompare.com" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/crypto_compare_price.sh
bandd add-data-source \
	"Binance price" \
	"The script that queries crypto price from https://www.binance.com/en" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/binance_price.sh
bandd add-data-source \
	"Open weather map" \
	"The script that queries weather information from https://api.openweathermap.org" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/open_weather_map.sh
bandd add-data-source \
	"Gold price" \
	"The script that queries current gold price" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/gold_price.sh
bandd add-data-source \
	"Alphavantage" \
	"The script that queries stock price from Alphavantage" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/alphavantage.sh
bandd add-data-source \
	"Bitcoin block count" \
	"The script that queries latest block height of Bitcoin" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/bitcoin_count.sh
bandd add-data-source \
	"Bitcoin block hash" \
	"The script that queries block hash of Bitcoin" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/bitcoin_hash.sh
bandd add-data-source \
	"Coingecko volume script" \
	"The script that queries crypto volume from Coingecko" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/coingecko_volume.sh
bandd add-data-source \
	"Crypto compare volume script" \
	"The script that queries crypto volume from Crypto compare" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/crypto_compare_volume.sh
bandd add-data-source \
	"ETH gas station" \
	"The script that queries current Ethereum gas price https://ethgasstation.info" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/ethgasstation.sh
bandd add-data-source \
	"Open sky network" \
	"The script that queries flight information from https://opensky-network.org" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/open_sky_network.sh
bandd add-data-source \
	"Quantum random numbers" \
	"The script that queries array of random number from https://qrng.anu.edu.au" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/qrng_anu.sh
bandd add-data-source \
	"Yahoo finance" \
	"The script that queries stock price from https://finance.yahoo.com" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/yahoo_finance.py

bandd add-oracle-script \
	"Crypto price script" \
	"Oracle script for getting the current an average cryptocurrency price from various sources." \
	"{symbol:string,multiplier:u64}/{px:u64}" \
    "https://ipfs.io/ipfs/QmUbdfoRR9ge6P39EoqDjBhQoDeaT6gJu76Ce9kKsz938N" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/crypto_price.wasm
bandd add-oracle-script \
	"Gold price script" \
	"Oracle script for getting the current average gold price in ATOMs" \
	"{multiplier:u64}/{px:u64}" \
    "https://ipfs.io/ipfs/Qmbcdr3UZXMrJeoRtHzTtHHepnzjyX1gWNhewWe6BXgmPm" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/gold_price.wasm
bandd add-oracle-script \
	"Alpha Vantage stock price script" \
	"Oracle script for getting the current price of a stock from Alpha Vantage" \
	"{symbol:string,api_key:string,multiplier:u64}/{px:u64}" \
    "https://ipfs.io/ipfs/QmPsSmJ9gEdBoeQqwtk6bJykyFtqSpztCeGb9J1VFW65av" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/alphavantage.wasm
bandd add-oracle-script \
	"Bitcoin block count" \
	"Oracle script for getting the latest Bitcoin block height" \
	"{_unused:u8}/{block_count:u64}" \
    "https://ipfs.io/ipfs/QmUkpTCvdKMEFxwgeTpjP9hszZ11e5ioXZAS7XLpQLbV2k" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/bitcoin_block_count.wasm
bandd add-oracle-script \
	"Bitcoin block hash" \
	"Oracle script for getting the Bitcoin block hash at the given block height" \
	"{block_height:u64}/{block_hash:string}" \
    "https://ipfs.io/ipfs/QmXu5NyUrtbcdPxut4WhVsRT4KjsPjy2NwJzEts7rjuEDf" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/bitcoin_block_hash.wasm
bandd add-oracle-script \
	"CoinGecko crypto volume" \
	"Oracle script for getting a cryptocurrency's average trading volume for the past day from Coingecko" \
	"{symbol:string,multiplier:u64}/{volume:u64}" \
    "https://ipfs.io/ipfs/QmVuYP5cSujNSv33ZNMiFbcSMoRtBa7WMzG2q55j21Vhxj" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/coingecko_volume.wasm
bandd add-oracle-script \
	"CryptoCompare crypto volume" \
	"Oracle script for getting a cryptocurrency's average trading volume for the past day from CryptoCompare" \
	"{symbol:string,multiplier:u64}/{volume:u64}" \
    "https://ipfs.io/ipfs/Qmf2e5VF3uscGzBMwfQRZbaYLWxAZPWAgwv3m3iSv6BoGE" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/crypto_compare_volume.wasm
bandd add-oracle-script \
	"Ethereum gas price" \
	"Oracle script for getting the current Ethereum gas price from ETH gas station" \
	"{gas_option:string}/{gweix10:u64}" \
    "https://ipfs.io/ipfs/QmP1i61XdPnfKSewh7vyh3xLgjxT42Gqpiv7CYLFK6V3Mg" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/eth_gas_station.wasm
bandd add-oracle-script \
	"Open sky network" \
	"Oracle script for checking whether a given flight number exists during the given time period" \
	"{flight_op:string,icao24:string,begin:string,end:string}/{flight_existence:u8}" \
    "https://ipfs.io/ipfs/QmST4us1xAXmfXZFqBRZqjsDpNhTMY1CRsgjNccwmB4FTX" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/open_sky_network.wasm
bandd add-oracle-script \
	"Open weather map" \
	"Oracle script for getting the current weather data of a location" \
	"{country:string,main_field:string,sub_field:string,multiplier:u64}/{value:u64}" \
    "https://ipfs.io/ipfs/QmNWvYfqZztrMNKjKyKLTATvnVPVUMPCdLJFkV5HfHBQoo" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/open_weather_map.wasm
bandd add-oracle-script \
	"Quantum random number generator" \
	"Oracle script for getting a big random number from quantum computer" \
	"{size:u64}/{random_bytes:string}" \
    "https://ipfs.io/ipfs/QmZ62dxgAmCtDnt5XcAs2zP4UjGzozDzdKiHoR1Wo9MVeV" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/qrng.wasm
bandd add-oracle-script \
	"Yahoo stock price" \
	"Oracle script for getting the current price of a stock from Yahoo Finance" \
	"{symbol:string,multiplier:u64}/{px:u64}" \
    "https://ipfs.io/ipfs/QmfEUKFoX9PY3LHnT7Deixwb8qRrgWvdf5v8MzTinTYXLu" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/yahoo_price.wasm
bandd add-oracle-script \
	"Fair price from 3 sources" \
	"Oracle script that query prices from many markets and then aggregate them together" \
	"{base_symbol:string,quote_symbol:string,aggregation_method:string,multiplier:u64}/{px:u64}" \
    "https://ipfs.io/ipfs/QmbnRei1WG8gdstsVuU7Qqq4PwqED9LuHFvjDgS5asShoM" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/fair_crypto_market_price.wasm

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
bandd add-genesis-account requester 10000000000000uband --keyring-backend test

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

# Recreate files volume
docker volume rm query-files
docker volume create --driver local \
    --opt type=none \
    --opt device=$HOME/.bandd/files \
    --opt o=bind query-files

cd ..

docker-compose up -d --build

sleep 30

for v in {1..4}
do
    rm -rf ~/.oracled
    bandoracled2 config chain-id bandchain
    bandoracled2 config node tcp://172.18.0.1$v:26657
    bandoracled2 config chain-rest-server http://172.18.0.20:1317
    bandoracled2 config validator $(bandcli keys show validator$v -a --bech val --keyring-backend test)

    for i in $(eval echo {1..5})
    do
    # add reporter key
    bandoracled2 keys add reporter$i

    # send band tokens to reporter
    echo "y" | bandcli tx send validator$v $(bandoracled2 keys show reporter$i) 1000000uband --keyring-backend test

    # wait for sending band tokens transaction success
    sleep 2

    # add reporter to bandchain
    echo "y" | bandcli tx oracle add-reporter $(bandoracled2 keys show reporter$i) --from validator$v --keyring-backend test

    # wait for addding reporter transaction success
    sleep 2
    done

    docker create --network bandchain_bandchain --name bandchain_oracle${v} band-validator:latest bandoracled2 r
    docker cp ~/.oracled bandchain_oracle${v}:/root/.oracled
    docker start bandchain_oracle${v}
done
