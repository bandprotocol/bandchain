module Styles = {
  open Css;

  let container =
    style([
      display(`flex),
      position(`relative),
      width(`percent(100.)),
      height(`px(7)),
      borderRadius(`px(7)),
      backgroundColor(Colors.blueGray2),
    ]);

  let outer =
    style([
      display(`flex),
      width(`percent(100.)),
      height(`px(7)),
      borderRadius(`px(7)),
      overflow(`hidden),
    ]);

  let inner = p =>
    style([
      display(`flex),
      width(`percent(p)),
      height(`px(7)),
      background(Colors.bandBlue),
    ]);

  let withWidth = (w, p) =>
    style([
      display(`flex),
      flexDirection(`column),
      position(`absolute),
      width(`px(w)),
      top(`px(-22)),
      left(`calc((`sub, `percent(p), `px(w / 2)))),
      alignItems(`center),
    ]);

  let arrowDown =
    style([
      width(`zero),
      height(`zero),
      borderLeft(`px(6), `solid, Colors.transparent),
      borderRight(`px(6), `solid, Colors.transparent),
      borderTop(`px(6), `solid, Colors.bandBlue),
    ]);

  let lowerText =
    style([
      position(`absolute),
      display(`flex),
      width(`percent(100.)),
      justifyContent(`center),
      top(`px(10)),
    ]);

  // Modern
  let barContainer =
    style([
      position(`relative),
      paddingTop(`px(20)),
      Media.mobile([display(`flex), alignItems(`center), paddingTop(`zero)]),
    ]);
  let progressOuter =
    style([
      position(`relative),
      width(`percent(100.)),
      height(`px(12)),
      borderRadius(`px(7)),
      border(`px(1), `solid, Colors.gray9),
      padding(`px(1)),
      overflow(`hidden),
    ]);
  let progressInner = (p, success) =>
    style([
      width(`percent(p)),
      height(`percent(100.)),
      borderRadius(`px(7)),
      background(success ? Colors.bandBlue : Colors.red4),
    ]);
  let leftText =
    style([
      position(`absolute),
      top(`zero),
      left(`zero),
      Media.mobile([
        position(`static),
        flexGrow(0.),
        flexShrink(0.),
        flexBasis(`px(50)),
        paddingRight(`px(10)),
      ]),
    ]);
  let rightText =
    style([
      position(`absolute),
      top(`zero),
      right(`zero),
      Media.mobile([
        position(`static),
        flexGrow(0.),
        flexShrink(0.),
        flexBasis(`px(70)),
        paddingLeft(`px(10)),
      ]),
    ]);
};

[@react.component]
let make = (~reportedValidators, ~minimumValidators, ~requestValidators) => {
  let minimumPercentage =
    (minimumValidators * 100 |> float_of_int) /. (requestValidators |> float_of_int);
  let progressPercentage =
    (reportedValidators * 100 |> float_of_int) /. (requestValidators |> float_of_int);
  <div className=Styles.container>
    <div className=Styles.outer> <div className={Styles.inner(progressPercentage)} /> </div>
    <div className={Styles.withWidth(120, minimumPercentage)}>
      <div className=Styles.arrowDown />
    </div>
    {reportedValidators < minimumValidators
       ? <div className=Styles.lowerText>
           <Text
             value={
               "Pending "
               ++ (minimumValidators - reportedValidators |> Format.iPretty)
               ++ " Validators"
             }
             color=Colors.gray7
             size=Text.Sm
           />
         </div>
       : React.null}
  </div>;
};

module Modern = {
  [@react.component]
  let make = (~reportedValidators, ~minimumValidators, ~requestValidators) => {
    let progressPercentage =
      (reportedValidators * 100 |> float_of_int) /. (requestValidators |> float_of_int);
    let success = reportedValidators >= minimumValidators ? true : false;

    <div className=Styles.barContainer>
      <div className=Styles.leftText>
        <Text value={"Min " ++ (minimumValidators |> Format.iPretty)} color=Colors.gray7 />
      </div>
      <div className=Styles.progressOuter>
        <div className={Styles.progressInner(progressPercentage, success)} />
      </div>
      <div className=Styles.rightText>
        <Text
          value={
            (reportedValidators |> Format.iPretty)
            ++ " of "
            ++ (requestValidators |> Format.iPretty)
          }
          color=Colors.gray7
        />
      </div>
    </div>;
  };
};
