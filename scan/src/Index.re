[@bs.val] external document: Js.t({..}) = "document";

let style = document##createElement("style");
document##head##appendChild(style);
style##innerHTML #= AppStyle.style;

Axios.setRpcUrl("https://cors-anywhere.herokuapp.com/http://d3n.bandprotocol.com:1317");
TimeAgos.setMomentRelativeTimeThreshold();

ReactDOMRe.render(<GlobalContext> <App /> </GlobalContext>, document##getElementById("root"));
