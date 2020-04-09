type t = {financial: PriceHook.Price.t};

let context = React.createContext(ContextHelper.default);

[@react.component]
let make = (~children) => {
  let financialOpt = PriceHook.get();

  let data = {
    let%Opt financial = financialOpt;
    Some({financial: financial});
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
