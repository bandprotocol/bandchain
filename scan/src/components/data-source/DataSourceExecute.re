module Styles = {
  open Css;

  let container = style([padding2(~h=`px(35), ~v=`px(20))]);

  let paramsContainer = style([display(`flex), flexDirection(`column)]);

  let listContainer = style([marginBottom(`px(25))]);

  let input =
    style([
      width(`percent(100.)),
      background(white),
      paddingLeft(`px(25)),
      fontSize(`px(12)),
      fontWeight(`num(500)),
      outline(`px(1), `none, white),
      height(`px(40)),
      borderRadius(`px(4)),
      boxShadow(
        Shadow.box(~inset=true, ~x=`zero, ~y=`zero, ~blur=`px(4), Css.rgba(0, 0, 0, 0.1)),
      ),
    ]);

  let buttonContainer = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let button =
    style([
      width(`px(110)),
      backgroundColor(Colors.btnGreen),
      borderRadius(`px(6)),
      fontSize(`px(12)),
      fontWeight(`num(500)),
      color(`hex("1D7C73")),
      cursor(`pointer),
      padding2(~v=Css.px(10), ~h=Css.px(10)),
      whiteSpace(`nowrap),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.1))),
      border(`zero, `solid, Colors.white),
    ]);

  let hFlex = style([display(`flex), flexDirection(`row)]);
};

let parameterInput = (name, placeholder) => {
  <div className=Styles.listContainer key=name>
    <Text value=name size=Text.Md color=Colors.graySubHeader />
    <VSpacing size=Spacing.xs />
    <input className=Styles.input type_="text" placeholder />
  </div>;
};

[@react.component]
let make = () => {
  <div className=Styles.container>
    <div className=Styles.hFlex>
      <Text value="Test data source execution with following" color=Colors.grayHeader />
      <HSpacing size=Spacing.sm />
      <Text value="parameters" color=Colors.grayHeader weight=Text.Bold />
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.paramsContainer>
      {parameterInput("Symbol", "BTC")}
      {parameterInput("Multiplier", "100")}
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.buttonContainer>
      <button className=Styles.button> {"Test Execution" |> React.string} </button>
      <HSpacing size=Spacing.xl />
    </div>
  </div>;
};
