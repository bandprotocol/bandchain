module Styles = {
  open Css;

  let wrapper = isContent =>
    style([
      display(`flex),
      justifyContent(`center),
      width(`percent(100.)),
      minHeight(isContent ? `calc((`sub, `vh(100.), `px(380))) : `unset),
    ]);

  let container = style([display(`flex), maxWidth(`px(1180)), width(`percent(100.))]);
};

[@react.component]
let make = (~isContent=false, ~children) => {
  <div className={Styles.wrapper(isContent)}>
    <div className=Styles.container> children </div>
  </div>;
};
