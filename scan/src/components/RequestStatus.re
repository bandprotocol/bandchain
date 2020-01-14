module Styles = {
  open Css;

  let badge =
    style([
      borderRadius(`px(6)),
      display(`inlineFlex),
      justifyContent(`center),
      alignItems(`center),
      padding4(~top=`px(11), ~right=`px(17), ~bottom=`px(11), ~left=`px(17)),
    ]);

  let comfirmed = style([backgroundColor(`hex("D7FFEC"))]);
  let logo = style([marginRight(`px(10))]);

  let pending = style([backgroundColor(Colors.lighterBlue)]);
};

[@react.component]
let make = (~reqID) => {
  let comfirmed = false; /* use hook right here */
  comfirmed
    ? <div className={Css.merge([Styles.badge, Styles.comfirmed])}>
        <img src=Images.checkIcon className=Styles.logo />
        <Text value="Comfirmed" size=Text.Lg color=Colors.darkGreen />
      </div>
    : <div className={Css.merge([Styles.badge, Styles.pending])}>
        <img src=Images.pendingIcon className=Styles.logo />
        <Text value="Pending Data Reports (3/4 Providers)" size=Text.Lg color=Colors.purpleBlue />
      </div>;
};
