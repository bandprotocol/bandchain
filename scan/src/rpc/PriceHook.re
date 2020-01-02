module Price = {
  type t = {
    usdPrice: float,
    usdMarketCap: float,
    usd24HrChange: float,
    btcPrice: float,
    btcMarketCap: float,
    btc24HrChange: float,
  };

  let decode = (usdJson, btcJson) =>
    JsonUtils.Decode.{
      usdPrice: usdJson |> at(["band-protocol", "usd"], JsonUtils.Decode.float),
      usdMarketCap: usdJson |> at(["band-protocol", "usd_market_cap"], JsonUtils.Decode.float),
      usd24HrChange: usdJson |> at(["band-protocol", "usd_24h_change"], JsonUtils.Decode.float),
      btcPrice: btcJson |> at(["band-protocol", "btc"], JsonUtils.Decode.float),
      btcMarketCap: btcJson |> at(["band-protocol", "btc_market_cap"], JsonUtils.Decode.float),
      btc24HrChange: btcJson |> at(["band-protocol", "btc_24h_change"], JsonUtils.Decode.float),
    };
};

let get = (~pollInterval=?, ()) => {
  let usdJson =
    Axios.use(
      "https://api.coingecko.com/api/v3/simple/price?ids=band-protocol&vs_currencies=usd&include_market_cap=true&include_24hr_change=true",
      ~pollInterval?,
      (),
    );
  let btcJson =
    Axios.use(
      "https://api.coingecko.com/api/v3/simple/price?ids=band-protocol&vs_currencies=btc&include_market_cap=true&include_24hr_change=true",
      ~pollInterval?,
      (),
    );
  usdJson |> Belt.Option.flatMap(_, u => btcJson |> Belt.Option.map(_, Price.decode(u)));
};
