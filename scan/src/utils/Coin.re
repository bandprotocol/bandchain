type t = {
  denom: string,
  amount: float,
};

let decodeCoin = json =>
  JsonUtils.Decode.{
    denom: json |> field("denom", string),
    amount: json |> field("amount", uamount),
  };

let newCoin = (denom, amount) => {denom, amount};

let getBandAmountFromCoins = coins =>
  coins
  ->Belt_List.keep(coin => coin.denom == "uband")
  ->Belt_List.get(0)
  ->Belt_Option.mapWithDefault(0., coin => coin.amount /. 1e6);

let getDescription = coin => {
  (coin.amount |> Format.fPretty)
  ++ " "
  ++ (
    switch (coin.denom.[0]) {
    | 'u' =>
      coin.denom->String.sub(_, 1, (coin.denom |> String.length) - 1) |> String.uppercase_ascii
    | _ => coin.denom
    }
  );
};

let toCoinsString = coins => {
  coins
  ->Belt_List.map(coin => coin->getDescription)
  ->Belt_List.reduceWithIndex("", (des, acc, i) =>
      acc ++ des ++ (i + 1 < coins->Belt_List.size ? ", " : "")
    );
};

let getFeeAmount = coins => {
  let coinOpt = coins->Belt_List.get(0);
  switch (coinOpt) {
  | Some(coin) => coin.amount
  | None => 0.
  };
};
