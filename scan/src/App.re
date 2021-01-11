module Styles = {
  open Css;

  let container =
    style([width(`percent(100.)), height(`percent(100.)), position(`relative)]);

  let innerContainer =
    style([
      marginLeft(`auto),
      marginRight(`auto),
      padding2(~v=`zero, ~h=`px(15)),
      Media.mobile([marginTop(`px(58))]),
    ]);

  let routeContainer =
    style([
      minHeight(`calc((`sub, `vh(100.), `px(200)))),
      paddingBottom(`px(20)),
      Media.mobile([paddingBottom(`zero)]),
    ]);
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
  let currentRoute = ReasonReactRouter.useUrl() |> Route.fromUrl;
  let (syncing, setSyncing) = React.useState(_ => false);
  let (_, dispatchModal) = React.useContext(ModalContext.context);
  let trackingSub = TrackingSub.use();

  // If database is syncing the state (when replayOffset = -2).
  React.useEffect2(
    () => {
      switch (trackingSub) {
      | Data({replayOffset}) when replayOffset != (-2) && !syncing =>
        Syncing->OpenModal->dispatchModal;
        setSyncing(_ => true);
      | _ => ()
      };
      None;
    },
    (trackingSub, syncing),
  );

  <div className=Styles.container>
    <Header />
    {Media.isMobile()
       ? <Section pt=16 pb=16 bg={currentRoute == HomePage ? Colors.highlightBg : Colors.bg}>
           <div className=CssHelper.container> <SearchBar /> </div>
         </Section>
       : React.null}
    <div className=Styles.routeContainer>
      {switch (currentRoute) {
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
       | ProposalHomePage => <ProposalHomePage />
       | ProposalIndexPage(proposalID) =>
         <ProposalIndexPage proposalID={ID.Proposal.ID(proposalID)} />
       | IBCHomePage => <IBCHomePage />
       | NotFound => <NotFound />
       }}
    </div>
    <Modal />
  </div>;
};
