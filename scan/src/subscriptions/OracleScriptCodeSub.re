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

// TODO: If we can get the schema out of IBCSub directly then this module is not necessary any more.
module SchemaByOracleScriptIDConfig = [%graphql
  {|
    subscription OracleScriptCode($oracleScriptID: bigint!) {
      oracle_script_codes(where: {oracle_scripts: {id: {_eq: $oracleScriptID}}}) {
        schema
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

// TODO: If we can get the schema out of IBCSub directly then this function is not necessary any more.
let getSchemaByOracleScriptID = oracleScriptID => {
  let (result, _) =
    ApolloHooks.useSubscription(
      SchemaByOracleScriptIDConfig.definition,
      ~variables=
        SchemaByOracleScriptIDConfig.makeVariables(
          ~oracleScriptID=oracleScriptID |> ID.OracleScript.toJson,
          (),
        ),
    );
  let%Sub x = result;
  x##oracle_script_codes
  ->Belt_Array.get(0)
  ->Belt_Option.mapWithDefault(ApolloHooks.Subscription.NoData, y =>
      switch (y##schema) {
      | Some(schema) => Sub.resolve(schema)
      | None => NoData
      }
    );
};
