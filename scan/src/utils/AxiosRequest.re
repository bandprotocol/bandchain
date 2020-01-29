[@bs.deriving abstract]
type t = {
  codeHash: string,
  params: Js.Dict.t(string),
};

let convert: t => Js.t('a) = [%bs.raw {|
function(data) {
  return data;
}
  |}];

let execute = (data: t) => {
  let%Promise response =
    Axios.postData("http://d3n-debug.bandprotocol.com:5000/request", convert(data));
  Promise.ret(response);
};
