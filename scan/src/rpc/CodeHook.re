let getCode = url => {
  let json = AxiosHooks.use(url);
  json |> Belt.Option.flatMap(_, x => x |> Js.Json.decodeString);
};
