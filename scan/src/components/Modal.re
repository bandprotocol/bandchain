module Styles = {
  open Css;

  let overlay = isFadeOut =>
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
      animation(
        ~duration=Config.modalFadingDutation,
        ~timingFunction=`cubicBezier((0.25, 0.46, 0.45, 0.94)),
        ~fillMode=`forwards,
        keyframes(
          isFadeOut
            ? [(0, [opacity(1.)]), (100, [opacity(0.)])]
            : [(0, [opacity(0.)]), (100, [opacity(1.)])],
        ),
      ),
    ]);

  let content = isFadeOut =>
    style([
      display(`table),
      marginTop(`vw(10.)),
      backgroundColor(Css_Colors.white),
      borderRadius(`px(5)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(8), ~blur=`px(32), Css.rgba(0, 0, 0, 0.5))),
      animation(
        ~duration=Config.modalFadingDutation,
        ~timingFunction=`cubicBezier((0.25, 0.46, 0.45, 0.94)),
        ~fillMode=`forwards,
        keyframes(
          isFadeOut
            ? [
              (0, [transform(translateY(`px(-30))), opacity(1.)]),
              (100, [transform(translateY(`zero)), opacity(0.)]),
            ]
            : [
              (0, [transform(translateY(`zero)), opacity(0.)]),
              (100, [transform(translateY(`px(-30))), opacity(1.)]),
            ],
        ),
      ),
    ]);

  let closeButton =
    style([
      width(`px(15)),
      position(`absolute),
      top(`px(20)),
      left(`px(605)),
      cursor(`pointer),
      zIndex(3),
    ]);
};

[@react.component]
let make = () => {
  let (modalStateOpt, dispatchModal) = React.useContext(ModalContext.context);

  let closeModal = () => {
    dispatchModal(CloseModal);
  };

  switch (modalStateOpt) {
  | None => React.null
  | Some({modal, canExit, closing}) =>
    <div className={Styles.overlay(closing)} onClick={_ => {canExit ? closeModal() : ()}}>
      <div
        className={Styles.content(closing)} onClick={e => ReactEvent.Mouse.stopPropagation(e)}>
        <img
          src=Images.closeButton
          onClick={_ => {canExit ? closeModal() : ()}}
          className=Styles.closeButton
        />
        {switch (modal) {
         | Connect(chainID) => <ConnectModal chainID />
         | SubmitTx(msg) => <SubmitTxModal msg />
         }}
      </div>
    </div>
  };
};
