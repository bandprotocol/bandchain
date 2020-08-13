let ellipsis = (~text, ~limit, ()) => {
  Js.String2.length(text) > limit ? Js.String.slice(~from=2, ~to_=limit, text) ++ "..." : text;
};
