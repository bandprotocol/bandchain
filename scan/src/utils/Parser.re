let checker = (str: string) => {
  let len = String.length(str);
  len >= 3 && str.[len - 3] == '=' && str.[len - 2] == '$' && str.[len - 1] == '1';
};

let trim = (str: string) => {
  let len = String.length(str);
  String.sub(str, 0, len - 3);
};

let getVariable = (str: string) => {
  String.split_on_char('\n', str) |> List.filter(checker) |> List.map(trim);
};
// external btoa : string -> string = "" [@@bs.val]

let parseExecutableScript = (buff: JsBuffer.t) => {
  buff |> JsBuffer._toString(_, "UTF-8") |> getVariable;
};
