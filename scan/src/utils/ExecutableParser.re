let checker = (str: string) => {
  let reg = ".*=[$][0-9]+" |> Js.Re.fromString;
  let t =
    reg
    |> Js.Re.exec_(_, str)
    |> Belt_Option.mapWithDefault(_, [||], Js.Re.captures)
    |> Belt_Array.length;
  t > 0;
};

let getElementInList = (l, idx) => {
  Belt_List.get(l, idx) |> Belt_Option.getExn;
};

let splitToPair = s => {
  let tmp = s |> String.split_on_char('=');
  let s0 = getElementInList(tmp, 0);
  let s1 =
    getElementInList(tmp, 1)
    |> String.split_on_char('$')
    |> getElementInList(_, 1)
    |> int_of_string;
  (s0, s1);
};

let comparePair = ((_, num1), (_, num2)) => {
  compare(num1, num2);
};

let checkValid = pairs => {
  let nums = pairs |> List.map(((_, x)) => x) |> List.sort_uniq(compare);
  let len = nums |> Belt_List.length;

  len > 0 && len == getElementInList(nums, len - 1) && 1 == getElementInList(nums, 0);
};

let getVariables = str => {
  let pairs =
    String.split_on_char('\n', str)
    |> List.filter(checker)
    |> List.map(splitToPair)
    |> List.sort(comparePair);

  pairs |> checkValid ? Some(pairs |> List.map(((x, _)) => x)) : None;
};

let parseExecutableScript = (buff: JsBuffer.t) => {
  buff |> JsBuffer._toString(_, "UTF-8") |> getVariables;
};
