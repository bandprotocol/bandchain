type btn_style_t =
  | Primary
  | Outline;

module Styles = {
  open Css;

  let btn = (~variant=Primary, ~fsize=12, ~px=25, ~py=13, ~pxSm=px, ~pySm=py, ~disable, ()) => {
    let base =
      style([
        display(`block),
        padding2(~v=`px(py), ~h=`px(px)),
        transition(~duration=200, "all"),
        borderRadius(`px(4)),
        fontSize(`px(fsize)),
        cursor(disable ? `default : `pointer),
        outlineStyle(`none),
        borderStyle(`none),
        margin(`zero),
        Media.mobile([padding2(~v=`px(pySm), ~h=`px(pxSm))]),
      ]);

    let custom =
      switch (variant) {
      | Primary =>
        style([
          backgroundColor(Colors.bandBlue),
          color(Colors.white),
          hover([backgroundColor(Colors.buttonBaseHover)]),
          active([backgroundColor(Colors.buttonBaseActive)]),
          disabled([backgroundColor(Colors.buttonDisabled), color(Colors.white)]),
        ])
      | Outline =>
        style([
          backgroundColor(Colors.white),
          color(Colors.bandBlue),
          border(`px(1), `solid, Colors.bandBlue),
          hover([backgroundColor(Colors.buttonOutlineHover)]),
          active([backgroundColor(Colors.buttonOutlineActive)]),
          disabled([borderColor(Colors.buttonDisabled), color(Colors.buttonDisabled)]),
        ])
      };
    merge([base, custom]);
  };
};

[@react.component]
let make =
    (
      ~variant=Primary,
      ~children,
      ~py=5,
      ~px=10,
      ~fsize=12,
      ~pySm=py,
      ~pxSm=px,
      ~onClick,
      ~style="",
      ~disabled=false,
    ) => {
  <button
    className={Css.merge([
      Styles.btn(~variant, ~px, ~py, ~pxSm, ~pySm, ~fsize, ~disable=disabled, ()),
      CssHelper.flexBox(~align=`center, ~justify=`center, ()),
      style,
    ])}
    onClick
    disabled>
    children
  </button>;
};
