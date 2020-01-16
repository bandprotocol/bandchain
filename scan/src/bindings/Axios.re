let setRpcUrl: string => unit = [%bs.raw
  {|
function(rpcUrl) {
  const AxiosHooks = require("axios-hooks");
  const Axios = require("axios");
  AxiosHooks.configure({
    axios: Axios.create({
      baseURL: rpcUrl,
    }),
  });
}
  |}
];

[@bs.deriving abstract]
type t = {
  data: Js.undefined(Js.Json.t),
  loading: bool,
};

[@bs.val] [@bs.module "axios-hooks"]
external _context: string => (t, (unit, unit) => unit) = "default";

let use = (url, ~pollInterval=600000, ()) => {
  let (rawdata, refetch) = _context(url);
  React.useEffect2(
    () => {
      let intervalId = Js.Global.setInterval(() => refetch((), ()), pollInterval);
      Some(() => Js.Global.clearInterval(intervalId));
    },
    (url, pollInterval),
  );
  Js.undefinedToOption(rawdata->dataGet);
};
