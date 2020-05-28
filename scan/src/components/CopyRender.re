module Styles = {
  open Css;

  let copy = w => style([width(`px(w)), cursor(`pointer), position(`relative), zIndex(2)]);
  let tick = w => style([width(`px(w))]);
};

[@react.component]
let make = (~width, ~message) => {
  let (copied, setCopy) = React.useState(_ => false);

  copied
    ? <img src=Images.tickIcon className={Styles.tick(width)} />
    : <img
        src=Images.copy
        className={Styles.copy(width)}
        onClick={_ => {
          Copy.copy(message);
          setCopy(_ => true);
          let _ = Js.Global.setTimeout(() => setCopy(_ => false), 700);
          ();
        }}
      />;
};
