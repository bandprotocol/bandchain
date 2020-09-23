type t = {
  usdPrice: float,
  usdMarketCap: float,
  usd24HrChange: float,
  btcPrice: float,
  btcMarketCap: float,
  btc24HrChange: float,
  // circulatingSupply: float,
};

module CoinGekco = {
  let get = () => {
    let decode = (usdJson, btcJson) =>
      JsonUtils.Decode.{
        usdPrice: usdJson |> at(["band-protocol", "usd"], JsonUtils.Decode.float),
        usdMarketCap: usdJson |> at(["band-protocol", "usd_market_cap"], JsonUtils.Decode.float),
        usd24HrChange:
          usdJson |> at(["band-protocol", "usd_24h_change"], JsonUtils.Decode.float),
        btcPrice: btcJson |> at(["band-protocol", "btc"], JsonUtils.Decode.float),
        btcMarketCap: btcJson |> at(["band-protocol", "btc_market_cap"], JsonUtils.Decode.float),
        btc24HrChange:
          btcJson |> at(["band-protocol", "btc_24h_change"], JsonUtils.Decode.float),
      };

    let usdJsonUrl = "https://api.coingecko.com/api/v3/simple/price?ids=band-protocol&vs_currencies=usd&include_market_cap=true&include_24hr_change=true";
    let btcJsonUrl = "https://api.coingecko.com/api/v3/simple/price?ids=band-protocol&vs_currencies=btc&include_market_cap=true&include_24hr_change=true";

    let (usdJson, usdReload) = AxiosHooks.useWithReload(usdJsonUrl);
    let (btcJson, btcReload) = AxiosHooks.useWithReload(btcJsonUrl);

    let reload = () => {
      usdReload((), ());
      btcReload((), ());
    };

    let data = {
      let%Opt usd = usdJson;
      let%Opt btc = btcJson;
      Some(decode(usd, btc));
    };

    (data, reload);
  };
};

module CrytoCompare = {
  let get = () => {
    // TODO: Find the formular to calulate this
    let circulatingSupply = 20494032.;
    let decode = (usdJson, btcJson) =>
      switch (
        JsonUtils.Decode.{
          usdPrice: usdJson |> at(["RAW", "BAND", "USD", "PRICE"], JsonUtils.Decode.float),
          usdMarketCap:
            (usdJson |> at(["RAW", "BAND", "USD", "PRICE"], JsonUtils.Decode.float))
            *. circulatingSupply,
          usd24HrChange:
            usdJson |> at(["RAW", "BAND", "USD", "CHANGEPCT24HOUR"], JsonUtils.Decode.float),
          btcPrice: btcJson |> at(["RAW", "BAND", "BTC", "PRICE"], JsonUtils.Decode.float),
          btcMarketCap:
            (btcJson |> at(["RAW", "BAND", "BTC", "PRICE"], JsonUtils.Decode.float))
            *. circulatingSupply,
          btc24HrChange:
            btcJson |> at(["RAW", "BAND", "BTC", "CHANGEPCT24HOUR"], JsonUtils.Decode.float),
        }
      ) {
      | result => Some(result)
      | exception _ => None
      };

    let usdJsonUrl = "https://min-api.cryptocompare.com/data/pricemultifull?fsyms=BAND&tsyms=USD";
    let btcJsonUrl = "https://min-api.cryptocompare.com/data/pricemultifull?fsyms=BAND&tsyms=BTC";

    let (usdJson, usdReload) = AxiosHooks.useWithReload(usdJsonUrl);
    let (btcJson, btcReload) = AxiosHooks.useWithReload(btcJsonUrl);

    let reload = () => {
      usdReload((), ());
      btcReload((), ());
    };

    let data = {
      let%Opt usd = usdJson;
      let%Opt btc = btcJson;
      let%Opt result = decode(usd, btc);
      Some(result);
    };

    (data, reload);
  };
};
