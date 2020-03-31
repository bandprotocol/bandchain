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
      backgroundColor(Colors.blue1),
      borderRadius(`px(15)),
      color(textColor),
      backgroundColor(bgColor),
      display(`inlineFlex),
    ]);

  let msgAmount =
    style([
      borderRadius(`percent(50.)),
      backgroundColor(Colors.gray4),
      marginLeft(`px(5)),
      display(`inlineFlex),
      alignItems(`center),
      justifyContent(`center),
      width(`px(14)),
      height(`px(14)),
    ]);

  let txTypeMapping = msg => {
    switch (msg) {
    | TxHook.Msg.CreateDataSource(_) => ("CREATE DATA SOURCE", Colors.green7, Colors.green1)
    | TxHook.Msg.EditDataSource(_) => ("EDIT DATA SOURCE", Colors.green7, Colors.green1)
    | TxHook.Msg.CreateOracleScript(_) => ("CREATE ORACLE SCRIPT", Colors.green7, Colors.green1)
    | TxHook.Msg.EditOracleScript(_) => ("EDIT ORACLE SCRIPT", Colors.green7, Colors.green1)
    | TxHook.Msg.Request(_) => ("DATA REQUEST", Colors.blue4, Colors.blue1)
    | TxHook.Msg.Send(_) => ("SEND TOKEN", Colors.purple6, Colors.purple1)
    | TxHook.Msg.Report(_) => ("DATA REPORT", Colors.orange6, Colors.orange1)
    | TxHook.Msg.AddOracleAddress(_) => ("ADD ORACLE ADDRESS", Colors.purple6, Colors.purple1)
    | TxHook.Msg.RemoveOracleAddress(_) => (
        "REMOVE ORACLE ADDRESS",
        Colors.purple6,
        Colors.purple1,
      )
    | TxHook.Msg.CreateValidator(_) => ("CREATE VALIDATOR", Colors.purple6, Colors.purple1)
    | TxHook.Msg.EditValidator(_) => ("EDIT VALIDATOR", Colors.purple6, Colors.purple1)
    | Unknown => ("Unknown", Colors.gray6, Colors.gray7)
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
