module Price = {
  type t = {
    usdPrice: float,
    usdMarketCap: float,
    usd24HrChange: float,
    btcPrice: float,
    btcMarketCap: float,
    btc24HrChange: float,
    circulatingSupply: float,
  };
  let decode = (usdJson, btcJson, bandJson) =>
    JsonUtils.Decode.{
      usdPrice: usdJson |> at(["band-protocol", "usd"], JsonUtils.Decode.float),
      usdMarketCap: usdJson |> at(["band-protocol", "usd_market_cap"], JsonUtils.Decode.float),
      usd24HrChange: usdJson |> at(["band-protocol", "usd_24h_change"], JsonUtils.Decode.float),
      btcPrice: btcJson |> at(["band-protocol", "btc"], JsonUtils.Decode.float),
      btcMarketCap: btcJson |> at(["band-protocol", "btc_market_cap"], JsonUtils.Decode.float),
      btc24HrChange: btcJson |> at(["band-protocol", "btc_24h_change"], JsonUtils.Decode.float),
      circulatingSupply:
        bandJson |> at(["market_data", "circulating_supply"], JsonUtils.Decode.float),
    };
};
let get = () => {
  let usdJson =
    AxiosHooks.use(
      "https://api.coingecko.com/api/v3/simple/price?ids=band-protocol&vs_currencies=usd&include_market_cap=true&include_24hr_change=true",
    );
  let btcJson =
    AxiosHooks.use(
      "https://api.coingecko.com/api/v3/simple/price?ids=band-protocol&vs_currencies=btc&include_market_cap=true&include_24hr_change=true",
    );
  let bandJson = AxiosHooks.use("https://api.coingecko.com/api/v3/coins/band-protocol");
  let%Opt usd = usdJson;
  let%Opt btc = btcJson;
  let%Opt band = bandJson;
  Some(Price.decode(usd, btc, band));
};
