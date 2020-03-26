let normalizeHexString = hexstr =>
  hexstr
  ->Js.Re.exec_("[0-9a-fA-F]+"->Js.Re.fromString, _)
  ->Belt_Option.mapWithDefault([||], result =>
      result->Js.Re.captures->Belt_Array.keepMap(Js.toOption)
    )
  ->Belt_Array.get(0)
  ->Belt_Option.getWithDefault(_, "");
