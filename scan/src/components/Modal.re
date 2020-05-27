module Styles = {
  open Css;

  let overlay =
    style([
      display(`flex),
      justifyContent(`center),
      position(`absolute),
      width(`percent(100.)),
      height(`percent(100.)),
      top(`zero),
      left(`zero),
      bottom(`zero),
      right(`zero),
      zIndex(10),
      backgroundColor(`rgba((0, 0, 0, 0.5))),
      position(`fixed),
    ]);

  let content =
    style([
      display(`table),
      marginTop(`vw(10.)),
      backgroundColor(Css_Colors.white),
      borderRadius(`px(5)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(8), ~blur=`px(32), Css.rgba(0, 0, 0, 0.5))),
      animation(
        ~duration=500,
        ~timingFunction=`cubicBezier((0.25, 0.46, 0.45, 0.94)),
        ~fillMode=`forwards,
        keyframes([
          (0, [transform(translateY(`zero)), opacity(0.)]),
          (100, [transform(translateY(`px(-30))), opacity(1.)]),
        ]),
      ),
    ]);

  let closeButton =
    style([
      width(`px(15)),
      position(`absolute),
      top(`px(20)),
      left(`px(605)),
      cursor(`pointer),
    ]);
};

[@react.component]
let make = () => {
  let (modalStateOpt, dispatchModal) = React.useContext(ModalContext.context);

  switch (modalStateOpt) {
  | None => React.null
  | Some({modal, canExit}) =>
    let body =
      switch (modal) {
      | Connect(chainID) => <ConnectModal chainID />
      | SubmitTx(msg) => <SubmitTxModal msg />
      };
    <div className=Styles.overlay onClick={_ => {canExit ? dispatchModal(CloseModal) : ()}}>
      <div className=Styles.content onClick={e => ReactEvent.Mouse.stopPropagation(e)}>
        <img
          src=Images.closeButton
          onClick={_ => {canExit ? dispatchModal(CloseModal) : ()}}
          className=Styles.closeButton
        />
        body
      </div>
    </div>;
  };
};
