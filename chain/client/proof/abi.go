package proof

var relayFormat = []byte(`
[
  {
    "internalType": "uint256",
    "name": "_blockHeight",
    "type": "uint256"
  },
  {
    "components": [
      {
        "internalType": "bytes32",
        "name": "accToMemCapStoresMerkleHash",
        "type": "bytes32"
      },
      {
        "internalType": "bytes32",
        "name": "mintStoresMerkleHash",
        "type": "bytes32"
      },
      {
        "internalType": "bytes32",
        "name": "oracleIAVLStateHash",
        "type": "bytes32"
      },
      {
        "internalType": "bytes32",
        "name": "paramsAndSlashingStoresMerkleHash",
        "type": "bytes32"
      },
      {
        "internalType": "bytes32",
        "name": "stakingAndUpgradeStoresMerkleHash",
        "type": "bytes32"
      }
    ],
    "internalType": "struct MultiStore.Data",
    "name": "_multiStore",
    "type": "tuple"
  },
  {
    "components": [
      {
        "internalType": "bytes32",
        "name": "versionAndChainIdHash",
        "type": "bytes32"
      },
      {
        "internalType": "bytes32",
        "name": "timeHash",
        "type": "bytes32"
      },
      {
        "internalType": "bytes32",
        "name": "lastBlockIDAndOther",
        "type": "bytes32"
      },
      {
        "internalType": "bytes32",
        "name": "nextValidatorHashAndConsensusHash",
        "type": "bytes32"
      },
      {
        "internalType": "bytes32",
        "name": "lastResultsHash",
        "type": "bytes32"
      },
      {
        "internalType": "bytes32",
        "name": "evidenceAndProposerHash",
        "type": "bytes32"
      }
    ],
    "internalType": "struct BlockHeaderMerkleParts.Data",
    "name": "_merkleParts",
    "type": "tuple"
  },
  {
    "internalType": "bytes",
    "name": "_signedDataPrefix",
    "type": "bytes"
  },
  {
    "components": [
      {
        "internalType": "bytes32",
        "name": "r",
        "type": "bytes32"
      },
      {
        "internalType": "bytes32",
        "name": "s",
        "type": "bytes32"
      },
      {
        "internalType": "uint8",
        "name": "v",
        "type": "uint8"
      },
      {
        "internalType": "bytes",
        "name": "signedDataSuffix",
        "type": "bytes"
      }
    ],
    "internalType": "struct TMSignature.Data[]",
    "name": "_signatures",
    "type": "tuple[]"
  }
]
`)

var verifyFormat = []byte(`
[
  {
    "internalType": "uint256",
    "name": "_blockHeight",
    "type": "uint256"
  },
  {
    "components": [
      {
        "internalType": "string",
        "name": "clientId",
        "type": "string"
      },
      {
        "internalType": "uint64",
        "name": "oracleScriptId",
        "type": "uint64"
      },
      {
        "internalType": "string",
        "name": "params",
        "type": "string"
      },
      {
        "internalType": "uint64",
        "name": "askCount",
        "type": "uint64"
      },
      {
        "internalType": "uint64",
        "name": "minCount",
        "type": "uint64"
      }
    ],
    "internalType": "struct IBridge.RequestPacket",
    "name": "_requestPacket",
    "type": "tuple"
  },
  {
    "components": [
      {
        "internalType": "string",
        "name": "clientId",
        "type": "string"
      },
      {
        "internalType": "uint64",
        "name": "requestId",
        "type": "uint64"
      },
      {
        "internalType": "uint64",
        "name": "ansCount",
        "type": "uint64"
      },
      {
        "internalType": "uint64",
        "name": "requestTime",
        "type": "uint64"
      },
      {
        "internalType": "uint64",
        "name": "resolveTime",
        "type": "uint64"
      },
      {
        "internalType": "uint8",
        "name": "resolveStatus",
        "type": "uint8"
      },
      {
        "internalType": "string",
        "name": "result",
        "type": "string"
      }
    ],
    "internalType": "struct IBridge.ResponsePacket",
    "name": "_responsePacket",
    "type": "tuple"
  },
  {
    "internalType": "uint256",
    "name": "_version",
    "type": "uint256"
  },
  {
    "components": [
      {
        "internalType": "bool",
        "name": "isDataOnRight",
        "type": "bool"
      },
      {
        "internalType": "uint8",
        "name": "subtreeHeight",
        "type": "uint8"
      },
      {
        "internalType": "uint256",
        "name": "subtreeSize",
        "type": "uint256"
      },
      {
        "internalType": "uint256",
        "name": "subtreeVersion",
        "type": "uint256"
      },
      {
        "internalType": "bytes32",
        "name": "siblingHash",
        "type": "bytes32"
      }
    ],
    "internalType": "struct IAVLMerklePath.Data[]",
    "name": "_merklePaths",
    "type": "tuple[]"
  }
]
`)
