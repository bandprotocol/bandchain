package main

var rawABI = []byte(`
[
 {
  "inputs": [
   {
    "components": [
     {
      "internalType": "address",
      "name": "addr",
      "type": "address"
     },
     {
      "internalType": "uint256",
      "name": "power",
      "type": "uint256"
     }
    ],
    "internalType": "struct Bridge.ValidatorWithPower[]",
    "name": "_validators",
    "type": "tuple[]"
   }
  ],
  "payable": false,
  "stateMutability": "nonpayable",
  "type": "constructor"
 },
 {
  "anonymous": false,
  "inputs": [
   {
    "indexed": true,
    "internalType": "address",
    "name": "previousOwner",
    "type": "address"
   },
   {
    "indexed": true,
    "internalType": "address",
    "name": "newOwner",
    "type": "address"
   }
  ],
  "name": "OwnershipTransferred",
  "type": "event"
 },
 {
  "constant": false,
  "inputs": [
   {
    "internalType": "bytes",
    "name": "_data",
    "type": "bytes"
   }
  ],
  "name": "relayAndVerify",
  "outputs": [
   {
    "components": [
     {
      "internalType": "uint64",
      "name": "oracleScriptId",
      "type": "uint64"
     },
     {
      "internalType": "uint64",
      "name": "requestTime",
      "type": "uint64"
     },
     {
      "internalType": "uint64",
      "name": "aggregationTime",
      "type": "uint64"
     },
     {
      "internalType": "uint64",
      "name": "requestedValidatorsCount",
      "type": "uint64"
     },
     {
      "internalType": "uint64",
      "name": "minCount",
      "type": "uint64"
     },
     {
      "internalType": "uint64",
      "name": "reportedValidatorsCount",
      "type": "uint64"
     },
     {
      "internalType": "bytes",
      "name": "params",
      "type": "bytes"
     },
     {
      "internalType": "bytes",
      "name": "data",
      "type": "bytes"
     }
    ],
    "internalType": "struct IBridge.VerifyOracleDataResult",
    "name": "result",
    "type": "tuple"
   }
  ],
  "payable": false,
  "stateMutability": "nonpayable",
  "type": "function"
 },
 {
  "constant": false,
  "inputs": [
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
  ],
  "name": "relayOracleState",
  "outputs": [],
  "payable": false,
  "stateMutability": "nonpayable",
  "type": "function"
 },
 {
  "constant": false,
  "inputs": [],
  "name": "renounceOwnership",
  "outputs": [],
  "payable": false,
  "stateMutability": "nonpayable",
  "type": "function"
 },
 {
  "constant": false,
  "inputs": [
   {
    "internalType": "address",
    "name": "newOwner",
    "type": "address"
   }
  ],
  "name": "transferOwnership",
  "outputs": [],
  "payable": false,
  "stateMutability": "nonpayable",
  "type": "function"
 },
 {
  "constant": false,
  "inputs": [
   {
    "components": [
     {
      "internalType": "address",
      "name": "addr",
      "type": "address"
     },
     {
      "internalType": "uint256",
      "name": "power",
      "type": "uint256"
     }
    ],
    "internalType": "struct Bridge.ValidatorWithPower[]",
    "name": "_validators",
    "type": "tuple[]"
   }
  ],
  "name": "updateValidatorPowers",
  "outputs": [],
  "payable": false,
  "stateMutability": "nonpayable",
  "type": "function"
 },
 {
  "constant": true,
  "inputs": [
   {
    "internalType": "bytes",
    "name": "_encodedData",
    "type": "bytes"
   }
  ],
  "name": "decodeResult",
  "outputs": [
   {
    "components": [
     {
      "internalType": "uint64",
      "name": "oracleScriptId",
      "type": "uint64"
     },
     {
      "internalType": "uint64",
      "name": "requestTime",
      "type": "uint64"
     },
     {
      "internalType": "uint64",
      "name": "aggregationTime",
      "type": "uint64"
     },
     {
      "internalType": "uint64",
      "name": "requestedValidatorsCount",
      "type": "uint64"
     },
     {
      "internalType": "uint64",
      "name": "minCount",
      "type": "uint64"
     },
     {
      "internalType": "uint64",
      "name": "reportedValidatorsCount",
      "type": "uint64"
     },
     {
      "internalType": "bytes",
      "name": "params",
      "type": "bytes"
     },
     {
      "internalType": "bytes",
      "name": "data",
      "type": "bytes"
     }
    ],
    "internalType": "struct IBridge.VerifyOracleDataResult",
    "name": "",
    "type": "tuple"
   }
  ],
  "payable": false,
  "stateMutability": "pure",
  "type": "function"
 },
 {
  "constant": true,
  "inputs": [],
  "name": "isOwner",
  "outputs": [
   {
    "internalType": "bool",
    "name": "",
    "type": "bool"
   }
  ],
  "payable": false,
  "stateMutability": "view",
  "type": "function"
 },
 {
  "constant": true,
  "inputs": [
   {
    "internalType": "uint256",
    "name": "",
    "type": "uint256"
   }
  ],
  "name": "oracleStates",
  "outputs": [
   {
    "internalType": "bytes32",
    "name": "",
    "type": "bytes32"
   }
  ],
  "payable": false,
  "stateMutability": "view",
  "type": "function"
 },
 {
  "constant": true,
  "inputs": [],
  "name": "owner",
  "outputs": [
   {
    "internalType": "address",
    "name": "",
    "type": "address"
   }
  ],
  "payable": false,
  "stateMutability": "view",
  "type": "function"
 },
 {
  "constant": true,
  "inputs": [],
  "name": "totalValidatorPower",
  "outputs": [
   {
    "internalType": "uint256",
    "name": "",
    "type": "uint256"
   }
  ],
  "payable": false,
  "stateMutability": "view",
  "type": "function"
 },
 {
  "constant": true,
  "inputs": [
   {
    "internalType": "address",
    "name": "",
    "type": "address"
   }
  ],
  "name": "validatorPowers",
  "outputs": [
   {
    "internalType": "uint256",
    "name": "",
    "type": "uint256"
   }
  ],
  "payable": false,
  "stateMutability": "view",
  "type": "function"
 },
 {
  "constant": true,
  "inputs": [
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
    "name": "_params",
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
  ],
  "name": "verifyOracleData",
  "outputs": [
   {
    "components": [
     {
      "internalType": "uint64",
      "name": "oracleScriptId",
      "type": "uint64"
     },
     {
      "internalType": "uint64",
      "name": "requestTime",
      "type": "uint64"
     },
     {
      "internalType": "uint64",
      "name": "aggregationTime",
      "type": "uint64"
     },
     {
      "internalType": "uint64",
      "name": "requestedValidatorsCount",
      "type": "uint64"
     },
     {
      "internalType": "uint64",
      "name": "minCount",
      "type": "uint64"
     },
     {
      "internalType": "uint64",
      "name": "reportedValidatorsCount",
      "type": "uint64"
     },
     {
      "internalType": "bytes",
      "name": "params",
      "type": "bytes"
     },
     {
      "internalType": "bytes",
      "name": "data",
      "type": "bytes"
     }
    ],
    "internalType": "struct IBridge.VerifyOracleDataResult",
    "name": "",
    "type": "tuple"
   }
  ],
  "payable": false,
  "stateMutability": "view",
  "type": "function"
 }
]
`)
