[@bs.val] external document: Js.t({..}) = "document";

let style = document##createElement("style");
document##head##appendChild(style);
style##innerHTML #= AppStyle.style;

Axios.setRpcUrl("http://localhost:8010/");
TimeAgos.setMomentLocale();

ReactDOMRe.render(<GlobalContext> <App /> </GlobalContext>, document##getElementById("root"));
