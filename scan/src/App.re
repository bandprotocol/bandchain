module Styles = {
  open Css;

  let pageWidth = style([maxWidth(`px(Config.pageWidth))]);

  let container =
    style([width(`percent(100.)), height(`percent(100.)), position(`relative)]);

  let innerContainer =
    style([marginLeft(`auto), marginRight(`auto), padding2(~v=`px(0), ~h=`px(15))]);

  let routeContainer =
    style([minHeight(`calc((`sub, `vh(100.), `px(200)))), paddingBottom(`px(20))]);
};

[@react.component]
let make = () => {
  exception WrongNetwork(string);
  switch (Env.network) {
  | "WENCHANG"
  | "GUANYU38"
  | "GUANYU" => ()
  | _ => raise(WrongNetwork("Incorrect or unspecified NETWORK environment variable"))
  };

  <div className=Styles.container>
    <TopBar />
    <div className={Css.merge([Styles.innerContainer, Styles.pageWidth])}>
      {Media.isMobile() ? React.null : <NavBar />}
      <div className=Styles.routeContainer>
        {switch (ReasonReactRouter.useUrl() |> Route.fromUrl) {
         | HomePage => <HomePage />
         | DataSourceHomePage => <DataSourceHomePage />
         | DataSourceIndexPage(dataSourceID, hashtag) =>
           <DataSourceIndexPage dataSourceID={ID.DataSource.ID(dataSourceID)} hashtag />
         | OracleScriptHomePage => <OracleScriptHomePage />
         | OracleScriptIndexPage(oracleScriptID, hashtag) =>
           <OracleScriptIndexPage oracleScriptID={ID.OracleScript.ID(oracleScriptID)} hashtag />
         | TxHomePage => <TxHomePage />
         | TxIndexPage(txHash) => <TxIndexPage txHash />
         | BlockHomePage => <BlockHomePage />
         | BlockIndexPage(height) => <BlockIndexPage height={ID.Block.ID(height)} />
         | ValidatorHomePage => <ValidatorHomePage />
         | ValidatorIndexPage(address, hashtag) => <ValidatorIndexPage address hashtag />
         | RequestHomePage => <RequestHomePage />
         | RequestIndexPage(reqID) => <RequestIndexPage reqID={ID.Request.ID(reqID)} />
         | AccountIndexPage(address, hashtag) => <AccountIndexPage address hashtag />
         | IBCHomePage => <IBCHomePage />
         | NotFound => <NotFound />
         }}
      </div>
    </div>
    <Modal />
  </div>;
};
