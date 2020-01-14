module Styles = {
  open Css;

  let container = style([display(`flex), flexDirection(`column)]);

  let listContainer =
    style([
      background(Colors.white),
      display(`inlineFlex),
      padding4(~top=`px(6), ~right=`px(17), ~bottom=`px(6), ~left=`px(17)),
      border(`px(1), `solid, `hex("EEEEEE")),
    ]);

  let keyContainer = style([width(`px(60)), marginRight(`px(15))]);
};

[@react.component]
let make = () =>
  <div className=Styles.container>
    {[("interval", "7-day"), ("method", "median")]
     ->Belt.List.mapWithIndex((idx, (key, value)) =>
         <div className=Styles.listContainer key={idx |> string_of_int}>
           <div className=Styles.keyContainer> <Text value=key size=Text.Lg /> </div>
           <Text value size=Text.Lg weight=Text.Semibold />
         </div>
       )
     ->Array.of_list
     ->React.array}
  </div>;
