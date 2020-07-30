module RenderDesktop = {
  module Styles = {
    open Css;

    let navContainer =
      style([
        paddingTop(Spacing.md),
        paddingBottom(Spacing.md),
        maxWidth(`px(970)),
        marginLeft(`auto),
        marginRight(`auto),
        minHeight(`px(70)),
      ]);

    let nav =
      style([
        paddingRight(Spacing.md),
        cursor(`pointer),
        color(Colors.blueGray4),
        fontSize(`px(11)),
        hover([color(Colors.blueGray4)]),
        active([color(Colors.blueGray4)]),
      ]);
  };

  [@react.component]
  let make = (~routes) => {
    <div className=Styles.navContainer>
      <Row justify=Row.Between alignItems=`flexStart>
        <Col>
          <Row>
            {routes
             ->Belt.List.map(((v, route)) =>
                 <Col key=v> <Link className=Styles.nav route> {v |> React.string} </Link> </Col>
               )
             ->Array.of_list
             ->React.array}
          </Row>
        </Col>
        <Col> <UserAccount /> </Col>
      </Row>
    </div>;
  };
};

module RenderMobile = {
  module Styles = {
    open Css;

    let navContainer = show =>
      style([
        display(`flex),
        flexDirection(`column),
        opacity(show ? 1. : 0.),
        zIndex(2),
        pointerEvents(show ? `auto : `none),
        width(`percent(100.)),
        position(`absolute),
        top(`px(62)),
        left(`zero),
        right(`zero),
        transition(~duration=400, "all"),
        backgroundColor(Colors.white),
        padding4(~top=`zero, ~left=`px(24), ~right=`px(24), ~bottom=`px(24)),
        boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
      ]);

    let nav = style([color(Colors.gray8), padding2(~v=`px(16), ~h=`zero)]);
    let menu = style([width(`px(20))]);
    let twitterLogo = style([width(`px(19))]);
    let telegramLogo = style([width(`px(17))]);
    let socialContainer = style([display(`flex), flexDirection(`row), marginTop(`px(10))]);
    let socialLink =
      style([
        display(`flex),
        flexDirection(`row),
        justifyContent(`center),
        alignItems(`center),
      ]);
    let backdropContainer = show =>
      style([
        width(`percent(100.)),
        height(`percent(100.)),
        backgroundColor(`rgba((0, 0, 0, 0.5))),
        position(`fixed),
        opacity(show ? 1. : 0.),
        pointerEvents(show ? `auto : `none),
        left(`zero),
        top(`px(62)),
        transition(~duration=400, "all"),
      ]);
  };

  [@react.component]
  let make = (~routes) => {
    let (show, setShow) = React.useState(_ => false);
    <>
      <img
        src={show ? Images.close : Images.menu}
        className=Styles.menu
        onClick={_ => setShow(prev => !prev)}
      />
      <div className={Styles.navContainer(show)}>
        {routes
         ->Belt.List.map(((v, route)) =>
             <Link key=v className=Styles.nav route onClick={_ => setShow(_ => false)}>
               <Text value=v size=Text.Lg />
             </Link>
           )
         ->Array.of_list
         ->React.array}
        <div className=Styles.socialContainer>
          <div className=Styles.socialLink>
            <a href="https://twitter.com/bandprotocol" target="_blank" rel="noopener">
              <img src=Images.twitterLogo className=Styles.twitterLogo />
            </a>
          </div>
          <HSpacing size={`px(24)} />
          <div className=Styles.socialLink>
            <a href="https://t.me/bandprotocol" target="_blank" rel="noopener">
              <img src=Images.telegramLogo className=Styles.telegramLogo />
            </a>
          </div>
        </div>
      </div>
      <div onClick={_ => setShow(prev => !prev)} className={Styles.backdropContainer(show)} />
    </>;
  };
};

[@react.component]
let make = () => {
  let wenchangRoutes = [
    ("Home", Route.HomePage),
    ("Validators", ValidatorHomePage),
    ("Blocks", BlockHomePage),
    ("Transactions", TxHomePage),
  ];

  exception WrongNetwork(string);
  let routes =
    switch (Env.network) {
    | "WENCHANG" => wenchangRoutes
    | "GUANYU38"
    | "GUANYU" =>
      wenchangRoutes->Belt.List.concat([
        ("Data Sources", DataSourceHomePage),
        ("Oracle Scripts", OracleScriptHomePage),
        ("Requests", RequestHomePage),
      ])
    | _ => raise(WrongNetwork("Incorrect or unspecified NETWORK environment variable"))
    };

  Media.isMobile() ? <RenderMobile routes /> : <RenderDesktop routes />;
};
