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
    ]);

  let content =
    style([
      display(`flex),
      marginTop(`vw(10.)),
      width(`px(640)),
      height(`px(480)),
      backgroundColor(Css_Colors.white),
      borderRadius(`px(5)),
      boxShadow(Css.Shadow.box(~x=`zero, ~y=`px(-10), ~blur=`px(100), `rgba((0, 0, 0, 0.3)))),
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
};

[@react.component]
let make = () => {
  let (modal, dispatchModal) = React.useContext(ModalContext.context);

  switch (modal) {
  | None => React.null
  | Some(m) =>
    let body =
      switch (m) {
      | Connect(value) => <ConnectModal value />
      };
    <div className=Styles.overlay onClick={_ => dispatchModal(CloseModal)}>
      <div className=Styles.content onClick={e => ReactEvent.Mouse.stopPropagation(e)}>
        body
      </div>
    </div>;
  };
};
