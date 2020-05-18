module Styles = {
  open Css;
  let msgContainer = overflowed =>
    style([
      position(`relative),
      height(overflowed ? `px(60) : `auto),
      overflow(overflowed ? `hidden : `visible),
    ]);
  let showButton =
    style([
      display(`flex),
      backgroundColor(Colors.gray3),
      borderRadius(`px(30)),
      width(`px(60)),
      alignItems(`center),
      justifyContent(`center),
      fontSize(`px(10)),
      alignItems(`center),
      cursor(`pointer),
      color(Colors.gray8),
      height(`px(20)),
    ]);
  let showContainer = style([display(`flex), marginTop(`px(10))]);
};

[@react.component]
let make = (~txHash: Hash.t, ~messages, ~width: int, ~success: bool, ~rawLog: string) => {
  let (overflowed, setOverflowed) = React.useState(_ => false);
  let (expanded, setExpanded) = React.useState(_ => false);

  let divID = "messageWrapper" ++ (txHash |> Hash.toHex);

  React.useEffect0(_ => {
    let x = ReactDOMRe._getElementById(divID) |> Belt.Option.getExn;
    let divHeight = ReactDOMRe.domElementToObj(x)##clientHeight;
    Js.Console.log(divHeight);
    divHeight > 60 ? setOverflowed(_ => true) : ();
    None;
  });
  <>
    <div id=divID className={Styles.msgContainer(overflowed)}>
      {messages
       ->Belt_List.toArray
       ->Belt_Array.mapWithIndex((i, msg) =>
           <React.Fragment key={(txHash |> Hash.toHex) ++ (i |> string_of_int)}>
             <VSpacing size=Spacing.sm />
             <Msg msg width />
             <VSpacing size=Spacing.sm />
           </React.Fragment>
         )
       ->React.array}
      {success
         ? React.null
         : <div> <Text value={"Error: " ++ rawLog} code=true size=Text.Sm breakAll=true /> </div>}
    </div>
    {overflowed || expanded
       ? <div>
           <div
             className=Styles.showContainer
             onClick={_ => {
               setOverflowed(_ => !overflowed);
               setExpanded(_ => !expanded);
             }}>
             {expanded
                ? <div className=Styles.showButton> {"show less" |> React.string} </div>
                : <div className=Styles.showButton> {"..." |> React.string} </div>}
           </div>
         </div>
       : React.null}
  </>;
};
