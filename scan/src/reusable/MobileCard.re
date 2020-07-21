module Styles = {
  open Css;
  let cardContainer =
    style([
      padding(`px(10)),
      backgroundColor(Colors.white),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
      selector("+ div", [marginTop(`px(10))]),
    ]);
  let cardItem =
    style([
      display(`flex),
      alignItems(`center),
      padding(`px(5)),
    ]);
  let cardItemHeading =
    style([display(`flex), alignItems(`center), flexBasis(`percent(25.))]);
};

[@react.component]
let make = (~values, ~idx) => {
  <div className=Styles.cardContainer>
    {values
     ->Belt_List.mapWithIndex((index, (heading, value)) => {
         <div className=Styles.cardItem key={idx ++ (index |> string_of_int)}>
           <div className=Styles.cardItemHeading>
             <Text
               value=heading
               size=Text.Xs
               weight=Text.Semibold
               color=Colors.gray6
               spacing={Text.Em(0.1)}
             />
           </div>
           <div> <InfoMobileCard info=value /> </div>
         </div>;
       })
     ->Belt.List.toArray
     ->React.array}
  </div>;
};

