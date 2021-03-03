module Styles = {
  open Css;

  let iconWrapper = msgType =>
    style([
      backgroundColor(
        switch (msgType) {
        | TxSub.Msg.TokenMsg => Colors.bandBlue
        | ValidatorMsg => Colors.blue12
        | ProposalMsg => Colors.blue13
        | DataMsg => Colors.blue14
        | IBCMsg => Colors.blue13
        | _ => Colors.bandBlue
        },
      ),
      width(`px(24)),
      height(`px(24)),
      borderRadius(`percent(50.)),
      position(`relative),
      selector(
        "> i",
        [
          position(`absolute),
          left(`percent(50.)),
          top(`percent(50.)),
          transform(translate(`percent(-50.), `percent(-50.))),
        ],
      ),
    ]);
};

[@react.component]
let make = (~category: TxSub.Msg.msg_cat_t) => {
  <div className={Styles.iconWrapper(category)}>
    {switch (category) {
     | TokenMsg => <Icon name="far fa-wallet" color=Colors.white size=14 />
     | ValidatorMsg => <Icon name="fas fa-user" color=Colors.white size=14 />
     | ProposalMsg => <Icon name="fal fa-file" color=Colors.white size=14 />
     | DataMsg => <Icon name="fal fa-globe" color=Colors.white size=14 />
     | IBCMsg => <Icon name="fal fa-exchange-alt" color=Colors.white size=14 />
     | _ => <Icon name="fal fa-question-circle" color=Colors.white size=14 />
     }}
  </div>;
};
