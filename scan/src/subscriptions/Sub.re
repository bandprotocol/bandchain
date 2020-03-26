let map = (result, f) =>
  switch (result) {
  | ApolloHooks.Subscription.Data(data) => ApolloHooks.Subscription.Data(data |> f)
  | Loading => Loading
  | Error(e) => Error(e)
  | NoData => NoData
  };

let resolve = data => ApolloHooks.Subscription.Data(data);

let default = (result, value) =>
  switch (result) {
  | ApolloHooks.Subscription.Data(data) => data
  | _ => value
  };

let toOption =
  fun
  | ApolloHooks.Subscription.Data(data) => Some(data)
  | _ => None;

let let_ = (result, f) =>
  switch (result) {
  | ApolloHooks.Subscription.Data(data) => f(data)
  | Loading => ApolloHooks.Subscription.Loading
  | Error(e) => Error(e)
  | NoData => NoData
  };
