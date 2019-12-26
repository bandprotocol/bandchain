[@react.component]
let make = () => {
  // let (data, loading) =
  //   Axios.use("http://localhost:8010/blocks/latest", ~pollInterval=1000, ());

  let data = BlockHook.latest(~pollInterval=1000, ());
  Js.Console.log(data);

  <div> {React.string(Js_json.stringifyAny(data)->Belt_Option.getWithDefault(""))} </div>;
};
