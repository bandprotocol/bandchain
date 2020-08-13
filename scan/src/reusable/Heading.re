type size =
  | H1
  | H2
  | H3
  | H4
  | H5;

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
    | H1 => style([fontSize(`px(24)), Media.smallMobile([fontSize(`px(20))])])
    | H2 => style([fontSize(`px(20)), Media.smallMobile([fontSize(`px(18))])])
    | H3 => style([fontSize(`px(18)), Media.smallMobile([fontSize(`px(16))])])
    | H4 => style([fontSize(`px(14)), Media.smallMobile([fontSize(`px(12))])])
    | H5 => style([fontSize(`px(12)), Media.smallMobile([fontSize(`px(11))])]);

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
  let textColor = color_ => {
    style([color(color_)]);
  };
  let mb = size => {
    style([marginBottom(`px(size))]);
  };
  let mbSm = size => {
    style([marginBottom(`px(size))]);
  };
};

[@react.component]
let make =
    (
      ~value,
      ~align=Left,
      ~weight=Semibold,
      ~size=H1,
      ~marginBottom=0,
      ~marginBottomSm=marginBottom,
      ~style="",
      ~color=Colors.gray7,
    ) => {
  let children_ = React.string(value);
  let style_ = size =>
    Css.merge(
      Styles.[
        fontSize(size),
        fontWeight(weight),
        textColor(color),
        textAlign(align),
        lineHeight,
        mb(marginBottom),
        mbSm(marginBottomSm),
        style,
      ],
    );

  switch (size) {
  | H1 => <h1 className={style_(size)}> children_ </h1>
  | H2 => <h2 className={style_(size)}> children_ </h2>
  | H3 => <h3 className={style_(size)}> children_ </h3>
  | H4 => <h4 className={style_(size)}> children_ </h4>
  | H5 => <h5 className={style_(size)}> children_ </h5>
  };
};
