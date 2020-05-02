// shell version
let shellChecker = str => {
  let reg = ".*=[$][1-9][0-9]*" |> Js.Re.fromString;
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
  let s1 = getElementInList(tmp, 1) |> String.split_on_char('$') |> getElementInList(_, 1);

  switch (s1 |> int_of_string) {
  // the number is too big then return -1
  | exception (Failure(_)) => (s0, (-1))
  | num => (s0, num)
  };
};

let comparePair = ((_, num1), (_, num2)) => {
  compare(num1, num2);
};

let checkValid = pairs => {
  let nums = pairs |> List.map(((_, x)) => x) |> List.sort_uniq(compare);
  let len = nums |> Belt_List.length;

  len > 0 && len == getElementInList(nums, len - 1) && 1 == getElementInList(nums, 0);
};

// python version
let pythonMatch = str => {
  let reg = "def main\(\s*([^)]+?)\s*\)" |> Js.Re.fromString;
  let rawResult =
    reg |> Js.Re.exec_(_, str) |> Belt_Option.mapWithDefault(_, [||], Js.Re.captures);

  switch (rawResult->Belt.Array.get(1)) {
  | Some(result) =>
    switch (result->Js.Nullable.toOption) {
    | Some(result') =>
      Some(result' |> String.split_on_char(',') |> Belt.List.map(_, String.trim))
    | None => None
    }
  | None => None
  };
};

let getVariables = str => {
  let splitedStr = String.split_on_char('\n', str);
  let%Opt program = splitedStr->Belt.List.get(0);

  switch (program) {
  | "#!/bin/bash"
  | "#!/bin/sh" =>
    let pairs =
      splitedStr |> List.filter(shellChecker) |> List.map(splitToPair) |> List.sort(comparePair);

    pairs |> checkValid ? Some(pairs |> List.map(((x, _)) => x)) : None;
  | "#!/usr/bin/env python3" => str |> pythonMatch
  | _ => None
  };
};

let parseExecutableScript = (buff: JsBuffer.t) => {
  buff |> JsBuffer._toString(_, "UTF-8") |> getVariables;
};
