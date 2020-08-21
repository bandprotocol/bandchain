type t('a) = ApolloHooks.Query.variant('a);

let map = (result, f) =>
  switch (result) {
  | ApolloHooks.Query.Data(data) => ApolloHooks.Query.Data(data |> f)
  | Loading => Loading
  | Error(e) => Error(e)
  | NoData => NoData
  };

let resolve = data => ApolloHooks.Query.Data(data);

let default = (result, value) =>
  switch (result) {
  | ApolloHooks.Query.Data(data) => data
  | _ => value
  };

let let_ = (result, f) =>
  switch (result) {
  | ApolloHooks.Query.Data(data) => f(data)
  | Loading => ApolloHooks.Query.Loading
  | Error(e) => Error(e)
  | NoData => NoData
  };

let all2 = (s1, s2) => let_(s1, s1' => let_(s2, s2' => Data((s1', s2'))));

let all3 = (s1, s2, s3) =>
  let_(s1, s1' => let_(s2, s2' => let_(s3, s3' => Data((s1', s2', s3')))));

let all4 = (s1, s2, s3, s4) =>
  let_(s1, s1' => let_(s2, s2' => let_(s3, s3' => let_(s4, s4' => Data((s1', s2', s3', s4'))))));

let all5 = (s1, s2, s3, s4, s5) =>
  let_(s1, s1' =>
    let_(s2, s2' =>
      let_(s3, s3' => let_(s4, s4' => let_(s5, s5' => Data((s1', s2', s3', s4', s5')))))
    )
  );
