module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row)]);

  let pageContainer =
    style([
      paddingTop(`px(50)),
      minHeight(`px(500)),
      alignItems(`center),
      justifyContent(`center),
      backgroundColor(Colors.white),
    ]);

  let logo = style([width(`px(100)), marginRight(`px(10))]);
};

[@react.component]
let make = () => {
  <div className=Styles.pageContainer>
    <Col> <img src=Images.notFoundBg className=Styles.logo /> </Col>
  </div>;
};
