module Styles = {
  open Css;
  let container = style([fontSize(`px(14)), marginLeft(`em(0.5)), cursor(`pointer)]);
};

[@react.component]
let make = (~text, ~to_) => {
  <div className=Styles.container onClick={_ => Js.Console.log(to_)}>
    {text |> React.string}
  </div>;
};
