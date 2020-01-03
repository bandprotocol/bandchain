[@bs.val] external document: Js.t({..}) = "document";

let style = document##createElement("style");
document##head##appendChild(style);
style##innerHTML #= AppStyle.style;

Axios.setRpcUrl("http://d3n.bandprotocol.com:1317/");
ReactDOMRe.render(<App />, document##getElementById("root"));
