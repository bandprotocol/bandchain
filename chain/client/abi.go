package rpc

var relayFormat = []byte(`
[
  {
    "internalType": "uint256",
    "name": "_blockHeight",
    "type": "uint256"
  },
  {
    "internalType": "bytes32",
    "name": "_oracleIAVLStateHash",
    "type": "bytes32"
  },
  {
    "internalType": "bytes32",
    "name": "_otherStoresMerkleHash",
    "type": "bytes32"
  },
  {
    "internalType": "bytes32",
    "name": "_supplyStoresMerkleHash",
    "type": "bytes32"
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
        "name": "txCountAndLastBlockInfoHash",
        "type": "bytes32"
      },
      {
        "internalType": "bytes32",
        "name": "consensusDataHash",
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
    "internalType": "bytes",
    "name": "_data",
    "type": "bytes"
  },
  {
    "internalType": "uint64",
    "name": "_requestId",
    "type": "uint64"
  },
  {
    "internalType": "uint64",
    "name": "_oracleScriptId",
    "type": "uint64"
  },
  {
    "internalType": "bytes",
    "name": "_calldata",
    "type": "bytes"
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
