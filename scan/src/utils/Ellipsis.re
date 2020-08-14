let format = (~text, ~limit, ()) => {
  Js.String2.length(text) > limit ? Js.String.slice(~from=0, ~to_=limit, text) ++ "..." : text;
};
