#!/bin/bash

# Set `DIR` to your path of genesis directory.
DIR=~/genesis_ds_os/genesis

bandd add-data-source \
	"CoinGecko Cryptocurrency Price" \
	"Retrieves current price of a cryptocurrency from https://www.coingecko.com" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/datasources/coingecko_price.py

bandd add-data-source \
	"CryptoCompare Cryptocurrency Price" \
	"Retrieves current price of a cryptocurrency from https://www.cryptocompare.com" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/datasources/cryptocompare_price.py

bandd add-data-source \
	"Binance Cryptocurrency Price" \
	"Retrieves current price of a cryptocurrency from https://www.binance.com/en" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/datasources/binance_price.py

bandd add-data-source \
	"Open Weather Map Weather Data" \
	"Retrieves current weather information from https://www.openweathermap.org" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/datasources/open_weather_map.py

bandd add-data-source \
	"FreeForexAPI Gold Price" \
	"Retrives current gold price from https://www.freeforexapi.com" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/datasources/gold_price.py

bandd add-data-source \
	"Alpha Vantage Stock Price" \
	"Retrives current price of a stock from https://www.alphavaage.co" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/datasources/alphavantage.py

bandd add-data-source \
	"Blockchain.info Bitcoin Block Count" \
	"Retrives latest Bitcoin block height from https://www.blockchain.info" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/datasources/bitcoin_count.py

bandd add-data-source \
	"BlockCypher Bitcoin Block Hash" \
	"Retrives Bitcoin block hash at a given block height from https://blockcypher.com" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/datasources/bitcoin_hash.py

bandd add-data-source \
	"CoinGecko Cryptocurrency Trading Volume" \
	"Retrieves current trading volume of a cryptocurrency from https://www.coingecko.com" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/datasources/coingecko_volume.py

bandd add-data-source \
	"CryptoCompare Cryptocurrency Trading Volume" \
	"Retrieves current trading volume of a cryptocurrency from https://www.cryptocompare.com" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/datasources/cryptocompare_volume.py

bandd add-data-source \
	"ETH Gas Station Current Ethereum Gas Price" \
	"Retrieves current Ethereum gas price from https://ethgasstation.info" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/datasources/ethgasstation.py

bandd add-data-source \
	"Open Sky Network Flight Data" \
	"Retrieves flight information from https://opensky-network.org" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/datasources/open_sky_network.py

bandd add-data-source \
	"Quantum Random Number Generator" \
	"Retrieves array of random number from https://qrng.anu.edu.au" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/datasources/qrng_anu.py

bandd add-data-source \
	"Yahoo Finance Stock Price" \
	"Retrieves current price of a stock from https://finance.yahoo.com" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/datasources/yahoo_finance.py

bandd add-oracle-script \
	"Cryptocurrency Price in USD" \
	"Oracle script that queries the average cryptocurrency price using current price data from CoinGecko, CryptoCompare, and Binance" \
	"{symbol:string,multiplier:u64}/{px:u64}" \
	"https://ipfs.io/ipfs/QmQqxHLszpbCy8Hk2ame3pPAxUUAyStBrVdGdDgrfAngAv" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/res/crypto_price.wasm

bandd add-oracle-script \
	"Gold Price in ATOMs" \
	"Oracle script that queries current average gold price in ATOMs using gold price data from FreeForexAPI and ATOM price from Binance" \
	"{multiplier:u64}/{px:u64}" \
	"https://ipfs.io/ipfs/QmWRhyKRbLB3t19EgC8PjhQr3A7UYsENetvyoAyXLtLRAk" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/res/gold_price.wasm

bandd add-oracle-script \
	"Stock Price (Alpha Vantage)" \
	"Oracle script that queries the current price of a stock from Alpha Vantage" \
	"{symbol:string,api_key:string,multiplier:u64}/{px:u64}" \
	"https://ipfs.io/ipfs/QmS4szs6irBJwyZXGnLnMNedBKhYitv4Q3AAiPyqFGnKGP" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/res/alphavantage.wasm

bandd add-oracle-script \
	"Latest Bitcoin Block Count" \
	"Oracle script that queries the latest Bitcoin block height from Blockchain.info" \
	"{_unused:u8}/{block_count:u64}" \
	"https://ipfs.io/ipfs/QmQFSuZasGZJrnPmo7q3mauLGeKT7fctXprGktRQ8kMGF6" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/res/bitcoin_block_count.wasm

bandd add-oracle-script \
	"Bitcoin Block Hash" \
	"Oracle script that queries the Bitcoin block hash at the given block height from BlockCypher" \
	"{block_height:u64}/{block_hash:string}" \
	"https://ipfs.io/ipfs/QmV8MMq75gZDCergfRwhthPNG9su2GjFYsjyfaw4w9jBrQ" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/res/bitcoin_block_hash.wasm

bandd add-oracle-script \
	"CoinGecko Cryptocurrency Volume" \
	"Oracle script that queries a cryptocurrency's average trading volume for the past day from Coingecko" \
	"{symbol:string,multiplier:u64}/{volume:u64}" \
	"https://ipfs.io/ipfs/QmQK2EXY345gCdAjb7soXp8qj88a2uNe3mcaq4y4QxjJAz" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/res/coingecko_volume.wasm

bandd add-oracle-script \
	"CryptoCompare Cryptocurrency Volume" \
	"Oracle script that queries a cryptocurrency's average trading volume for the past day from CryptoCompare" \
	"{symbol:string,multiplier:u64}/{volume:u64}" \
	"https://ipfs.io/ipfs/QmQECnCsuTbMuYyeWEFgse2nTHHHZizxegxoyTQDZ2vHgX" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/res/crypto_compare_volume.wasm

bandd add-oracle-script \
	"Ethereum Gas Price" \
	"Oracle script that queries the current Ethereum gas price from ETH gas station" \
	"{gas_option:string}/{gweix10:u64}" \
	"https://ipfs.io/ipfs/QmSmynKucQFiBnU5Fbs8N61VXE7mf6b27HCq6rpEAkqx3m" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/res/eth_gas_station.wasm

bandd add-oracle-script \
	"Open Sky Network" \
	"Oracle script for checking whether a given flight number exists during the given time period from OpenSky Network" \
	"{flight_op:string,airport:string,icao24:string,begin:string,end:string}/{flight_existence:bool}" \
	"https://ipfs.io/ipfs/QmRgqLyJgQaUvA1u4Y37u5GxnFz1uXT26kh92bWkHtYWwe" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/res/open_sky_network.wasm

bandd add-oracle-script \
	"Open Weather Map" \
	"Oracle script that queries the current weather data of a location from OpenWeatherMap" \
	"{country:string,main_field:string,sub_field:string,multiplier:u64}/{value:u64}" \
	"https://ipfs.io/ipfs/QmSFdZfnBFEFVyac2WrpugTMpZueA119crrwFCtSRr7sj2" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/res/open_weather_map.wasm

bandd add-oracle-script \
	"Quantum Random Number Generator" \
	"Oracle script that queries a large random number from  Australia's National University Quantum Random API." \
	"{size:u64}/{random_bytes:string}" \
	"https://ipfs.io/ipfs/QmSUirhmjeFMsFao7WttLFGFymse8NHRe7XuWFKT78Ekem" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/res/qrng.wasm

bandd add-oracle-script \
	"Yahoo Stock Price" \
	"Oracle script that queries the current price of a stock from Yahoo Finance" \
	"{symbol:string,multiplier:u64}/{px:u64}" \
	"https://ipfs.io/ipfs/QmNmeuopTDgMUq2D1TEe3V1Hi6WJeLe3bVe9GpCRsDgmwm" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/res/yahoo_price.wasm

bandd add-oracle-script \
	"Fair Cryptocurrency Market Price" \
	"Oracle script that queries the current price of a cryptocurrency from CoinGecko, CryptoCompare, and Binance and aggregates them using the user-selected method in the selected  base currency" \
	"{base_symbol:string,quote_symbol:string,aggregation_method:string,multiplier:u64}/{px:u64}" \
	"https://ipfs.io/ipfs/QmduiYT3GZkkYnhvbr6uK3NDhG9QKDZuacLkRGsYPs8f2Q" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/res/fair_crypto_market_price.wasm

