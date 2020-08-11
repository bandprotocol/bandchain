open Css;

let flexBox = (~align=`center, ~justify=`flexStart, ~wrap=`wrap, ~direction=`row, ()) =>
  style([
    display(`flex),
    alignItems(align),
    justifyContent(justify),
    flexDirection(direction),
    flexWrap(wrap),
  ]);
