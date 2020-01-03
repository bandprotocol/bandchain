type t = {
  financial: PriceHook.Price.t,
  latestBlock: BlockHook.Block.t,
  latestBlocks: list(BlockHook.Block.t),
};

let context = React.createContext(ContextHelper.default);

[@react.component]
let make = (~children) => {
  let financialOpt = PriceHook.get(~pollInterval=300000, ()); // 5-min
  let latestBlocksOpt = BlockHook.latest(~pollInterval=3000, ()); // 3-sec

  let data = {
    let%Opt financial = financialOpt;
    let%Opt latestBlocks = latestBlocksOpt;
    let%Opt latestBlock = latestBlocks->Belt.List.get(0);
    Some({financial, latestBlock, latestBlocks});
  };

  React.createElement(React.Context.provider(context), {"value": data, "children": children});
};
