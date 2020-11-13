module RenderDesktop = {
  module Styles = {
    open Css;

    let nav = isActive =>
      style([
        padding2(~v=`px(16), ~h=`zero),
        cursor(`pointer),
        fontSize(`px(12)),
        hover([color(Colors.gray6)]),
        active([color(Colors.gray6)]),
        transition(~duration=400, "all"),
        color(isActive ? Colors.gray7 : Colors.gray6),
        borderBottom(`px(4), `solid, isActive ? Colors.bandBlue : Colors.white),
      ]);
  };

  [@react.component]
  let make = (~routes) => {
    let currentRoute = ReasonReactRouter.useUrl() |> Route.fromUrl;

    <div className={CssHelper.flexBox(~justify=`spaceBetween, ())} id="navigationBar">
      {routes
       ->Belt.List.map(((v, route)) =>
           <div key=v className={CssHelper.flexBox(~justify=`spaceBetween, ())}>
             <Link className={Styles.nav(currentRoute == route)} route>
               {v |> React.string}
             </Link>
           </div>
         )
       ->Array.of_list
       ->React.array}
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
        boxShadow(
          Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, `num(0.08))),
        ),
      ]);

    let nav = style([color(Colors.gray8), padding2(~v=`px(16), ~h=`zero)]);
    let menuContainer =
      style([padding4(~top=`px(10), ~bottom=`px(10), ~left=`px(10), ~right=`px(5))]);
    let menu = style([width(`px(20)), display(`block)]);
    let cmcLogo = style([width(`px(19)), height(`px(19))]);
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
        backgroundColor(`rgba((0, 0, 0, `num(0.5)))),
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
      <div className=Styles.menuContainer onClick={_ => setShow(prev => !prev)}>
        <img src={show ? Images.close : Images.menu} className=Styles.menu />
      </div>
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
              <Icon name="fab fa-twitter" color=Colors.bandBlue size=20 />
            </a>
          </div>
          <HSpacing size={`px(24)} />
          <div className=Styles.socialLink>
            <a href="https://t.me/bandprotocol" target="_blank" rel="noopener">
              <Icon name="fab fa-telegram-plane" color=Colors.bandBlue size=21 />
            </a>
          </div>
          <HSpacing size={`px(24)} />
          <div className=Styles.socialLink>
            <a
              href="https://coinmarketcap.com/currencies/band-protocol"
              target="_blank"
              rel="noopener">
              <img src=Images.cmcLogo className=Styles.cmcLogo />
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
    ("Proposals", ProposalHomePage),
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
