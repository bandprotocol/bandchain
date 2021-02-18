type t = {financial: PriceHook.t};

let context = React.createContext(ContextHelper.default);

[@react.component]
let make = (~children) => {
  let client = BandChainJS.createClient("https://api-gm-lb.bandchain.org");
  let (financialOpt, setFinancialOpt) = React.useState(_ => None);

  React.useEffect0(() => {
    let fetchData = () => {
      let _ =
        PriceHook.getBandInfo(client)
        |> Js.Promise.then_(bandInfoOpt => {
             setFinancialOpt(_ => bandInfoOpt);
             Promise.ret();
           });
      ();
    };

    fetchData();
    let intervalID = Js.Global.setInterval(fetchData, 60000);
    Some(() => Js.Global.clearInterval(intervalID));
  });

  let data = {
    switch (financialOpt) {
    | Some(financial) => Some({financial: financial})
    | _ => None
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
