type t = {codeText: option(string)};

module Config = [%graphql
  {|
    subscription OracleScriptCode($codeHash: bytea!) {
        oracle_script_codes_by_pk(code_hash: $codeHash) @bsRecord {
            codeText: code_text
        }
    }
  |}
];

let get = codeHash => {
  let (result, _) =
    ApolloHooks.useSubscription(
      Config.definition,
      ~variables=
        Config.makeVariables(
          ~codeHash=codeHash |> Hash.toHex |> (x => "\\x" ++ x) |> Js.Json.string,
          (),
        ),
    );
  let%Sub x = result;
  switch (x##oracle_script_codes_by_pk) {
  | Some(data) => Sub.resolve(data)
  | None => NoData
  };
};
