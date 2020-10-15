module Styles = {
  open Css;
  let buttonContainer =
    style([selector("> button", [margin(`zero)]), Media.mobile([width(`percent(100.))])]);
  let baseBtn =
    style([
      border(`px(1), `solid, Colors.bandBlue),
      textAlign(`center),
      Media.mobile([flexGrow(0.), flexShrink(0.), flexBasis(`percent(50.))]),
    ]);

  let activeBtn = isActive => {
    style([
      borderTopRightRadius(`zero),
      borderBottomRightRadius(`zero),
      backgroundColor(isActive ? Colors.bandBlue : Colors.white),
      hover([backgroundColor(isActive ? Colors.bandBlue : Colors.white)]),
    ]);
  };
  let inActiveBtn = isActive => {
    style([
      borderTopLeftRadius(`zero),
      borderBottomLeftRadius(`zero),
      backgroundColor(isActive ? Colors.white : Colors.bandBlue),
      hover([backgroundColor(isActive ? Colors.white : Colors.bandBlue)]),
    ]);
  };
};

[@react.component]
let make = (~isActive, ~setIsActive) => {
  <div className={Css.merge([CssHelper.flexBox(), Styles.buttonContainer])}>
    <Button
      px=22
      py=5
      onClick={_ => setIsActive(_ => true)}
      style={Css.merge([Styles.baseBtn, Styles.activeBtn(isActive)])}>
      <Text value="Active" color={isActive ? Colors.white : Colors.bandBlue} />
    </Button>
    <Button
      px=22
      py=5
      onClick={_ => setIsActive(_ => false)}
      style={Css.merge([Styles.baseBtn, Styles.inActiveBtn(isActive)])}>
      <Text value="Inactive" color={isActive ? Colors.bandBlue : Colors.white} />
    </Button>
  </div>;
};
