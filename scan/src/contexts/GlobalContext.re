type t = {
  financial: PriceHook.Price.t,
  validators: list(ValidatorHook.Validator.t),
};

let context = React.createContext(ContextHelper.default);

[@react.component]
let make = (~children) => {
  let financialOpt = PriceHook.get();
  let validatorsOpt = ValidatorHook.getList();

  let data = {
    let%Opt financial = financialOpt;
    let%Opt validators = validatorsOpt;
    Some({financial, validators});
  };

  React.createElement(React.Context.provider(context), {"value": data, "children": children});
};
