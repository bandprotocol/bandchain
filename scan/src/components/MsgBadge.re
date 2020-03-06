type position_t =
  | Header
  | Body;

module Styles = {
  open Css;

  let vFlex = style([display(`inlineFlex), alignItems(`center)]);

  let getPadding =
    fun
    | Header => style([padding2(~v=`px(5), ~h=`px(12))])
    | Body => style([padding2(~v=`px(3), ~h=`px(5))]);

  let getSize =
    fun
    | Header => style([width(`px(20)), height(`px(20))])
    | Body => style([width(`px(14)), height(`px(14))]);

  let oval = (textColor, bgColor) =>
    style([
      backgroundColor(Colors.lightBlue),
      borderRadius(`px(15)),
      color(textColor),
      backgroundColor(bgColor),
      display(`inlineFlex),
    ]);

  let msgAmount =
    style([
      borderRadius(`percent(50.)),
      backgroundColor(Colors.lightGray),
      marginLeft(`px(5)),
      display(`inlineFlex),
      alignItems(`center),
      justifyContent(`center),
      width(`px(14)),
      height(`px(14)),
    ]);

  let txTypeMapping = msg => {
    switch (msg) {
    | TxHook.Msg.CreateDataSource(_) => (
        "CREATE DATA SOURCE",
        Colors.darkGreen,
        Colors.lightGreen,
      )
    | TxHook.Msg.EditDataSource(_) => ("EDIT DATA SOURCE", Colors.darkGreen, Colors.lightGreen)
    | TxHook.Msg.CreateOracleScript(_) => (
        "CREATE ORACLE SCRIPT",
        Colors.darkGreen,
        Colors.lightGreen,
      )
    | TxHook.Msg.EditOracleScript(_) => (
        "EDIT ORACLE SCRIPT",
        Colors.darkGreen,
        Colors.lightGreen,
      )
    | TxHook.Msg.Request(_) => ("DATA REQUEST", Colors.darkBlue, Colors.lightBlue)
    | TxHook.Msg.Send(_) => ("SEND TOKEN", Colors.purple, Colors.lightPurple)
    | TxHook.Msg.Report(_) => ("DATA REPORT", Colors.darkIndigo, Colors.lightIndigo)
    | Unknown => ("Unknown", Colors.darkGrayText, Colors.grayHeader)
    };
  };
};

[@react.component]
let make = (~msgs: list(TxHook.Msg.t), ~position=Body) =>
  if (msgs->Belt.List.length > 0) {
    let firstMsg = msgs->Belt.List.getExn(0);
    let (typeName, textColor, bgColor) = Styles.txTypeMapping(firstMsg.action);

    <div className=Styles.vFlex>
      <div
        className={Css.merge([Styles.oval(textColor, bgColor), Styles.getPadding(position)])}>
        <Text value=typeName size=Text.Xs block=true />
      </div>
      {msgs->Belt.List.length > 1
         ? <div className={Css.merge([Styles.msgAmount, Styles.getSize(position)])}>
             <Text
               value={"+" ++ (msgs->Belt.List.length - 1 |> string_of_int)}
               block=true
               size=Text.Xs
             />
           </div>
         : React.null}
    </div>;
  } else {
    React.null;
  };
