module Styles = {
  open Css;

  let fontSize = size => style([fontSize(`px(size))]);
  let fontColor = color_ => style([color(color_)]);
};

[@react.component]
let make = (~name, ~color=Colors.gray1, ~size=12) => {
  <i className={Css.merge([name, Styles.fontColor(color), Styles.fontSize(size)])} />;
};
