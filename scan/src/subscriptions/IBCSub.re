module Request = {
  type t = {
    id: ID.Request.t,
    oracleScriptID: ID.OracleScript.t,
    oracleScriptName: string,
    calldata: JsBuffer.t,
    requestedValidatorCount: int,
    sufficientValidatorCount: int,
    expiration: int,
    prepareGas: int,
    executeGas: int,
  };
};

module Response = {
  type status_t =
    | Success
    | Fail;

  type t = {
    requestID: ID.Request.t,
    oracleScriptID: ID.OracleScript.t,
    oracleScriptName: string,
    status: status_t,
    result: option(JsBuffer.t),
  };
};

type packet_direction_t =
  | Incoming
  | Outgoing;

type packet_t =
  | Unknown
  | Request(Request.t)
  | Response(Response.t);

type t = {
  direction: packet_direction_t,
  chainID: string,
  channel: string,
  port: string,
  yourChainID: string,
  yourChannel: string,
  yourPort: string,
  blockHeight: ID.Block.t,
  packet: packet_t,
};

module Internal = {
  type t = {
    isIncoming: bool,
    blockHeight: ID.Block.t,
    channel: string,
    port: string,
    yourChainID: string,
    yourChannel: string,
    yourPort: string,
    packetType: string,
    packetDetail: Js.Json.t,
  };

  let toExternal =
      (
        {
          isIncoming,
          blockHeight,
          channel,
          port,
          yourChainID,
          yourChannel,
          yourPort,
          packetType,
          packetDetail,
        },
      ) => {
    direction: isIncoming ? Incoming : Outgoing,
    chainID: "bandchain",
    channel,
    port,
    yourChainID,
    yourChannel,
    yourPort,
    blockHeight,
    packet:
      switch (packetType) {
      | "ORACLE REQUEST" =>
        Request(
          JsonUtils.Decode.{
            id: ID.Request.ID(packetDetail |> at(["extra", "requestID"], int)),
            oracleScriptID:
              ID.OracleScript.ID(packetDetail |> at(["value", "oracleScriptID"], intstr)),
            oracleScriptName: packetDetail |> at(["extra", "oracleScriptName"], string),
            calldata: packetDetail |> at(["value", "calldata"], string) |> JsBuffer.fromHex,
            requestedValidatorCount:
              packetDetail |> at(["value", "requestedValidatorCount"], intstr),
            sufficientValidatorCount:
              packetDetail |> at(["value", "sufficientValidatorCount"], intstr),
            expiration: packetDetail |> at(["value", "expiration"], intstr),
            prepareGas: packetDetail |> at(["value", "prepareGas"], intstr),
            executeGas: packetDetail |> at(["value", "executeGas"], intstr),
          },
        )
      | _ =>
        let status =
          packetDetail
          |> JsonUtils.Decode.at(["extra", "resolveStatus"], JsonUtils.Decode.string)
          == "Success"
            ? Response.Success : Response.Fail;
        Response(
          JsonUtils.Decode.{
            requestID: ID.Request.ID(packetDetail |> at(["value", "requestID"], intstr)),
            oracleScriptID:
              ID.OracleScript.ID(packetDetail |> at(["extra", "oracleScriptID"], int)),
            oracleScriptName: packetDetail |> at(["extra", "oracleScriptName"], string),
            status,
            result:
              status == Success
                ? Some(packetDetail |> at(["value", "result"], string) |> JsBuffer.fromHex)
                : None,
          },
        );
      },
  };

  module MultiPacketsConfig = [%graphql
    {|
    subscription Requests($limit: Int!, $offset: Int!) {
      packets(limit: $limit, offset: $offset, order_by: {block_height: desc}) @bsRecord {
        isIncoming: is_incoming
        blockHeight: block_height @bsDecoder(fn: "ID.Block.fromJson")
        channel: my_channel
        port: my_port
        yourChainID: your_chain_id
        yourChannel: your_channel
        yourPort: your_port
        packetType: type
        packetDetail: detail
      }
    }
  |}
  ];
};

module PacketCountConfig = [%graphql
  {|
  subscription PacketsCount {
    packets_aggregate{
      aggregate{
        count @bsDecoder(fn: "Belt_Option.getExn")
      }
    }
  }
|}
];

let getList = (~page=1, ~pageSize=10, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useSubscription(
      Internal.MultiPacketsConfig.definition,
      ~variables=Internal.MultiPacketsConfig.makeVariables(~limit=pageSize, ~offset, ()),
    );
  result |> Sub.map(_, x => x##packets->Belt_Array.map(Internal.toExternal));
};

let count = () => {
  let (result, _) = ApolloHooks.useSubscription(PacketCountConfig.definition);
  result
  |> Sub.map(_, x => x##packets_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count));
};
