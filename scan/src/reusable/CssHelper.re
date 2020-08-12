open Css;

let flexBox = (~align=`center, ~justify=`flexStart, ~wrap=`wrap, ~direction=`row, ()) =>
  style([
    display(`flex),
    alignItems(align),
    justifyContent(justify),
    flexDirection(direction),
    flexWrap(wrap),
  ]);

let btn = (~fsize=12, ()) => {
  style([
    display(`block),
    fontSize(`px(fsize)),
    backgroundColor(Colors.bandBlue),
    color(Colors.white),
    padding2(~v=`px(13), ~h=`px(25)),
    transition(~duration=200, "all"),
    borderRadius(`px(4)),
    hover([backgroundColor(Colors.bandBlue)]),
  ]);
};
