module Styles = {
  open Css;
  let msgBadge = msgType =>
    style([
      backgroundColor(
        switch (msgType) {
        | TxSub.Msg.TokenMsg => Colors.bandBlue
        | ValidatorMsg => Colors.blue12
        | ProposalMsg => Colors.blue13
        | DataMsg => Colors.blue14
        | _ => Colors.bandBlue
        },
      ),
      borderRadius(`px(50)),
      margin2(~v=`zero, ~h=`px(5)),
      padding2(~v=`px(3), ~h=`px(8)),
    ]);
};

[@react.component]
let make = (~msgType: TxSub.Msg.msg_cat_t, ~name) => {
  <div
    className={Css.merge([
      Styles.msgBadge(msgType),
      CssHelper.flexBox(~wrap=`nowrap, ~justify=`center, ()),
    ])}>
    <Text value=name size=Text.Xs color=Colors.white transform=Text.Uppercase align=Text.Center />
  </div>;
};
