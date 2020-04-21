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
  let%Promise response = Axios.postData(Env.faucet, convert(data));
  Promise.ret(response);
};
