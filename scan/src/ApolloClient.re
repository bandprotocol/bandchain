let wsLink = ApolloLinks.webSocketLink(~uri=Env.graphql, ~reconnect=true, ());

let client =
  ReasonApollo.createApolloClient(
    ~link=wsLink,
    ~cache=ApolloInMemoryCache.createInMemoryCache(),
    (),
  );

[@react.component]
let make = (~children) => {
  <ApolloHooks.Provider client> children </ApolloHooks.Provider>;
};
