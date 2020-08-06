module Price = {
  type t = {
    usdPrice: float,
    usdMarketCap: float,
    usd24HrChange: float,
    btcPrice: float,
    btcMarketCap: float,
    btc24HrChange: float,
    // circulatingSupply: float,
  };
  let decode = (usdJson, btcJson) =>
    JsonUtils.Decode.{
      usdPrice: usdJson |> at(["RAW", "BAND", "USD", "PRICE"], JsonUtils.Decode.float),
      usdMarketCap: usdJson |> at(["RAW", "BAND", "USD", "MKTCAP"], JsonUtils.Decode.float),
      usd24HrChange:
        usdJson |> at(["RAW", "BAND", "USD", "CHANGEPCT24HOUR"], JsonUtils.Decode.float),
      btcPrice: btcJson |> at(["RAW", "BAND", "BTC", "PRICE"], JsonUtils.Decode.float),
      btcMarketCap: btcJson |> at(["RAW", "BAND", "BTC", "MKTCAP"], JsonUtils.Decode.float),
      btc24HrChange:
        btcJson |> at(["RAW", "BAND", "BTC", "CHANGEPCT24HOUR"], JsonUtils.Decode.float),
      // circulatingSupply: btcJson |> at(["RAW", "BAND", "BTC", "CHANGEPCT24HOUR"], JsonUtils.Decode.float),
    };
};
let get = () => {
  let (usdJson, usdReload) =
    AxiosHooks.useWithReload(
      "https://min-api.cryptocompare.com/data/pricemultifull?fsyms=BAND&tsyms=USD",
    );
  let (btcJson, btcReload) =
    AxiosHooks.useWithReload(
      "https://min-api.cryptocompare.com/data/pricemultifull?fsyms=BAND&tsyms=BTC",
    );

  let reload = () => {
    usdReload((), ());
    btcReload((), ());
  };

  let data = {
    let%Opt usd = usdJson;
    let%Opt btc = btcJson;
    Some(Price.decode(usd, btc));
  };

  (data, reload);
};
