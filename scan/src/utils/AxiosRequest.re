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
    Axios.postData("https://d3n.bandprotocol.com/bandsv/request", convert(data));
  Promise.ret(response);
};
