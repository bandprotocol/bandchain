[@bs.deriving abstract]
type t = {
  codeHash: string,
  params: Js.Dict.t(string),
};

let execute: t => unit = [%bs.raw
  {|
function(data) {
  console.log(data);
  const Axios = require("axios");
  Axios.post("http://d3n-debug.bandprotocol.com:5000/request", data)
    .then(response => console.log(response))
    .catch(err => console.log(err));
}
  |}
];

let exexx = (data: t) => {
  // let%Promise response = Axios.post("http://d3n-debug.bandprotocol.com:5000/request", data);
  // if (response##status == 200) {
  //   dispatchTx(AddFaucet(response##data##result));
  //   dispatchAnalytics(TrackFaucet);
  // } else if (response##status == 208) {
  //   "Sorry, you got rate limited. Please try again in the next 24 hours." |> Window.alert;
  // };
  // Promise.ret $$ ();
};
