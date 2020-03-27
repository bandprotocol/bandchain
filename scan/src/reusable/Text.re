type size =
  | Xs
  | Sm
  | Md
  | Lg
  | Xl
  | Xxl
  | Xxxl;

type weight =
  | Thin
  | Regular
  | Medium
  | Semibold
  | Bold;

type align =
  | Center
  | Right;

type spacing =
  | Unset
  | Em(float);

type lineHeight =
  | Px(int)
  | PxFloat(float);

module Styles = {
  open Css;
  open Belt.Option;

  let fontSize =
    mapWithDefault(
      _,
      style([fontSize(px(12))]),
      fun
      | Xs => style([fontSize(px(8)), letterSpacing(`em(0.07))])
      | Sm => style([fontSize(px(10)), letterSpacing(`em(0.05))])
      | Md => style([fontSize(px(12))])
      | Lg => style([fontSize(px(14)), lineHeight(`px(18))])
      | Xl => style([fontSize(px(16)), lineHeight(`px(18))])
      | Xxl => style([fontSize(px(18))])
      | Xxxl => style([fontSize(px(24))]),
    );

  let fontWeight =
    mapWithDefault(
      _,
      style([]),
      fun
      | Thin => style([fontWeight(`num(300))])
      | Regular => style([fontWeight(`num(400))])
      | Medium => style([fontWeight(`num(500))])
      | Semibold => style([fontWeight(`num(600))])
      | Bold => style([fontWeight(`num(700))]),
    );

  let lineHeight =
    mapWithDefault(
      _,
      style([]),
      fun
      | Px(height) => style([lineHeight(`px(height))])
      | PxFloat(height) => style([lineHeight(`pxFloat(height))]),
    );

  let letterSpacing =
    mapWithDefault(
      _,
      style([letterSpacing(`unset)]),
      fun
      | Unset => style([letterSpacing(`unset)])
      | Em(spacing) => style([letterSpacing(`em(spacing))]),
    );

  let noWrap = style([whiteSpace(`nowrap)]);
  let block = style([display(`block)]);
  let ellipsis = style([overflow(`hidden), textOverflow(`ellipsis), whiteSpace(`nowrap)]);
  let underline = style([textDecoration(`underline)]);
  let textAlign =
    mapWithDefault(
      _,
      style([textAlign(`left)]),
      fun
      | Center => style([textAlign(`center)])
      | Right => style([textAlign(`right)]),
    );

  let code =
    style([
      fontFamilies([
        `custom("IBM Plex Mono"),
        `custom("cousine"),
        `custom("sfmono-regular"),
        `custom("Consolas"),
        `custom("Menlo"),
        `custom("liberation mono"),
        `custom("ubuntu mono"),
        `custom("Courier"),
        `monospace,
      ]),
    ]);
};

[@react.component]
let make =
    (
      ~size=?,
      ~weight=?,
      ~align=?,
      ~spacing=?,
      ~height=?,
      ~nowrap=false,
      ~color=?,
      ~block=false,
      ~code=false,
      ~ellipsis=false,
      ~underline=false,
      ~value,
    ) => {
  <span
    className={Css.merge([
      Styles.fontSize(size),
      Styles.fontWeight(weight),
      Styles.textAlign(align),
      Styles.letterSpacing(spacing),
      Styles.lineHeight(height),
      nowrap ? Styles.noWrap : "",
      block ? Styles.block : "",
      code ? Styles.code : "",
      color->Belt.Option.mapWithDefault("", c => Css.style([Css.color(c)])),
      ellipsis ? Styles.ellipsis : "",
      underline ? Styles.underline : "",
    ])}>
    {React.string(value)}
  </span>;
};
