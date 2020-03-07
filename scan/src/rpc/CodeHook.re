module Code = {
  type file_t = {
    name: string,
    content: string,
  };

  type t = list(file_t);

  let decodeCode = json =>
    JsonUtils.Decode.{
      name: json |> field("name", string),
      content: json |> field("content", string),
    };

  let decodeCodes = json => JsonUtils.Decode.(json |> list(decodeCode));
};

let getCode = (codeHash: Hash.t) => {
  let codeHashHex = codeHash |> Hash.toHex(~upper=true);
  let json =
    AxiosHooks.use(
      {j|https://s3.ap-southeast-1.amazonaws.com/code.d3n.bandprotocol.com/$codeHashHex|j},
    );
  json |> Belt.Option.map(_, Code.decodeCodes);
};
