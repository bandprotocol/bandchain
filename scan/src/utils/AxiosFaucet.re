type t = {
  address: string,
  amount: int,
};

let convert: t => Js.t('a) = [%bs.raw {|
function(data) {
  return {...data};
}
  |}];

let request = (data: t) => {
  let%Promise response =
    Axios.postData("https://d3n.bandprotocol.com/faucet/request", convert(data));
  Promise.ret(response);
};
