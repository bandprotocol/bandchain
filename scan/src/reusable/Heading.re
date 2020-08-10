type size =
  | H1
  | H2
  | H3;
type weight =
  | Thin
  | Regular
  | Medium
  | Semibold
  | Bold;

type align =
  | Center
  | Right
  | Left;

module Styles = {
  open Css;
  let lineHeight = style([lineHeight(`em(1.2))]);
  let fontSize =
    fun
    | H1 => style([fontSize(`px(24)), Media.mobile([fontSize(`px(20))])])
    | H2 => style([fontSize(`px(20)), Media.mobile([fontSize(`px(18))])])
    | H3 => style([fontSize(`px(18)), Media.mobile([fontSize(`px(16))])]);

  let fontWeight =
    fun
    | Thin => style([fontWeight(`num(300))])
    | Regular => style([fontWeight(`num(400))])
    | Medium => style([fontWeight(`num(500))])
    | Semibold => style([fontWeight(`num(600))])
    | Bold => style([fontWeight(`num(700))]);

  let textAlign =
    fun
    | Center => style([textAlign(`center)])
    | Right => style([textAlign(`right)])
    | Left => style([textAlign(`left)]);

  let marginBottom = size => {
    style([marginBottom(`px(size))]);
  };
};

[@react.component]
let make = (~value, ~align=Left, ~weight=Semibold, ~size=H1, ~marginBottom=0) => {
  switch (size) {
  | H1 =>
    <h1
      className={Css.merge([
        Styles.fontSize(size),
        Styles.fontWeight(weight),
        Styles.textAlign(align),
        Styles.lineHeight,
        Styles.marginBottom(marginBottom),
      ])}>
      {React.string(value)}
    </h1>
  | H2 =>
    <h2
      className={Css.merge([
        Styles.fontSize(size),
        Styles.fontWeight(weight),
        Styles.textAlign(align),
        Styles.lineHeight,
        Styles.marginBottom(marginBottom),
      ])}>
      {React.string(value)}
    </h2>
  | H3 =>
    <h3
      className={Css.merge([
        Styles.fontSize(size),
        Styles.fontWeight(weight),
        Styles.textAlign(align),
        Styles.lineHeight,
        Styles.marginBottom(marginBottom),
      ])}>
      {React.string(value)}
    </h3>
  };
};
