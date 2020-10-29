type t = {
  packetType: string,
  isIncoming: bool,
};

// module MultiConfig = [%graphql
//   {|
//   subscription Packets($limit: Int!, $offset: Int!) {
//     packets(offset: $offset, limit: $limit, order_by: {block_height: desc}) @bsRecord {
//       packetType : type
//       isIncoming: is_incoming
//     }
//   }
// |}
// ];

// let getList = (~page, ~pageSize, ()) => {
// let offset = (page - 1) * pageSize;
// let (result, _) =
//   ApolloHooks.useSubscription(
//     MultiConfig.definition,
//     ~variables=MultiConfig.makeVariables(~limit=pageSize, ~offset, ()),
//   );
// result |> Sub.map(_, x => x##packets);
// Sub.resolve([||]);
// };
