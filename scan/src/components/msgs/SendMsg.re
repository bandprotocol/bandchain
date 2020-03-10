module Styles = {
  open Css;
  let fullWidth = style([width(`percent(100.0)), display(`flex)]);
  let hashContainer = style([maxWidth(`px(140))]);
  let statusContainer =
    style([maxWidth(`px(95)), display(`flex), flexDirection(`row), alignItems(`center)]);
  let logo = style([width(`px(20)), marginLeft(`auto), marginRight(`px(15))]);
};

[@react.component]
let make = (~send: TxHook.Msg.Send.t) => {
  <div> <Text value={send.fromAddress |> Address.toHex} /> </div>;
};
