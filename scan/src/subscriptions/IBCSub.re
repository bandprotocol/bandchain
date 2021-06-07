module Request = {
  type t = {
    idOpt: option(ID.Request.t),
    oracleScriptID: ID.OracleScript.t,
    oracleScriptNameOpt: option(string),
    clientID: string,
    calldata: JsBuffer.t,
    askCount: int,
    minCount: int,
  };
};

module Response = {
  type t = {
    requestID: ID.Request.t,
    oracleScriptID: ID.OracleScript.t,
    oracleScriptName: string,
    status: RequestStatus.t,
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
  srcChannel: string,
  srcPort: string,
  chainID: string,
  dstChannel: string,
  dstPort: string,
  blockHeight: ID.Block.t,
  packet: packet_t,
};

module Internal = {
  type t = {
    isIncoming: bool,
    blockHeight: ID.Block.t,
    srcChannel: string,
    srcPort: string,
    dstChannel: string,
    dstPort: string,
    packetType: string,
    packetDetail: Js.Json.t,
    acknowledgement: option(Js.Json.t),
  };

  let toExternal =
      (
        {
          isIncoming,
          blockHeight,
          srcChannel,
          srcPort,
          dstChannel,
          dstPort,
          packetType,
          packetDetail,
          acknowledgement,
        },
      ) => {
    direction: isIncoming ? Incoming : Outgoing,
    srcChannel,
    srcPort,
    chainID: "bandchain",
    dstChannel,
    dstPort,
    blockHeight,
    packet:
      switch (packetType) {
      | "oracle request" =>
        Request(
          JsonUtils.Decode.{
            idOpt: {
              let x = acknowledgement->Belt.Option.getExn;
              let success = x |> field("success", bool);
              success ? Some(ID.Request.ID(x |> at(["request_id"], int))) : None;
            },
            oracleScriptID: ID.OracleScript.ID(packetDetail |> at(["oracle_script_id"], int)),
            oracleScriptNameOpt: packetDetail |> optional(at(["oracle_script_name"], string)),
            clientID: packetDetail |> at(["client_id"], string),
            calldata: packetDetail |> at(["calldata"], string) |> JsBuffer.fromHex,
            askCount: packetDetail |> at(["ask_count"], int),
            minCount: packetDetail |> at(["min_count"], int),
          },
        )
      | "oracle response" =>
        let status =
          packetDetail
          |> JsonUtils.Decode.at(["resolve_status"], JsonUtils.Decode.int)
          |> RequestStatus.fromInt;
        Response(
          JsonUtils.Decode.{
            requestID: ID.Request.ID(packetDetail |> at(["request_id"], int)),
            oracleScriptID: ID.OracleScript.ID(packetDetail |> at(["oracle_script_id"], int)),
            oracleScriptName: packetDetail |> at(["oracle_script_name"], string),
            status,
            result:
              status == Success
                ? Some(packetDetail |> at(["result"], string) |> JsBuffer.fromHex) : None,
          },
        );
      | _ => Unknown
      },
  };
  module MultiPacketsConfig = [%graphql
    {|
    subscription Packets($limit: Int!, $offset: Int!) {
      packets(limit: $limit, offset: $offset, order_by: {block_height: desc}) @bsRecord {
        isIncoming: is_incoming
        blockHeight: block_height @bsDecoder(fn: "ID.Block.fromInt")
        srcChannel: src_channel
        srcPort: src_port
        dstChannel: dst_channel
        dstPort: dst_port
        packetType: type
        packetDetail: data
        acknowledgement
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

let getList = (~page=1, ~pageSize=10, ()): ApolloHooks.Subscription.variant(array(t)) => {
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
