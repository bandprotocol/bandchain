module Styles = {
  open Css;

  let container = style([display(`flex), flexDirection(`column)]);

  let listContainer =
    style([
      background(Colors.white),
      display(`inlineFlex),
      padding4(~top=`px(6), ~right=`px(17), ~bottom=`px(6), ~left=`px(17)),
      border(`px(1), `solid, Colors.lightGray),
    ]);

  let keyContainer = style([display(`inlineFlex), marginRight(`px(15))]);
};

[@react.component]
let make = (~result) =>
  <div className=Styles.container>
    {result
     ->Belt.Array.mapWithIndex((idx, (key, value)) =>
         <div className=Styles.listContainer key={idx |> string_of_int}>
           <div className=Styles.keyContainer> <Text value=key size=Text.Lg /> </div>
           <Text value={value |> Js.Json.stringify} size=Text.Lg weight=Text.Semibold />
         </div>
       )
     ->React.array}
  </div>;
