// in memory cache for caching GraphQL data
let graphQL = "d3n-debug.bandprotocol.com:5433/v1/graphql";
let wsLink = ApolloLinks.webSocketLink(~uri="ws://" ++ graphQL, ~reconnect=true, ());

/* Create an InMemoryCache */
let inMemoryCache = ApolloInMemoryCache.createInMemoryCache();
// /* Create an HTTP Link */
let client = ReasonApollo.createApolloClient(~link=wsLink, ~cache=inMemoryCache, ());

[@react.component]
let make = (~children) => {
  <ApolloHooks.Provider client> children </ApolloHooks.Provider>;
};
