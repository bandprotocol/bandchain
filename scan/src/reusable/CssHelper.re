open Css;

let flexBox = (~align=`center, ~justify=`flexStart, ~wrap=`wrap, ~direction=`row, ()) =>
  style([
    display(`flex),
    alignItems(align),
    justifyContent(justify),
    flexDirection(direction),
    flexWrap(wrap),
  ]);
let flexBoxSm = (~align=`center, ~justify=`flexStart, ~wrap=`wrap, ~direction=`row, ()) =>
  style([
    Media.mobile([
      display(`flex),
      alignItems(align),
      justifyContent(justify),
      flexDirection(direction),
      flexWrap(wrap),
    ]),
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
      outlineStyle(`none),
      borderStyle(`none),
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

let mb = (~size=8, ()) => {
  style([marginBottom(`px(size))]);
};
let mbSm = (~size=8, ()) => {
  style([Media.mobile([marginBottom(`px(size))])]);
};

let px = (~size=0, ()) => {
  style([paddingLeft(`px(size)), paddingRight(`px(size))]);
};

let pxSm = (~size=0, ()) => {
  style([Media.mobile([paddingLeft(`px(size)), paddingRight(`px(size))])]);
};
