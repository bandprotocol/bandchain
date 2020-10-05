module Styles = {
  open Css;
  let addressWrapper = style([width(`px(120))]);
};

[@react.component]
let make = (~msgType, ~name, ~fromAddress) => {
  <div className={Css.merge([CssHelper.flexBox(~wrap=`nowrap, ())])}>
    <div className=Styles.addressWrapper> <AddressRender address=fromAddress /> </div>
    <MsgBadge msgType={TxSub.Msg.getCatVarientbyMsgType(msgType)} name />
  </div>;
};
