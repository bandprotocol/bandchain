DIR=$(dirname "$0")
bandd add-data-source \
	"CoinGecko Cryptocurrency Price" \
	"Retrieves current price of a cryptocurrency from https://www.coingecko.com" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../datasources/001_coingecko_price.py

bandd add-data-source \
	"CryptoCompare Cryptocurrency Price" \
	"Retrieves current price of a cryptocurrency from https://www.cryptocompare.com" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../datasources/002_crypto_compare_price.sh

bandd add-data-source \
	"Binance Cryptocurrency Price" \
	"Retrieves current price of a cryptocurrency from https://www.binance.com/en" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../datasources/003_binance_price.sh

bandd add-data-source \
	"Open Weather Map Weather Data" \
	"Retrieves current weather information from https://www.openweathermap.org" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../datasources/004_open_weather_map.sh

bandd add-data-source \
	"FreeForexAPI Gold Price" \
	"Retrives current gold price from https://www.freeforexapi.com" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../datasources/005_gold_price.sh

bandd add-data-source \
	"Alpha Vantage Stock Price" \
	"Retrives current price of a stock from https://www.alphavaage.co" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../datasources/006_alphavantage.sh

bandd add-data-source \
	"Blockchain.info Bitcoin Block Count" \
	"Retrives latest Bitcoin block height from https://www.blockchain.info" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../datasources/007_bitcoin_count.sh

bandd add-data-source \
	"BlockCypher Bitcoin Block Hash" \
	"Retrives Bitcoin block hash at a given block height from https://blockcypher.com" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../datasources/008_bitcoin_hash.sh

bandd add-data-source \
	"CoinGecko Cryptocurrency Trading Volume" \
	"Retrieves current trading volume of a cryptocurrency from https://www.coingecko.com" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../datasources/009_coingecko_volume.sh

bandd add-data-source \
	"CryptoCompare Cryptocurrency Trading Volume" \
	"Retrieves current trading volume of a cryptocurrency from https://www.cryptocompare.com" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../datasources/010_crypto_compare_volume.sh

bandd add-data-source \
	"ETH Gas Station Current Ethereum Gas Price" \
	"Retrieves current Ethereum gas price from https://ethgasstation.info" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../datasources/011_ethgasstation.sh

bandd add-data-source \
	"Open Sky Network Flight Data" \
	"Retrieves flight information from https://opensky-network.org" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../datasources/012_open_sky_network.sh

bandd add-data-source \
	"Quantum Random Number Generator" \
	"Retrieves array of random number from https://qrng.anu.edu.au" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../datasources/013_qrng_anu.sh

bandd add-data-source \
	"Yahoo Finance Stock Price" \
	"Retrieves current price of a stock from https://finance.yahoo.com" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../datasources/014_yahoo_finance.py

bandd add-oracle-script \
	"Cryptocurrency Price in USD" \
	"Oracle script that queries the average cryptocurrency price using current price data from CoinGecko, CryptoCompare, and Binance" \
	"{symbol:string,multiplier:u64}/{px:u64}" \
	"https://ipfs.io/ipfs/QmdMKT62HYaaYH44DrW1UkQNhsd76nZXej6KXWjYtR9c5m" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../pkg/owasm/res/os001_crypto_price.wasm

bandd add-oracle-script \
	"Gold Price in ATOMs" \
	"Oracle script that queries current average gold price in ATOMs using gold price data from FreeForexAPI and ATOM price from Binance" \
	"{multiplier:u64}/{px:u64}" \
	"https://ipfs.io/ipfs/QmZ4oq69uvxBgcm5Lby415Fycc6upRpEe9cjnZzadFAYEq" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../pkg/owasm/res/os002_gold_price.wasm

bandd add-oracle-script \
	"Stock Price (Alpha Vantage)" \
	"Oracle script that queries the current price of a stock from Alpha Vantage" \
	"{symbol:string,api_key:string,multiplier:u64}/{px:u64}" \
	"https://ipfs.io/ipfs/QmWohk9LcGSLTGdmUddCXMCa4Y96KAYkXQwba2Gfr2MTGw" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../pkg/owasm/res/os003_alphavantage.wasm

bandd add-oracle-script \
	"Latest Bitcoin Block Count" \
	"Oracle script that queries the latest Bitcoin block height from Blockchain.info" \
	"{_unused:u8}/{block_count:u64}" \
	"https://ipfs.io/ipfs/QmdjEMcoUJfVSLp5nYNjuozGqr3p8X2Td2t1XgyBEspcFm" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../pkg/owasm/res/os004_bitcoin_block_count.wasm

bandd add-oracle-script \
	"Bitcoin Block Hash" \
	"Oracle script for getting the Bitcoin block hash at the given block heightOracle script that queries the Bitcoin block hash at the given block height from BlockCypher" \
	"{block_height:u64}/{block_hash:string}" \
	"https://ipfs.io/ipfs/QmVYbGi5r7o5Aup97jZ6XMAp45tmiTjfwuSppPTZxcNzWQ" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../pkg/owasm/res/os005_bitcoin_block_hash.wasm

bandd add-oracle-script \
	"CoinGecko Cryptocurrency Volume" \
	"Oracle script that queries a cryptocurrency's average trading volume for the past day from Coingecko" \
	"{symbol:string,multiplier:u64}/{volume:u64}" \
	"https://ipfs.io/ipfs/QmWWV2uS2nj9wzz29YhAezXqpy5sm48kwQnifFcNHNCM4V" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../pkg/owasm/res/os006_coingecko_volume.wasm

bandd add-oracle-script \
	"CryptoCompare Cryptocurrency Volume" \
	"Oracle script that queries a cryptocurrency's average trading volume for the past day from CryptoCompare" \
	"{symbol:string,multiplier:u64}/{volume:u64}" \
	"https://ipfs.io/ipfs/Qmc4NxRirdYZyyou8JuVWLMWgor26WaRLXhoVFDEvv5MN5" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../pkg/owasm/res/os007_crypto_compare_volume.wasm

bandd add-oracle-script \
	"Ethereum Gas Price" \
	"Oracle script that queries the current Ethereum gas price from ETH gas station" \
	"{gas_option:string}/{gweix10:u64}" \
	"https://ipfs.io/ipfs/QmX3FE4QnGNJw9a5M4cgfXb9bBiTHdKzy2bN3Q3QUPz5KE" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../pkg/owasm/res/os008_eth_gas_station.wasm

bandd add-oracle-script \
	"Open Sky Network" \
	"Oracle script for checking whether a given flight number exists during the given time period from OpenSky Network" \
	"{flight_op:string,airport:string,icao24:string,begin:string,end:string}/{flight_existence:bool}" \
	"https://ipfs.io/ipfs/QmcHVJbQUPzKLYpBgMU8Gm99QvGo3fJabPMeCd4KFbEuQj" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../pkg/owasm/res/os009_open_sky_network.wasm

bandd add-oracle-script \
	"Open Weather Map" \
	"Oracle script that queries the current weather data of a location from OpenWeatherMap" \
	"{country:string,main_field:string,sub_field:string,multiplier:u64}/{value:u64}" \
	"https://ipfs.io/ipfs/QmdPcvkXZk47AzkV5NxxYm3fnToYLewfsvP8BRSMpySQGK" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../pkg/owasm/res/os010_open_weather_map.wasm

bandd add-oracle-script \
	"Quantum Random Number Generator" \
	"Oracle script that queries a large random number from  Australia's National University Quantum Random API." \
	"{size:u64}/{random_bytes:string}" \
	"https://ipfs.io/ipfs/QmSoZQXsW5GqJivazw4LYB93pj5zW8qyNaWeBM2NfXkeHc" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../pkg/owasm/res/os011_qrng.wasm

bandd add-oracle-script \
	"Yahoo Stock Price" \
	"Oracle script that queries the current price of a stock from Yahoo Finance" \
	"{symbol:string,multiplier:u64}/{px:u64}" \
	"https://ipfs.io/ipfs/QmQZ6UZwFYQLxSBvBS2sNm6p4D1hUkXxxEp5kiDCjwMFzH" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../pkg/owasm/res/os012_yahoo_price.wasm

bandd add-oracle-script \
	"Fair Cryptocurrency Market Price" \
	"Oracle script that queries the current price of a cryptocurrency from CoinGecko, CryptoCompare, and Binance and aggregates them using the user-selected method in the selected  base currency" \
	"{base_symbol:string,quote_symbol:string,aggregation_method:string,multiplier:u64}/{px:u64}" \
	"https://ipfs.io/ipfs/QmSNvRVjTsSaGSyFgXnbbDyoS7q4YxJRdysrnC7MhUHTNi" \
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \
	$DIR/../pkg/owasm/res/os013_fair_crypto_market_price.wasm

