module Styles = {
  open Css;
  let addressWrapper = style([width(`px(120)), Media.smallMobile([width(`px(80))])]);
};

[@react.component]
let make = (~msgType, ~name, ~fromAddress) => {
  <div className={Css.merge([CssHelper.flexBox(~wrap=`nowrap, ())])}>
    <div className=Styles.addressWrapper> <AddressRender address=fromAddress /> </div>
    <MsgBadge msgType name />
  </div>;
};
