type t = {
  txHash: Hash.t,
  blockHeight: ID.Block.t,
  success: bool,
  gasFee: list(TxHook.Coin.t),
  gasLimit: float,
  gasUsed: float,
  sender: Address.t,
  timestamp: MomentRe.Moment.t,
};

// subscription {
//   transactions {
//     tx_hash
//     block_height
//     success
//     gas_fee
//   }
// }

module SingleConfig = [%graphql
  {|
  subscription Transactions($tx_hash:String!) {
    transactions_by_pk(tx_hash: $tx_hash) @bsRecord {
        txHash @bsDecoder(fn: "GraphQLParser.hash")
        blockHeight @bsDecoder(fn: "ID.Block.fromJson")
        success @bsDecoder(fn: "GraphQLParser.bool")
        gasFee @bsDecoder(fn: "GraphQLParser.coins")
        gasLimit @bsDecoder(fn: "GraphQLParser.float")
        gasUsed @bsDecoder(fn: "GraphQLParser.float")
        sender gasLimit @bsDecoder(fn: "GraphQLParser.fromBech32")
        timestamp : last_updated @bsDecoder(fn: "GraphQLParser.time")
    }
  },
|}
];
