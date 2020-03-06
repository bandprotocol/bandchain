type t = {
  financial: PriceHook.Price.t,
  latestBlock: BlockHook.Block.t,
  latestBlocks: list(BlockHook.Block.t),
  validators: list(ValidatorHook.Validator.t),
};

let context = React.createContext(ContextHelper.default);

[@react.component]
let make = (~children) => {
  let financialOpt = PriceHook.get();
  let latestBlocksOpt = BlockHook.latest();
  let validatorsOpt = ValidatorHook.getValidators();
  let data = {
    let%Opt financial = financialOpt;
    let%Opt latestBlocks = latestBlocksOpt;
    let%Opt latestBlock = latestBlocks->Belt.List.get(0);
    let%Opt validators = validatorsOpt;
    Some({financial, latestBlock, latestBlocks, validators});
  };

  React.createElement(React.Context.provider(context), {"value": data, "children": children});
};
