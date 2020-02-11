[@bs.val] external document: Js.t({..}) = "document";

let style = document##createElement("style");
document##head##appendChild(style);
style##innerHTML #= AppStyle.style;

AxiosHooks.setRpcUrl("http://rpc.alpha.bandchain.org");

TimeAgos.setMomentRelativeTimeThreshold();

ReactDOMRe.render(<GlobalContext> <App /> </GlobalContext>, document##getElementById("root"));
