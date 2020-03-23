// in memory cache for caching GraphQL data
let graphQL = "d3n-debug.bandprotocol.com:5433/v1/graphql";
let wsLink = ApolloLinks.webSocketLink(~uri="wss://" ++ graphQL, ~reconnect=true, ());
let httpLink = ApolloLinks.createHttpLink(~uri="http://" ++ graphQL, ());
let link =
  ApolloLinks.split(
    e => {
      let definition = ApolloUtilities.getMainDefinition(e##query);
      definition##kind === "OperationDefinition" && definition##operation === "subscription";
    },
    wsLink,
    httpLink,
  );

/* Create an InMemoryCache */
let inMemoryCache = ApolloInMemoryCache.createInMemoryCache();
// /* Create an HTTP Link */
let client = ReasonApollo.createApolloClient(~link, ~cache=inMemoryCache, ());

[@react.component]
let make = (~children) => {
  <ApolloHooks.Provider client> children </ApolloHooks.Provider>;
};
