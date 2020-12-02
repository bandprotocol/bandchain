module Styles = {
  open Css;

  let container =
    style([
      display(`flex),
      justifyContent(`center),
      width(`px(640)),
      height(`px(275)),
      position(`relative),
    ]);

  let bg =
    style([
      position(`absolute),
      width(`percent(100.)),
      height(`percent(100.)),
      backgroundColor(Css.rgb(249, 249, 251)),
      backgroundImage(`url(Images.waitingModalBg)),
      backgroundRepeat(`noRepeat),
      borderRadius(`px(8)),
      zIndex(-1),
    ]);

  let innerContainer =
    style([
      display(`flex),
      flexDirection(`column),
      width(`percent(100.)),
      justifyContent(`center),
      alignItems(`center),
    ]);

  let modalTitle = style([display(`flex), justifyContent(`center)]);

  let icon = style([width(`px(30))]);

  let errorContainer =
    style([
      display(`flex),
      flexDirection(`row),
      justifyContent(`center),
      alignItems(`center),
    ]);
};

type result_t =
  | Loading
  | Error
  | Success;

[@react.component]
let make = _ => {
  // TODO: remove later
  let x = Js.Math.random_int(0, 3);
  let result =
    switch (x) {
    | 0 => Loading
    | 1 => Error
    | _ => Success
    };

  <div className=Styles.container>
    <div className=Styles.bg />
    <div className=Styles.innerContainer>
      {switch (result) {
       | Loading =>
         <>
           <Text
             value="Waiting for user to sign the transaction"
             weight=Text.Bold
             size=Text.Xxxl
             color=Colors.gray8
           />
           <VSpacing size=Spacing.xxl />
           <Loading width={`px(100)} />
         </>
       | Success =>
         <>
           <Text
             value="User Signed Successfully"
             weight=Text.Bold
             size=Text.Xxxl
             color=Colors.gray8
           />
           <VSpacing size=Spacing.xxl />
           <img src=Images.success2 className=Styles.icon />
         </>
       | Error =>
         <>
           <Text value="Signing Fail" weight=Text.Bold size=Text.Xxxl color=Colors.gray8 />
           <VSpacing size=Spacing.xxl />
           <div className=Styles.errorContainer>
             <img src=Images.fail2 className=Styles.icon />
             <HSpacing size=Spacing.md />
             <Text
               value="An error occurred"
               weight=Text.Medium
               color=Colors.red4
               spacing={Text.Em(0.03)}
             />
           </div>
         </>
       }}
    </div>
  </div>;
};
