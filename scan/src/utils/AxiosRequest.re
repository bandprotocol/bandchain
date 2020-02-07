[@bs.deriving abstract]
type t = {
  codeHash: string,
  params: Js.Dict.t(string),
};

/* TODO: FIX THIS MESS */
let convert: t => Js.t('a) = [%bs.raw
  {|
function(data) {
  let params = data.params;
  let ret = {};
  for (let key of Object.keys(params)) {
    if (isNaN(params[key])) ret[key] = params[key]
    else ret[key] = parseInt(params[key], 10)
  }
  return {...data, params: ret, type: "SYNCHRONOUS"};
}
  |}
];

let execute = (data: t) => {
  let%Promise response =
    Axios.postData("http://rpc.alpha.d3n.xyz/bandsv/request", convert(data));
  Promise.ret(response);
};
