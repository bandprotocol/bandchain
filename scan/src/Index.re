[@bs.val] external document: Js.t({..}) = "document";

let style = document##createElement("style");
document##head##appendChild(style);
style##innerHTML #= AppStyle.style;

AxiosHooks.setRpcUrl(Env.rpc);

TimeAgos.setMomentRelativeTimeThreshold();

let setupSentry: unit => unit = [%bs.raw
  {|
function() {
  const Sentry = require("@sentry/browser");
  Sentry.init({dsn: "https://6f05376ceab44557943d1864072a37ae@o270592.ingest.sentry.io/5260152"});
}
  |}
];
setupSentry();

ReactDOMRe.render(
  <ApolloClient>
    <GlobalContext>
      <TimeContext>
        <ModalContext> <AccountContext> <App /> </AccountContext> </ModalContext>
      </TimeContext>
    </GlobalContext>
  </ApolloClient>,
  document##getElementById("root"),
);
