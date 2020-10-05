type t = {
  usdPrice: float,
  usdMarketCap: float,
  usd24HrChange: float,
  btcPrice: float,
  btcMarketCap: float,
  btc24HrChange: float,
  circulatingSupply: float,
};

module CoinGekco = {
  let get = () => {
    let decode = (usdJson, btcJson, supplyJson) =>
      JsonUtils.Decode.{
        usdPrice: usdJson |> at(["band-protocol", "usd"], JsonUtils.Decode.float),
        usdMarketCap: usdJson |> at(["band-protocol", "usd_market_cap"], JsonUtils.Decode.float),
        usd24HrChange:
          usdJson |> at(["band-protocol", "usd_24h_change"], JsonUtils.Decode.float),
        btcPrice: btcJson |> at(["band-protocol", "btc"], JsonUtils.Decode.float),
        btcMarketCap: btcJson |> at(["band-protocol", "btc_market_cap"], JsonUtils.Decode.float),
        btc24HrChange:
          btcJson |> at(["band-protocol", "btc_24h_change"], JsonUtils.Decode.float),
        circulatingSupply: supplyJson,
      };

    let usdJsonUrl = "https://api.coingecko.com/api/v3/simple/price?ids=band-protocol&vs_currencies=usd&include_market_cap=true&include_24hr_change=true";
    let btcJsonUrl = "https://api.coingecko.com/api/v3/simple/price?ids=band-protocol&vs_currencies=btc&include_market_cap=true&include_24hr_change=true";
    let supplyJsonUrl = "https://supply.bandchain.org/circulating";

    let (usdJson, usdReload) = AxiosHooks.useWithReload(usdJsonUrl);
    let (btcJson, btcReload) = AxiosHooks.useWithReload(btcJsonUrl);
    let (supplyJson, supplyReload) = AxiosHooks.useWithReload(supplyJsonUrl);

    let reload = () => {
      usdReload((), ());
      btcReload((), ());
      supplyReload((), ());
    };

    let data = {
      let%Opt usd = usdJson;
      let%Opt btc = btcJson;
      let%Opt supply =
        switch (supplyJson) {
        | None => None
        | Some(sp) => sp |> Js.Json.decodeNumber
        };
      Some(decode(usd, btc, supply));
    };

    (data, reload);
  };
};

module CrytoCompare = {
  let get = () => {
    // TODO: Find the formular to calulate this
    let decode = (usdJson, btcJson, supplyJson) =>
      switch (
        JsonUtils.Decode.{
          usdPrice: usdJson |> at(["RAW", "BAND", "USD", "PRICE"], JsonUtils.Decode.float),
          usdMarketCap:
            (usdJson |> at(["RAW", "BAND", "USD", "PRICE"], JsonUtils.Decode.float))
            *. supplyJson,
          usd24HrChange:
            usdJson |> at(["RAW", "BAND", "USD", "CHANGEPCT24HOUR"], JsonUtils.Decode.float),
          btcPrice: btcJson |> at(["RAW", "BAND", "BTC", "PRICE"], JsonUtils.Decode.float),
          btcMarketCap:
            (btcJson |> at(["RAW", "BAND", "BTC", "PRICE"], JsonUtils.Decode.float))
            *. supplyJson,
          btc24HrChange:
            btcJson |> at(["RAW", "BAND", "BTC", "CHANGEPCT24HOUR"], JsonUtils.Decode.float),
          circulatingSupply: supplyJson,
        }
      ) {
      | result => Some(result)
      | exception _ => None
      };

    let usdJsonUrl = "https://min-api.cryptocompare.com/data/pricemultifull?fsyms=BAND&tsyms=USD";
    let btcJsonUrl = "https://min-api.cryptocompare.com/data/pricemultifull?fsyms=BAND&tsyms=BTC";
    let supplyJsonUrl = "https://supply.bandchain.org/circulating";

    let (usdJson, usdReload) = AxiosHooks.useWithReload(usdJsonUrl);
    let (btcJson, btcReload) = AxiosHooks.useWithReload(btcJsonUrl);
    let (supplyJson, supplyReload) = AxiosHooks.useWithReload(supplyJsonUrl);

    let reload = () => {
      usdReload((), ());
      btcReload((), ());
      supplyReload((), ());
    };

    let data = {
      let%Opt usd = usdJson;
      let%Opt btc = btcJson;
      let%Opt supply =
        switch (supplyJson) {
        | None => None
        | Some(sp) => sp |> Js.Json.decodeNumber
        };
      let%Opt result = decode(usd, btc, supply);
      Some(result);
    };

    (data, reload);
  };
};
