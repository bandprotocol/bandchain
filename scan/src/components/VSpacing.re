open Css;

[@react.component]
let make = (~size) => {
  <div className={style([paddingTop(size)])} />;
};