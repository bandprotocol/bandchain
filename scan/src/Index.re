[@bs.val] external document: Js.t({..}) = "document";

let style = document##createElement("style");
document##head##appendChild(style);
style##innerHTML #= AppStyle.style;

AxiosHooks.setRpcUrl(Env.rpc);

TimeAgos.setMomentRelativeTimeThreshold();

ReactDOMRe.render(
  <ApolloClient>
    <GlobalContext>
      <ModalContext> <AccountContext> <App /> </AccountContext> </ModalContext>
    </GlobalContext>
  </ApolloClient>,
  document##getElementById("root"),
);
