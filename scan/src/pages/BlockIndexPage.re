module Styles = {
  open Css;

  let vFlex =
    style([
      display(`flex),
      flexDirection(`column),
      alignItems(`center),
      justifyContent(`center),
      height(`px(600)),
    ]);
};

[@react.component]
let make = (~height, ~hashtag) => {
  <div className=Styles.vFlex>
    <Text value="Block Index Page" size=Text.Xxl weight=Text.Bold nowrap=true />
    <Text value={j|block height is $height|j} size=Text.Lg weight=Text.Bold nowrap=true />
    <Text
      value={hashtag != "" ? {j|Hashtag is $hashtag|j} : "No Hashtag"}
      size=Text.Xl
      weight=Text.Bold
      nowrap=true
    />
  </div>;
};
