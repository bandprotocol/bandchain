module ValidatorReport = {
  type oracle_script_t = {
    oracleScriptID: ID.OracleScript.t,
    name: string,
  };

  type request_t = {
    id: ID.Request.t,
    oracleScript: oracle_script_t,
  };

  type data_source_t = {
    dataSourceID: ID.DataSource.t,
    dataSourceName: string,
  };

  type raw_request_t = {
    calldata: JsBuffer.t,
    dataSource: data_source_t,
  };

  type report_details_t = {
    externalID: string,
    exitCode: string,
    data: JsBuffer.t,
    rawRequest: option(raw_request_t),
  };

  type transaction_t = {hash: Hash.t};

  type internal_t = {
    request: request_t,
    transaction: option(transaction_t),
    reportDetails: array(report_details_t),
  };

  type t = {
    txHash: option(Hash.t),
    request: request_t,
    reportDetails: array(report_details_t),
  };

  let toExternal = ({request, transaction, reportDetails}) => {
    txHash: transaction->Belt.Option.map(({hash}) => hash),
    request,
    reportDetails,
  };

  module MultiConfig = [%graphql
    {|
      subscription Reports ($limit: Int!, $offset: Int!, $validator: String!) {
        validators_by_pk(operator_address: $validator) {
          reports (limit: $limit, offset: $offset, order_by: {request_id: desc}) @bsRecord {
              request @bsRecord {
                id @bsDecoder (fn: "ID.Request.fromInt")
                oracleScript: oracle_script @bsRecord {
                  oracleScriptID: id @bsDecoder (fn: "ID.OracleScript.fromInt")
                  name
                }
              }
              transaction @bsRecord{
                hash @bsDecoder (fn: "GraphQLParser.hash")
              }
              reportDetails: raw_reports @bsRecord {
                externalID: external_id @bsDecoder (fn:"GraphQLParser.string")
                exitCode: exit_code  @bsDecoder (fn:"GraphQLParser.string")
                data @bsDecoder (fn: "GraphQLParser.buffer")
                rawRequest: raw_request @bsRecord {
                  calldata @bsDecoder(fn: "GraphQLParser.buffer")
                  dataSource: data_source @bsRecord {
                    dataSourceID: id @bsDecoder (fn: "ID.DataSource.fromInt")
                    dataSourceName: name
                  }
                }
              }
            }
          }
        }
      |}
  ];

  module ReportCountConfig = [%graphql
    {|
    subscription ReportsCount ($validator: String!) {
      validators_by_pk(operator_address: $validator) {
        reports_aggregate {
          aggregate{
            count @bsDecoder(fn: "Belt_Option.getExn")
          }
        }
      }
    }
  |}
  ];

  let getListByValidator = (~page=1, ~pageSize=5, ~validator) => {
    let offset = (page - 1) * pageSize;
    let (result, _) =
      ApolloHooks.useSubscription(
        MultiConfig.definition,
        ~variables=MultiConfig.makeVariables(~limit=pageSize, ~offset, ~validator, ()),
      );
    result
    |> Sub.map(_, x => {
         switch (x##validators_by_pk) {
         | Some(x') => x'##reports->Belt_Array.map(toExternal)
         | None => [||]
         }
       });
  };

  let count = validator => {
    let (result, _) =
      ApolloHooks.useSubscription(
        ReportCountConfig.definition,
        ~variables=ReportCountConfig.makeVariables(~validator, ()),
      );
    result
    |> Sub.map(_, x => {
         switch (x##validators_by_pk) {
         | Some(x') => x'##reports_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count)
         | None => 0
         }
       });
  };
};
