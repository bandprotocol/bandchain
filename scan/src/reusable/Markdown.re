module Styles = {
  open Css;

  let container = style([
    selector("a", [wordBreak(`breakAll)]),
  ])
}


[@react.component]
let make = (~value) => {
  <div className=Styles.container>
    value->MarkedJS.marked->MarkedJS.parse
  </div>
};
