bandd add-data-source "001_coingecko_price.py" \ 
	"001_coingecko_price.py" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/datasources/001_coingecko_price.py

bandd add-data-source "002_crypto_compare_price.sh" \ 
	"002_crypto_compare_price.sh" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/datasources/002_crypto_compare_price.sh

bandd add-data-source "003_binance_price.sh" \ 
	"003_binance_price.sh" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/datasources/003_binance_price.sh

bandd add-data-source "004_open_weather_map.sh" \ 
	"004_open_weather_map.sh" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/datasources/004_open_weather_map.sh

bandd add-data-source "005_gold_price.sh" \ 
	"005_gold_price.sh" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/datasources/005_gold_price.sh

bandd add-data-source "006_alphavantage.sh" \ 
	"006_alphavantage.sh" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/datasources/006_alphavantage.sh

bandd add-data-source "007_bitcoin_count.sh" \ 
	"007_bitcoin_count.sh" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/datasources/007_bitcoin_count.sh

bandd add-data-source "008_bitcoin_hash.sh" \ 
	"008_bitcoin_hash.sh" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/datasources/008_bitcoin_hash.sh

bandd add-data-source "009_coingecko_volume.sh" \ 
	"009_coingecko_volume.sh" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/datasources/009_coingecko_volume.sh

bandd add-data-source "010_crypto_compare_volume.sh" \ 
	"010_crypto_compare_volume.sh" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/datasources/010_crypto_compare_volume.sh

bandd add-data-source "011_ethgasstation.sh" \ 
	"011_ethgasstation.sh" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/datasources/011_ethgasstation.sh

bandd add-data-source "012_open_sky_network.sh" \ 
	"012_open_sky_network.sh" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/datasources/012_open_sky_network.sh

bandd add-data-source "013_qrng_anu.sh" \ 
	"013_qrng_anu.sh" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/datasources/013_qrng_anu.sh

bandd add-data-source "014_yahoo_finance.py" \ 
	"014_yahoo_finance.py" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/datasources/014_yahoo_finance.py

bandd add-oracle-script "001_crypto_price" \ 
	"001_crypto_price" \ 
	"{symbol:string,multiplier:u64}/{px:u64}" \ 
	"https://ipfs.io/ipfs/QmdMKT62HYaaYH44DrW1UkQNhsd76nZXej6KXWjYtR9c5m" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/pkg/owasm/res/001_crypto_price.wasm

bandd add-oracle-script "002_gold_price" \ 
	"002_gold_price" \ 
	"{multiplier:u64}/{px:u64}" \ 
	"https://ipfs.io/ipfs/QmZ4oq69uvxBgcm5Lby415Fycc6upRpEe9cjnZzadFAYEq" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/pkg/owasm/res/002_gold_price.wasm

bandd add-oracle-script "003_alphavantage" \ 
	"003_alphavantage" \ 
	"{symbol:string,api_key:string,multiplier:u64}/{px:u64}" \ 
	"https://ipfs.io/ipfs/QmWohk9LcGSLTGdmUddCXMCa4Y96KAYkXQwba2Gfr2MTGw" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/pkg/owasm/res/003_alphavantage.wasm

bandd add-oracle-script "004_bitcoin_block_count" \ 
	"004_bitcoin_block_count" \ 
	"{_unused:u8}/{block_count:u64}" \ 
	"https://ipfs.io/ipfs/QmdjEMcoUJfVSLp5nYNjuozGqr3p8X2Td2t1XgyBEspcFm" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/pkg/owasm/res/004_bitcoin_block_count.wasm

bandd add-oracle-script "005_bitcoin_block_hash" \ 
	"005_bitcoin_block_hash" \ 
	"{block_height:u64}/{block_hash:string}" \ 
	"https://ipfs.io/ipfs/QmVYbGi5r7o5Aup97jZ6XMAp45tmiTjfwuSppPTZxcNzWQ" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/pkg/owasm/res/005_bitcoin_block_hash.wasm

bandd add-oracle-script "006_coingecko_volume" \ 
	"006_coingecko_volume" \ 
	"{symbol:string,multiplier:u64}/{volume:u64}" \ 
	"https://ipfs.io/ipfs/QmWWV2uS2nj9wzz29YhAezXqpy5sm48kwQnifFcNHNCM4V" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/pkg/owasm/res/006_coingecko_volume.wasm

bandd add-oracle-script "007_crypto_compare_volume" \ 
	"007_crypto_compare_volume" \ 
	"{symbol:string,multiplier:u64}/{volume:u64}" \ 
	"https://ipfs.io/ipfs/Qmc4NxRirdYZyyou8JuVWLMWgor26WaRLXhoVFDEvv5MN5" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/pkg/owasm/res/007_crypto_compare_volume.wasm

bandd add-oracle-script "008_eth_gas_station" \ 
	"008_eth_gas_station" \ 
	"{gas_option:string}/{gweix10:u64}" \ 
	"https://ipfs.io/ipfs/QmX3FE4QnGNJw9a5M4cgfXb9bBiTHdKzy2bN3Q3QUPz5KE" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/pkg/owasm/res/008_eth_gas_station.wasm

bandd add-oracle-script "009_open_sky_network" \ 
	"009_open_sky_network" \ 
	"{flight_op:string,airport:string,icao24:string,begin:string,end:string}/{flight_existence:bool}" \ 
	"https://ipfs.io/ipfs/QmcHVJbQUPzKLYpBgMU8Gm99QvGo3fJabPMeCd4KFbEuQj" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/pkg/owasm/res/009_open_sky_network.wasm

bandd add-oracle-script "010_open_weather_map" \ 
	"010_open_weather_map" \ 
	"{country:string,main_field:string,sub_field:string,multiplier:u64}/{value:u64}" \ 
	"https://ipfs.io/ipfs/QmdPcvkXZk47AzkV5NxxYm3fnToYLewfsvP8BRSMpySQGK" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/pkg/owasm/res/010_open_weather_map.wasm

bandd add-oracle-script "011_qrng" \ 
	"011_qrng" \ 
	"{size:u64}/{random_bytes:string}" \ 
	"https://ipfs.io/ipfs/QmSoZQXsW5GqJivazw4LYB93pj5zW8qyNaWeBM2NfXkeHc" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/pkg/owasm/res/011_qrng.wasm

bandd add-oracle-script "012_yahoo_price" \ 
	"012_yahoo_price" \ 
	"{symbol:string,multiplier:u64}/{px:u64}" \ 
	"https://ipfs.io/ipfs/QmQZ6UZwFYQLxSBvBS2sNm6p4D1hUkXxxEp5kiDCjwMFzH" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/pkg/owasm/res/012_yahoo_price.wasm

bandd add-oracle-script "013_fair_crypto_market_price" \ 
	"013_fair_crypto_market_price" \ 
	"{base_symbol:string,quote_symbol:string,aggregation_method:string,multiplier:u64}/{px:u64}" \ 
	"https://ipfs.io/ipfs/QmSNvRVjTsSaGSyFgXnbbDyoS7q4YxJRdysrnC7MhUHTNi" \ 
	band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs \ 
	/Users/beeb/bandchain/chain/pkg/owasm/res/013_fair_crypto_market_price.wasm

