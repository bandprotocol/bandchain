[@bs.val] external document: Js.t({..}) = "document";

let style = document##createElement("style");
document##head##appendChild(style);
style##innerHTML #= AppStyle.style;

AxiosHooks.setRpcUrl(Env.rpc);

TimeAgos.setMomentRelativeTimeThreshold();

ReactDOMRe.render(
  <ApolloClient>
    <GlobalContext> <AccountContext> <App /> </AccountContext> </GlobalContext>
  </ApolloClient>,
  document##getElementById("root"),
);
