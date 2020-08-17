open Css;

let flexBox = (~align=`center, ~justify=`flexStart, ~wrap=`wrap, ~direction=`row, ()) =>
  style([
    display(`flex),
    alignItems(align),
    justifyContent(justify),
    flexDirection(direction),
    flexWrap(wrap),
  ]);

// TODO: abstract later
type btn_style_t =
  | Primary
  | Secondary
  | Outline;

let btn = (~variant=Primary, ~fsize=12, ~px=25, ~py=13, ()) => {
  let base =
    style([
      display(`block),
      padding2(~v=`px(py), ~h=`px(px)),
      transition(~duration=200, "all"),
      borderRadius(`px(4)),
      fontSize(`px(fsize)),
      cursor(`pointer),
    ]);

  let custom =
    switch (variant) {
    | Primary => style([backgroundColor(Colors.bandBlue), color(Colors.white)])
    | Secondary => style([]) // TODO: add later
    | Outline =>
      style([
        backgroundColor(Colors.white),
        color(Colors.bandBlue),
        border(`px(1), `solid, Colors.bandBlue),
      ])
    };
  merge([base, custom]);
};

let mobileSpacing = style([Media.mobile([paddingBottom(`px(20))])]);

let clickable = style([cursor(`pointer)]);

let container = "container";
