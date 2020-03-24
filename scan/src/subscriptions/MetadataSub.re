type t = {
  lastProcessedHeight: ID.Block.t,
  chainID: string,
};

module Config = [%graphql
  {|
  subscription Metadata {
    metadata {
      key
      value
    }
  }
|}
];

let find = (arr, field) => {
  arr->Belt.Array.keepMap(a => {a##key == field ? Some(a##value) : None})[0];
};

let decode = raw => {
  let metadata = raw##metadata;
  {
    lastProcessedHeight:
      ID.Block.ID(metadata |> find(_, "last_processed_height") |> int_of_string),
    chainID: metadata |> find(_, "chain_id"),
  };
};

let use = () => {
  let (result, _) = ApolloHooks.useSubscription(Config.definition);
  result |> Sub.map(_, decode);
};
