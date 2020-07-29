module Styles = {
  open Css;

  let container =
    style([
      display(`flex),
      alignItems(`center),
      Media.mobile([marginTop(`px(16)), marginBottom(`px(36))]),
    ]);
  let word = style([display(`flex), cursor(`pointer)]);

  let toggle = isActive =>
    style([
      display(`flex),
      justifyContent(isActive ? `flexStart : `flexEnd),
      backgroundColor(Colors.gray2),
      borderRadius(`px(15)),
      padding2(~v=`px(2), ~h=`px(3)),
      width(`px(45)),
      cursor(`pointer),
      boxShadow(
        Shadow.box(
          ~inset=true,
          ~x=`zero,
          ~y=`zero,
          ~blur=`px(4),
          isActive ? Colors.purple2 : Colors.gray7,
        ),
      ),
      margin2(~h=`px(6), ~v=`zero),
      minHeight(`px(20)),
      Media.mobile([width(`px(60)), minHeight(`px(30)), margin2(~h=`px(12), ~v=`zero)]),
    ]);

  let imgLogo = style([width(`px(15)), Media.mobile([width(`px(24))])]);
};

[@react.component]
let make = (~isActive, ~setIsActive) => {
  <div className=Styles.container>
    <div onClick={_ => setIsActive(_ => true)} className=Styles.word>
      <Text value="Active" color=Colors.purple8 />
    </div>
    <div className={Styles.toggle(isActive)} onClick={_ => setIsActive(oldVal => !oldVal)}>
      <img
        src={isActive ? Images.activeValidatorLogo : Images.inactiveValidatorLogo}
        className=Styles.imgLogo
      />
    </div>
    <div onClick={_ => setIsActive(_ => false)} className=Styles.word>
      <Text value="Inactive" />
    </div>
  </div>;
};
