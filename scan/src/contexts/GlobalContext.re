type t = {financial: PriceHook.t};

let context = React.createContext(ContextHelper.default);

[@react.component]
let make = (~children) => {
  let (financialOptCG, reloadCG) = PriceHook.CoinGekco.get();
  let (financialOptCC, reloadCC) = PriceHook.CrytoCompare.get();

  let reload = () => {
    reloadCG();
    reloadCC();
  };

  React.useEffect0(() => {
    let intervalID = Js.Global.setInterval(reload, 60000);
    Some(() => Js.Global.clearInterval(intervalID));
  });

  let data = {
    switch (financialOptCG, financialOptCC) {
    | (Some(financial), _) => Some({financial: financial})
    | (_, Some(financial)) => Some({financial: financial})
    | (_, _) => None
    };
  };

  React.createElement(
    React.Context.provider(context),
    {
      "value":
        switch (data) {
        | Some(x) => Sub.resolve(x)
        | None => Loading
        },
      "children": children,
    },
  );
};
