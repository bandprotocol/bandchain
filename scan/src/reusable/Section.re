module Styles = {
  open Css;

  let base = (~pt, ~pb, ()) => style([paddingTop(`px(pt)), paddingBottom(`px(pb))]);

  let bgColor = color => style([backgroundColor(color)]);
};

[@react.component]
let make = (~children, ~pt=0, ~pb=0, ~bg=Colors.white) => {
  let css = Css.merge([Styles.bgColor(bg), Styles.base(~pt, ~pb, ())]);

  <section className=css> children </section>;
};
