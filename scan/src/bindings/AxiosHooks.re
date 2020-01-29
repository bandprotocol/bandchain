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
type context_config_t = {useCache: bool};

[@bs.deriving abstract]
type t = {
  data: Js.undefined(Js.Json.t),
  loading: bool,
};

[@bs.val] [@bs.module "axios-hooks"]
external _context: (string, context_config_t) => (t, (unit, unit) => unit) = "default";

let use = (url, ~pollInterval=600000, ()) => {
  let (rawdata, refetch) = _context(url, context_config_t(~useCache=false));
  React.useEffect2(
    () => {
      let intervalId = Js.Global.setInterval(() => refetch((), ()), pollInterval);
      Some(() => Js.Global.clearInterval(intervalId));
    },
    (url, pollInterval),
  );
  Js.undefinedToOption(rawdata->dataGet);
};
