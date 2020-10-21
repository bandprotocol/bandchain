module Styles = {
  open Css;

  let loading = (w, h, mb) =>
    style([
      width(w),
      height(h),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      marginBottom(mb),
    ]);
};

[@react.component]
let make = (~width=`px(65), ~height=`px(20), ~marginBottom=`unset) => {
  <img
    src=Images.loadingCircles
    className={Css.merge([Styles.loading(width, height, marginBottom)])}
  />;
};
