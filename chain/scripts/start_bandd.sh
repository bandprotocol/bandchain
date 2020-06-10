DIR=`dirname "$0"`

mkdir -p pkg/owasm/res

# release mode
cd ../owasm/chaintests

for f in *; do
    if [ -d "$f" ]; then
        RUSTFLAGS='-C link-arg=-s' cargo build --target wasm32-unknown-unknown --release --package $f
        cp ../target/wasm32-unknown-unknown/release/$f.wasm ../../chain/pkg/owasm/res
    fi
done
cd ../../chain

rm -rf ~/.band*
dropdb my_db
createdb my_db

# initial new node
bandd init validator --chain-id bandchain
echo "lock nasty suffer dirt dream fine fall deal curtain plate husband sound tower mom crew crawl guard rack snake before fragile course bacon range" \
    | bandcli keys add validator --recover --keyring-backend test
echo "smile stem oven genius cave resource better lunar nasty moon company ridge brass rather supply used horn three panic put venue analyst leader comic" \
    | bandcli keys add requester --recover --keyring-backend test

# add accounts to genesis
bandd add-genesis-account validator 10000000000000uband --keyring-backend test
bandd add-genesis-account requester 10000000000000uband --keyring-backend test

# add data sources to genesis
bandd add-data-source \
    "CoinGecko Cryptocurrency Price" \
    "Retrieves current price of a cryptocurrency from https://www.coingecko.com" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./datasources/coingecko_price.py
bandd add-data-source \
	"CryptoCompare Cryptocurrency Price" \
	"Retrieves current price of a cryptocurrency from https://www.cryptocompare.com" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/crypto_compare_price.sh
bandd add-data-source \
	"Binance Cryptocurrency Price" \
	"Retrieves current price of a cryptocurrency from https://www.binance.com/en" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/binance_price.sh
bandd add-data-source \
	"Open Weather Map Weather Data" \
	"Retrieves current weather information from https://www.openweathermap.org" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/open_weather_map.sh
bandd add-data-source \
	"FreeForexAPI Gold Price" \
	"Retrives current gold price from https://www.freeforexapi.com" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/gold_price.sh
bandd add-data-source \
	"Alpha Vantage Stock Price" \
	"Retrives current price of a stock from https://www.alphavaage.co" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/alphavantage.sh
bandd add-data-source \
	"Blockchain.info Bitcoin Block Count" \
	"Retrives latest Bitcoin block height from https://www.blockchain.info" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/bitcoin_count.sh
bandd add-data-source \
	"BlockCypher Bitcoin Block Hash" \
	"Retrives Bitcoin block hash at a given block height from https://blockcypher.com" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/bitcoin_hash.sh
bandd add-data-source \
	"CoinGecko Cryptocurrency Trading Volume" \
	"Retrieves current trading volume of a cryptocurrency from https://www.coingecko.com" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/coingecko_volume.sh
bandd add-data-source \
	"CryptoCompare Cryptocurrency Trading Volume" \
	"Retrieves current trading volume of a cryptocurrency from https://www.cryptocompare.com" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/crypto_compare_volume.sh
bandd add-data-source \
	"ETH Gas Station Current Ethereum Gas Price" \
	"Retrieves current Ethereum gas price from https://ethgasstation.info" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/ethgasstation.sh
bandd add-data-source \
	"Open Sky Network Flight Data" \
	"Retrieves flight information from https://opensky-network.org" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/open_sky_network.sh
bandd add-data-source \
	"Quantum Random Number Generator" \
	"Retrieves array of random number from https://qrng.anu.edu.au" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/qrng_anu.sh
bandd add-data-source \
	"Yahoo Finance Stock Price" \
	"Retrieves current price of a stock from https://finance.yahoo.com" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	./datasources/yahoo_finance.py

bandd add-oracle-script \
	"Cryptocurrency Price in USD" \
	"Oracle script that queries the average cryptocurrency price using current price data from CoinGecko, CryptoCompare, and Binance" \
	"{symbol:string,multiplier:u64}/{px:u64}" \
    "https://ipfs.io/ipfs/QmUbdfoRR9ge6P39EoqDjBhQoDeaT6gJu76Ce9kKsz938N" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/crypto_price.wasm
bandd add-oracle-script \
	"Gold Price in ATOMs" \
	"Oracle script that queries current average gold price in ATOMs using gold price data from FreeForexAPI and ATOM price from Binance" \
	"{multiplier:u64}/{px:u64}" \
    "https://ipfs.io/ipfs/Qmbcdr3UZXMrJeoRtHzTtHHepnzjyX1gWNhewWe6BXgmPm" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/gold_price.wasm
bandd add-oracle-script \
	"Stock Price (Alpha Vantage)" \
	"Oracle script that queries the current price of a stock from Alpha Vantage" \
	"{symbol:string,api_key:string,multiplier:u64}/{px:u64}" \
    "https://ipfs.io/ipfs/QmPsSmJ9gEdBoeQqwtk6bJykyFtqSpztCeGb9J1VFW65av" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/alphavantage.wasm
bandd add-oracle-script \
	"Latest Bitcoin Block Count" \
	"Oracle script that queries the latest Bitcoin block height from Blockchain.info" \
	"{_unused:u8}/{block_count:u64}" \
    "https://ipfs.io/ipfs/QmUkpTCvdKMEFxwgeTpjP9hszZ11e5ioXZAS7XLpQLbV2k" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/bitcoin_block_count.wasm
bandd add-oracle-script \
	"Bitcoin Block Hash" \
	"Oracle script for getting the Bitcoin block hash at the given block height""Oracle script that queries the Bitcoin block hash at the given block height from BlockCypher"\
	"{block_height:u64}/{block_hash:string}" \
    "https://ipfs.io/ipfs/QmXu5NyUrtbcdPxut4WhVsRT4KjsPjy2NwJzEts7rjuEDf" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/bitcoin_block_hash.wasm
bandd add-oracle-script \
	"CoinGecko Cryptocurrency Volume" \
	"Oracle script that queries a cryptocurrency's average trading volume for the past day from Coingecko" \
	"{symbol:string,multiplier:u64}/{volume:u64}" \
    "https://ipfs.io/ipfs/QmVuYP5cSujNSv33ZNMiFbcSMoRtBa7WMzG2q55j21Vhxj" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/coingecko_volume.wasm
bandd add-oracle-script \
	"CryptoCompare Cryptocurrency Volume" \
	"Oracle script that queries a cryptocurrency's average trading volume for the past day from CryptoCompare" \
	"{symbol:string,multiplier:u64}/{volume:u64}" \
    "https://ipfs.io/ipfs/Qmf2e5VF3uscGzBMwfQRZbaYLWxAZPWAgwv3m3iSv6BoGE" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/crypto_compare_volume.wasm
bandd add-oracle-script \
	"Ethereum Gas Price" \
	"Oracle script that queries the current Ethereum gas price from ETH gas station" \
	"{gas_option:string}/{gweix10:u64}" \
    "https://ipfs.io/ipfs/QmP1i61XdPnfKSewh7vyh3xLgjxT42Gqpiv7CYLFK6V3Mg" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/eth_gas_station.wasm
bandd add-oracle-script \
	"Open Sky Network" \
	"Oracle script for checking whether a given flight number exists during the given time period from OpenSky Network" \
	"{flight_op:string,icao24:string,begin:string,end:string}/{flight_existence:u8}" \
    "https://ipfs.io/ipfs/QmST4us1xAXmfXZFqBRZqjsDpNhTMY1CRsgjNccwmB4FTX" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/open_sky_network.wasm
bandd add-oracle-script \
	"Open Weather Map" \
	"Oracle script that queries the current weather data of a location from OpenWeatherMap" \
	"{country:string,main_field:string,sub_field:string,multiplier:u64}/{value:u64}" \
    "https://ipfs.io/ipfs/QmNWvYfqZztrMNKjKyKLTATvnVPVUMPCdLJFkV5HfHBQoo" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/open_weather_map.wasm
bandd add-oracle-script \
	"Quantum Random Number Generator" \
	"Oracle script that queries a large random number from  Australia's National University Quantum Random API." \
	"{size:u64}/{random_bytes:string}" \
    "https://ipfs.io/ipfs/QmZ62dxgAmCtDnt5XcAs2zP4UjGzozDzdKiHoR1Wo9MVeV" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/qrng.wasm
bandd add-oracle-script \
	"Yahoo Stock Price" \
	"Oracle script that queries the current price of a stock from Yahoo Finance" \
	"{symbol:string,multiplier:u64}/{px:u64}" \
    "https://ipfs.io/ipfs/QmfEUKFoX9PY3LHnT7Deixwb8qRrgWvdf5v8MzTinTYXLu" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/yahoo_price.wasm
bandd add-oracle-script \
	"Fair Cryptocurrency Market Price" \
	"Oracle script that queries the current price of a cryptocurrency from CoinGecko, CryptoCompare, and Binance and aggregates them using the user-selected method in the selected  base currency" \
	"{base_symbol:string,quote_symbol:string,aggregation_method:string,multiplier:u64}/{px:u64}" \
    "https://ipfs.io/ipfs/QmbnRei1WG8gdstsVuU7Qqq4PwqED9LuHFvjDgS5asShoM" \
    band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
    ./pkg/owasm/res/fair_crypto_market_price.wasm

# genesis configurations
bandcli config chain-id bandchain
bandcli config output json
bandcli config indent true
bandcli config trust-node true

# register initial validators
bandd gentx \
    --amount 100000000uband \
    --node-id 11392b605378063b1c505c0ab123f04bd710d7d7 \
    --pubkey bandvalconspub1addwnpepq06h7wvh5n5pmrejr6t3pyn7ytpwd5c0kmv0wjdfujs847em8dusjl96sxg \
    --name validator \
    --keyring-backend test

# collect genesis transactions
bandd collect-gentxs
cp ./docker-config/single-validator/priv_validator_key.json ~/.bandd/config/priv_validator_key.json
cp ./docker-config/single-validator/node_key.json ~/.bandd/config/node_key.json

# start bandchain
bandd start --with-db "postgres: port=5432 user=$USER dbname=my_db sslmode=disable" \
  --rpc.laddr tcp://0.0.0.0:26657 --pruning=nothing
