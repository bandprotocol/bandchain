module Styles = {
  open Css;

  let loading = (w, mb) =>
    style([
      width(w),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      marginBottom(mb),
    ]);
};

[@react.component]
let make = (~width=`px(65), ~marginBottom=`unset) => {
  <img src=Images.loadingCircles className={Css.merge([Styles.loading(width, marginBottom)])} />;
};
