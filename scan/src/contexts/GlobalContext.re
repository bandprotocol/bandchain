type t = {
  financial: PriceHook.Price.t,
  latestBlock: BlockSub.t,
  latestBlocks: list(BlockSub.t),
  validators: list(ValidatorHook.Validator.t),
};

let context = React.createContext(ContextHelper.default);

[@react.component]
let make = (~children) => {
  let financialOpt = PriceHook.get();
  let validatorsOpt = ValidatorHook.getList();
  let latestBlocksOpt =
    switch (BlockSub.getList(~pageSize=10, ~page=1, ())) {
    | ApolloHooks.Subscription.Data(data) => Some(data)
    | _ => None
    };

  let data = {
    let%Opt financial = financialOpt;
    let%Opt latestBlocks = latestBlocksOpt;
    let%Opt latestBlock = latestBlocks->Belt_Array.get(0);
    let%Opt validators = validatorsOpt;
    Some({financial, latestBlock, latestBlocks: latestBlocks->Belt_List.fromArray, validators});
  };

  React.createElement(React.Context.provider(context), {"value": data, "children": children});
};
