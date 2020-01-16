type size =
  | Xs
  | Sm
  | Md
  | Lg
  | Xl
  | Xxl;

type weight =
  | Regular
  | Semibold
  | Bold;

type align =
  | Center
  | Right;

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
      | Xl => style([fontSize(px(16))])
      | Xxl => style([fontSize(px(24))]),
    );

  let fontWeight =
    mapWithDefault(
      _,
      style([]),
      fun
      | Regular => style([fontWeight(`normal)])
      | Semibold => style([fontWeight(`medium)])
      | Bold => style([fontWeight(`bold)]),
    );

  let noWrap = style([whiteSpace(`nowrap)]);
  let block = style([display(`block)]);
  let ellipsis = style([overflow(`hidden), textOverflow(`ellipsis), whiteSpace(`nowrap)]);
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
      fontFamily(
        "Fira Code, cousine, sfmono-regular,Consolas,Menlo,liberation mono,ubuntu mono,Courier,monospace",
      ),
    ]);
};

[@react.component]
let make =
    (
      ~size=?,
      ~weight=?,
      ~align=?,
      ~nowrap=false,
      ~color=?,
      ~block=false,
      ~code=false,
      ~ellipsis=false,
      ~value,
    ) => {
  <span
    className={Css.merge([
      Styles.fontSize(size),
      Styles.fontWeight(weight),
      Styles.textAlign(align),
      nowrap ? Styles.noWrap : "",
      block ? Styles.block : "",
      code ? Styles.code : "",
      color->Belt.Option.mapWithDefault("", c => Css.style([Css.color(c)])),
      ellipsis ? Styles.ellipsis : "",
    ])}>
    {React.string(value)}
  </span>;
};
