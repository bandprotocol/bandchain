type request_t =
  | Request(RequestSub.t)
  | RequestMini(RequestSub.Mini.t);

type request_packet_t = {
  clientID: string,
  oracleScriptID: int,
  calldata: JsBuffer.t,
  askCount: int,
  minCount: int,
};

type response_packet_t = {
  clientID: string,
  requestID: int,
  ansCount: int,
  requestTime: int,
  resolveTime: int,
  resolveStatus: int,
  result: JsBuffer.t,
};

type iavl_merkle_path_t = {
  isDataOnRight: bool,
  subtreeHeight: int,
  subtreeSize: int,
  subtreeVersion: int,
  siblingHash: JsBuffer.t,
};

type oracle_data_proof_t = {
  requestPacket: request_packet_t,
  responsePacket: response_packet_t,
  version: int,
  iavlMerklePaths: list(iavl_merkle_path_t),
};

type multi_store_proof_t = {
  accToGovStoresMerkleHash: JsBuffer.t,
  mainAndMintStoresMerkleHash: JsBuffer.t,
  oracleIAVLStateHash: JsBuffer.t,
  paramsStoresMerkleHash: JsBuffer.t,
  slashingToUpgradeStoresMerkleHash: JsBuffer.t,
};

type block_header_merkle_parts_t = {
  versionAndChainIdHash: JsBuffer.t,
  timeSecond: int,
  timeNanoSecond: int,
  lastBlockIDAndOther: JsBuffer.t,
  nextValidatorHashAndConsensusHash: JsBuffer.t,
  lastResultsHash: JsBuffer.t,
  evidenceAndProposerHash: JsBuffer.t,
};

type tm_signature_t = {
  r: JsBuffer.t,
  s: JsBuffer.t,
  v: int,
  signedDataPrefix: JsBuffer.t,
  signedDataSuffix: JsBuffer.t,
};

type block_relay_proof_t = {
  multiStoreProof: multi_store_proof_t,
  blockHeaderMerkleParts: block_header_merkle_parts_t,
  signatures: list(tm_signature_t),
};

type proof_t = {
  blockHeight: int,
  oracleDataProof: oracle_data_proof_t,
  blockRelayProof: block_relay_proof_t,
};

let decodeRequestPacket = (json: Js.Json.t) => {
  JsonUtils.Decode.{
    clientID: json |> optional(field("client_id", string)) |> Belt.Option.getWithDefault(_, ""),
    oracleScriptID: json |> field("oracle_script_id", intstr),
    calldata: json |> field("calldata", string) |> JsBuffer.fromBase64,
    askCount: json |> field("ask_count", intstr),
    minCount: json |> field("min_count", intstr),
  };
};

let decodeResponsePacket = json => {
  JsonUtils.Decode.{
    clientID: json |> optional(field("client_id", string)) |> Belt.Option.getWithDefault(_, ""),
    requestID: json |> field("request_id", intstr),
    ansCount: json |> field("ans_count", intstr),
    requestTime: json |> field("request_time", intstr),
    resolveTime: json |> field("resolve_time", intstr),
    resolveStatus: json |> field("resolve_status", int),
    result: json |> field("result", string) |> JsBuffer.fromBase64,
  };
};

let decodeIAVLMerklePath = json => {
  JsonUtils.Decode.{
    isDataOnRight: json |> field("isDataOnRight", bool),
    subtreeHeight: json |> field("subtreeHeight", int),
    subtreeSize: json |> field("subtreeSize", intstr),
    subtreeVersion: json |> field("subtreeVersion", intstr),
    siblingHash: json |> field("siblingHash", string) |> JsBuffer.fromHex,
  };
};

let decodeOracleDataProof = json => {
  JsonUtils.Decode.{
    requestPacket: json |> field("requestPacket", decodeRequestPacket),
    responsePacket: json |> field("responsePacket", decodeResponsePacket),
    version: json |> field("version", intstr),
    iavlMerklePaths: json |> field("merklePaths", list(decodeIAVLMerklePath)),
  };
};

let decodeMultiStoreProof = json => {
  JsonUtils.Decode.{
    accToGovStoresMerkleHash:
      json |> field("accToGovStoresMerkleHash", string) |> JsBuffer.fromHex,
    mainAndMintStoresMerkleHash:
      json |> field("mainAndMintStoresMerkleHash", string) |> JsBuffer.fromHex,
    oracleIAVLStateHash: json |> field("oracleIAVLStateHash", string) |> JsBuffer.fromHex,
    paramsStoresMerkleHash: json |> field("paramsStoresMerkleHash", string) |> JsBuffer.fromHex,
    slashingToUpgradeStoresMerkleHash:
      json |> field("slashingToUpgradeStoresMerkleHash", string) |> JsBuffer.fromHex,
  };
};

let decodeBlockHeaderMerkleParts = json => {
  JsonUtils.Decode.{
    versionAndChainIdHash: json |> field("versionAndChainIdHash", string) |> JsBuffer.fromHex,
    timeSecond: json |> field("timeSecond", intstr),
    timeNanoSecond: json |> field("timeNanoSecond", int),
    lastBlockIDAndOther: json |> field("lastBlockIDAndOther", string) |> JsBuffer.fromHex,
    nextValidatorHashAndConsensusHash:
      json |> field("nextValidatorHashAndConsensusHash", string) |> JsBuffer.fromHex,
    lastResultsHash: json |> field("lastResultsHash", string) |> JsBuffer.fromHex,
    evidenceAndProposerHash: json |> field("evidenceAndProposerHash", string) |> JsBuffer.fromHex,
  };
};

let decodeTMSignature = json => {
  JsonUtils.Decode.{
    r: json |> field("r", string) |> JsBuffer.fromHex,
    s: json |> field("s", string) |> JsBuffer.fromHex,
    v: json |> field("v", int),
    signedDataPrefix: json |> field("signedDataPrefix", string) |> JsBuffer.fromHex,
    signedDataSuffix: json |> field("signedDataSuffix", string) |> JsBuffer.fromHex,
  };
};

let decodeBlockRelayProof = json => {
  JsonUtils.Decode.{
    multiStoreProof: json |> field("multiStoreProof", decodeMultiStoreProof),
    blockHeaderMerkleParts: json |> field("blockHeaderMerkleParts", decodeBlockHeaderMerkleParts),
    signatures: json |> field("signatures", list(decodeTMSignature)),
  };
};

let decodeProof = json =>
  JsonUtils.Decode.{
    blockHeight: json |> field("blockHeight", intstr),
    oracleDataProof: json |> field("oracleDataProof", decodeOracleDataProof),
    blockRelayProof: json |> field("blockRelayProof", decodeBlockRelayProof),
  };

let obi_encode_int = (i, n) =>
  Obi.encode(
    "{x: " ++ n ++ "}/{_:u64}",
    "input",
    [|{fieldName: "x", fieldValue: i |> string_of_int}|],
  )
  |> Belt_Option.getExn;

type variant_of_proof_t =
  | RequestPacket(request_packet_t)
  | ResponsePacket(response_packet_t)
  | IAVLMerklePath(iavl_merkle_path_t)
  | IAVLMerklePaths(list(iavl_merkle_path_t))
  | MultiStoreProof(multi_store_proof_t)
  | BlockHeaderMerkleParts(block_header_merkle_parts_t)
  | Signature(tm_signature_t)
  | Signatures(list(tm_signature_t))
  | Proof(proof_t);

let rec encode =
  fun
  | RequestPacket({clientID, oracleScriptID, calldata, askCount, minCount}) => {
      Obi.encode(
        "{clientID: string, oracleScriptID: u64, calldata: bytes, askCount: u64, minCount: u64}/{_:u64}",
        "input",
        [|
          {fieldName: "clientID", fieldValue: clientID},
          {fieldName: "oracleScriptID", fieldValue: oracleScriptID |> string_of_int},
          {fieldName: "calldata", fieldValue: calldata |> JsBuffer.toHex(~with0x=true)},
          {fieldName: "askCount", fieldValue: askCount |> string_of_int},
          {fieldName: "minCount", fieldValue: minCount |> string_of_int},
        |],
      );
    }
  | ResponsePacket({
      clientID,
      requestID,
      ansCount,
      requestTime,
      resolveTime,
      resolveStatus,
      result,
    }) => {
      Obi.encode(
        "{clientID: string, requestID: u64, ansCount: u64, requestTime: u64, resolveTime: u64, resolveStatus: u32, result: bytes}/{_:u64}",
        "input",
        [|
          {fieldName: "clientID", fieldValue: clientID},
          {fieldName: "requestID", fieldValue: requestID |> string_of_int},
          {fieldName: "ansCount", fieldValue: ansCount |> string_of_int},
          {fieldName: "requestTime", fieldValue: requestTime |> string_of_int},
          {fieldName: "resolveTime", fieldValue: resolveTime |> string_of_int},
          {fieldName: "resolveStatus", fieldValue: resolveStatus |> string_of_int},
          {fieldName: "result", fieldValue: result |> JsBuffer.toHex(~with0x=true)},
        |],
      );
    }
  | IAVLMerklePath({isDataOnRight, subtreeHeight, subtreeSize, subtreeVersion, siblingHash}) => {
      Obi.encode(
        "{isDataOnRight: u8, subtreeHeight: u8, subtreeSize: u64, subtreeVersion: u64, siblingHash: bytes}/{_:u64}",
        "input",
        [|
          {fieldName: "isDataOnRight", fieldValue: isDataOnRight ? "1" : "0"},
          {fieldName: "subtreeHeight", fieldValue: subtreeHeight |> string_of_int},
          {fieldName: "subtreeSize", fieldValue: subtreeSize |> string_of_int},
          {fieldName: "subtreeVersion", fieldValue: subtreeVersion |> string_of_int},
          {fieldName: "siblingHash", fieldValue: siblingHash |> JsBuffer.toHex(~with0x=true)},
        |],
      );
    }
  | IAVLMerklePaths(iavl_merkle_paths) => {
      iavl_merkle_paths
      |> Belt_List.map(_, x => encode(IAVLMerklePath(x)))
      |> Belt_List.reduce(_, Some(JsBuffer.from([||])), (a, b) =>
           switch (a, b) {
           | (Some(acc), Some(elem)) => Some(JsBuffer.concat([|acc, elem|]))
           | _ => None
           }
         )
      |> Belt_Option.map(_, x =>
           JsBuffer.concat([|obi_encode_int(iavl_merkle_paths |> Belt_List.length, "u32"), x|])
         );
    }
  | MultiStoreProof({
      accToGovStoresMerkleHash,
      mainAndMintStoresMerkleHash,
      oracleIAVLStateHash,
      paramsStoresMerkleHash,
      slashingToUpgradeStoresMerkleHash,
    }) => {
      Some(
        JsBuffer.concat([|
          accToGovStoresMerkleHash,
          mainAndMintStoresMerkleHash,
          oracleIAVLStateHash,
          paramsStoresMerkleHash,
          slashingToUpgradeStoresMerkleHash,
        |]),
      );
    }
  | BlockHeaderMerkleParts({
      versionAndChainIdHash,
      timeSecond,
      timeNanoSecond,
      lastBlockIDAndOther,
      nextValidatorHashAndConsensusHash,
      lastResultsHash,
      evidenceAndProposerHash,
    }) => {
      Obi.encode(
        "{versionAndChainIdHash: bytes, timeSecond: u64, timeNanoSecond: u64, lastBlockIDAndOther: bytes, nextValidatorHashAndConsensusHash: bytes, lastResultsHash: bytes, evidenceAndProposerHash: bytes}/{_:u64}",
        "input",
        [|
          {
            fieldName: "versionAndChainIdHash",
            fieldValue: versionAndChainIdHash |> JsBuffer.toHex(~with0x=true),
          },
          {fieldName: "timeSecond", fieldValue: timeSecond |> string_of_int},
          {fieldName: "timeNanoSecond", fieldValue: timeNanoSecond |> string_of_int},
          {
            fieldName: "lastBlockIDAndOther",
            fieldValue: lastBlockIDAndOther |> JsBuffer.toHex(~with0x=true),
          },
          {
            fieldName: "nextValidatorHashAndConsensusHash",
            fieldValue: nextValidatorHashAndConsensusHash |> JsBuffer.toHex(~with0x=true),
          },
          {
            fieldName: "lastResultsHash",
            fieldValue: lastResultsHash |> JsBuffer.toHex(~with0x=true),
          },
          {
            fieldName: "evidenceAndProposerHash",
            fieldValue: evidenceAndProposerHash |> JsBuffer.toHex(~with0x=true),
          },
        |],
      );
    }
  | Signature({r, s, v, signedDataPrefix, signedDataSuffix}) => {
      Obi.encode(
        "{r: bytes, s: bytes, v: u8, signedDataPrefix: bytes, signedDataSuffix: bytes}/{_:u64}",
        "input",
        [|
          {fieldName: "r", fieldValue: r |> JsBuffer.toHex(~with0x=true)},
          {fieldName: "s", fieldValue: s |> JsBuffer.toHex(~with0x=true)},
          {fieldName: "v", fieldValue: v |> string_of_int},
          {
            fieldName: "signedDataPrefix",
            fieldValue: signedDataPrefix |> JsBuffer.toHex(~with0x=true),
          },
          {
            fieldName: "signedDataSuffix",
            fieldValue: signedDataSuffix |> JsBuffer.toHex(~with0x=true),
          },
        |],
      );
    }
  | Signatures(tm_signatures) => {
      tm_signatures
      |> Belt_List.map(_, x => encode(Signature(x)))
      |> Belt_List.reduce(_, Some(JsBuffer.from([||])), (a, b) =>
           switch (a, b) {
           | (Some(acc), Some(elem)) => Some(JsBuffer.concat([|acc, elem|]))
           | _ => None
           }
         )
      |> Belt_Option.map(_, x =>
           JsBuffer.concat([|obi_encode_int(tm_signatures |> Belt_List.length, "u32"), x|])
         );
    }
  | Proof({
      blockHeight,
      oracleDataProof: {requestPacket, responsePacket, version, iavlMerklePaths},
      blockRelayProof: {multiStoreProof, blockHeaderMerkleParts, signatures},
    }) => {
      let%Opt encodeMultiStore = encode(MultiStoreProof(multiStoreProof));
      let%Opt encodeBlockHeaderMerkleParts =
        encode(BlockHeaderMerkleParts(blockHeaderMerkleParts));
      let%Opt encodeSignatures = encode(Signatures(signatures));
      let%Opt encodeReq = encode(RequestPacket(requestPacket));
      let%Opt encodeRes = encode(ResponsePacket(responsePacket));
      let%Opt encodeIAVLMerklePaths = encode(IAVLMerklePaths(iavlMerklePaths));
      Obi.encode(
        "{blockHeight: u64, multiStore: bytes, blockMerkleParts: bytes, signatures: bytes, packet: bytes, version: u64, iavlPaths: bytes}/{_:u64}",
        "input",
        [|
          {fieldName: "blockHeight", fieldValue: blockHeight |> string_of_int},
          {
            fieldName: "multiStore",
            fieldValue: encodeMultiStore |> JsBuffer.toHex(~with0x=true),
          },
          {
            fieldName: "blockMerkleParts",
            fieldValue: encodeBlockHeaderMerkleParts |> JsBuffer.toHex(~with0x=true),
          },
          {
            fieldName: "signatures",
            fieldValue: encodeSignatures |> JsBuffer.toHex(~with0x=true),
          },
          {
            fieldName: "packet",
            fieldValue:
              JsBuffer.concat([|encodeReq, encodeRes|]) |> JsBuffer.toHex(~with0x=true),
          },
          {fieldName: "version", fieldValue: version |> string_of_int},
          {
            fieldName: "iavlPaths",
            fieldValue: encodeIAVLMerklePaths |> JsBuffer.toHex(~with0x=true),
          },
        |],
      );
    };

let createProofFromJson = (proof: Js.Json.t) => {
  switch (Proof(proof |> decodeProof)) {
  | result => result |> encode
  | exception _ => None
  };
};
