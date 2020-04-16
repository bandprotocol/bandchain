module Styles = {
  open Css;

  let version =
    style([
      display(`flex),
      borderRadius(`px(10)),
      backgroundColor(Colors.blue1),
      padding2(~v=`pxFloat(4.3), ~h=`px(8)),
      justifyContent(`center),
      alignItems(`center),
      marginLeft(Spacing.xs),
      marginTop(`px(1)),
    ]);

  let versionLoading =
    style([
      display(`flex),
      borderRadius(`px(10)),
      backgroundColor(Colors.blue1),
      overflow(`hidden),
      height(`px(16)),
      width(`px(60)),
      justifyContent(`center),
      alignItems(`center),
      marginLeft(Spacing.xs),
      marginTop(`px(1)),
    ]);
};

[@react.component]
let make = () =>
  {
    let metadataSub = MetadataSub.use();
    let%Sub metadata = metadataSub;
    <div className=Styles.version>
      <Text
        value={metadata.chainID}
        size=Text.Xs
        color=Colors.blue6
        nowrap=true
        weight=Text.Semibold
        spacing={Text.Em(0.03)}
      />
    </div>
    |> Sub.resolve;
  }
  |> Sub.default(
       _,
       <div className=Styles.versionLoading>
         <LoadingCensorBar width=80 height=16 colorBase=Colors.blue1 colorLighter=Colors.white />
       </div>,
     );
