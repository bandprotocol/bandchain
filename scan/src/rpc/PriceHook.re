type t = {
  usdPrice: float,
  usdMarketCap: float,
  usd24HrChange: float,
  btcPrice: float,
  btcMarketCap: float,
  circulatingSupply: float,
};

let getBandUsd24Change = () => {
  let coingeckoPromise =
    Axios.get(
      "https://api.coingecko.com/api/v3/simple/price?ids=band-protocol&vs_currencies=usd&include_market_cap=true&include_24hr_change=true",
    )
    |> Js.Promise.then_(result =>
         Promise.ret(
           result##data
           |> JsonUtils.Decode.at(["band-protocol", "usd_24h_change"], JsonUtils.Decode.float),
         )
       );

  let cryptocomparePromise =
    Axios.get("https://min-api.cryptocompare.com/data/pricemultifull?fsyms=BAND&tsyms=USD")
    |> Js.Promise.then_(result =>
         Promise.ret(
           result##data
           |> JsonUtils.Decode.at(
                ["RAW", "BAND", "USD", "CHANGEPCT24HOUR"],
                JsonUtils.Decode.float,
              ),
         )
       );

  Js.Promise.race([|coingeckoPromise, cryptocomparePromise|]);
};

let getBandInfo = client => {
  let ratesPromise = client->BandChainJS.getReferenceData([|"BAND/USD", "BAND/BTC"|]);
  // let supplyPromise = Axios.get("https://supply.bandchain.org/circulating");
  let usd24HrChangePromise = getBandUsd24Change();

  let%Promise (rates, usd24HrChange) = Js.Promise.all2((ratesPromise, usd24HrChangePromise));
  
  let bandInfoOpt = {
    let%Opt {rate: bandUsd} = rates->Belt.Array.get(0);
    let%Opt {rate: bandBtc} = rates->Belt.Array.get(1);
    // let supply = supplyData##data;
    let supply = 35191821.;

    Some({
      usdPrice: bandUsd,
      usdMarketCap: bandUsd *. supply,
      usd24HrChange,
      btcPrice: bandBtc,
      btcMarketCap: bandBtc *. supply,
      circulatingSupply: supply,
    });
  };
  bandInfoOpt->Promise.ret;
};
