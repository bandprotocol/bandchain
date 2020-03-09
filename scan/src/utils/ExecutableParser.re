let checker = (str: string) => {
  let reg = ".*=[$][0-9]+" |> Js.Re.fromString;
  let t =
    reg
    |> Js.Re.exec(_, str)
    |> Belt_Option.mapWithDefault(_, [||], Js.Re.captures)
    |> Belt_Array.length;
  t > 0;
};

let trim = (str: string) => {
  str |> String.split_on_char('=') |> Belt_List.get(_, 0) |> Belt_Option.getExn;
};

let getVariables = (str: string) => {
  String.split_on_char('\n', str) |> List.filter(checker) |> List.map(trim);
};

let parseExecutableScript = (buff: JsBuffer.t) => {
  buff |> JsBuffer._toString(_, "UTF-8") |> getVariables;
};
