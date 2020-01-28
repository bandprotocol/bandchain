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
