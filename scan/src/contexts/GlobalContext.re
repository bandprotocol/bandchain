type t = {financial: PriceHook.Price.t};

let context = React.createContext(ContextHelper.default);

[@react.component]
let make = (~children) => {
  let (financialOpt, reload) = PriceHook.get();

  let data = {
    let%Opt financial = financialOpt;
    Some({financial: financial});
  };

  React.useEffect1(
    () => {
      let intervalID = Js.Global.setInterval(reload, 60000);
      Some(() => Js.Global.clearInterval(intervalID));
    },
    [|financialOpt|],
  );

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
