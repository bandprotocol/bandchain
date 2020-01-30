[@bs.deriving abstract]
type t = {
  codeHash: string,
  params: Js.Dict.t(string),
};

let convert: t => Js.t('a) = [%bs.raw
  {|
function(data) {
  let params = data.params;
  let ret = {};
  for (let key of Object.keys(params)) {
    if (isNaN(params[key])) ret[key] = params[key]
    else ret[key] = parseInt(params[key], 10)
  }
  return {...data, params: ret};
}
  |}
];

let execute = (data: t) => {
  let%Promise response =
    Axios.postData("https://d3n.bandprotocol.com/bandsv/request", convert(data));
  Promise.ret(response);
};
