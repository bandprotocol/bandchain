module Styles = {
  open Css;

  let pageContainer = style([paddingTop(`px(50))]);

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let logo = style([width(`px(27)), marginRight(`px(10))]);

  let sourceContainer = style([marginTop(`px(15))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.grayHeader),
    ]);

  let codeVerifiedBadge =
    style([
      backgroundColor(`hex("D7FFEC")),
      borderRadius(`px(6)),
      display(`inlineFlex),
      justifyContent(`center),
      alignItems(`center),
      padding4(~top=`px(10), ~bottom=`px(10), ~left=`px(13), ~right=`px(13)),
    ]);

  let checkLogo = style([marginRight(`px(10))]);

  let tableContainer = style([border(`px(1), `solid, Colors.lightGray)]);

  let tableHeader = style([backgroundColor(Colors.white), padding(`px(20))]);

  let tableLowerContainer =
    style([
      padding(`px(20)),
      backgroundImage(
        `linearGradient((
          deg(0.0),
          [(`percent(0.0), Colors.white), (`percent(100.0), Colors.lighterGray)],
        )),
      ),
    ]);
};

[@react.component]
let make = (~codeHash, ~hashtag: Route.script_tab_t) => {
  // let step = 10;
  // let (limit, setLimit) = React.useState(_ => step);
  // let scriptOpt = ScriptHook.getInfo(codeHash);
  // let (txs, totalCount) =
  //   TxHook.withCodehash(~codeHash, ~limit, ())
  //   ->Belt.Option.mapWithDefault(([], 0), ({txs, totalCount}) =>
  //       (txs |> Belt.List.reverse, totalCount)
  //     );
  <div className=Styles.pageContainer>
    <Row justify=Row.Between>
      <Col>
        <div className=Styles.vFlex>
          <img src=Images.newScript className=Styles.logo />
          <Text
            value="DATA ORACLE SCRIPT"
            weight=Text.Semibold
            size=Text.Lg
            nowrap=true
            color=Colors.grayHeader
            block=true
          />
          <HSpacing size=Spacing.sm />
          <div className=Styles.seperatedLine />
        </div>
      </Col>
    </Row>
    // {switch (scriptOpt) {
    //  | Some(script) =>
    //    <TimeAgos time={script.createdAtTime} size=Text.Lg weight=Text.Regular />
    //  | None => React.null
    //  }}
    // <Col>
    //   {switch (scriptOpt) {
    //    | Some(_) =>
    //      <div className=Styles.codeVerifiedBadge>
    //        <img src=Images.checkIcon className=Styles.checkLogo />
    //        <Text
    //          value="Code Verified"
    //          size=Text.Lg
    //          weight=Text.Semibold
    //          color=Colors.darkGreen
    //        />
    //      </div>
    //    | None => React.null
    //    }}
    // </Col>
    // <div className=Styles.sourceContainer>
    //   <Text
    //     value={
    //       switch (scriptOpt) {
    //       | Some(script) => script.info.name
    //       | None => "?"
    //       }
    //     }
    //     size=Text.Xxxl
    //     weight=Text.Bold
    //     nowrap=true
    //   />
    // </div>
    <VSpacing size=Spacing.xl />
    // <InfoHL
    //   info={
    //     InfoHL.DataSources(
    //       scriptOpt->Belt_Option.mapWithDefault([], script =>
    //         script.info.dataSources->Belt_List.map(source => source.name)
    //       ),
    //     )
    //   }
    //   header="DATA SOURCES"
    // />
    <VSpacing size=Spacing.xl />
    <Row>
      <Col>
        <InfoHL info={InfoHL.Hash(codeHash, Colors.brightPurple)} header="SCRIPT HASH" />
      </Col>
      <HSpacing size=Spacing.xl />
      <HSpacing size=Spacing.xl />
    </Row>
    // <Col>
    //   {switch (scriptOpt) {
    //    | Some(script) =>
    //      <InfoHL
    //        info={InfoHL.Address(script.info.creator, Colors.brightPurple)}
    //        header="CREATOR"
    //      />
    //    | None => React.null
    //    }}
    // </Col>
    <VSpacing size=Spacing.xl />
    // <Tab
    //   tabs=[|
    //     {name: "TRANSACTIONS", route: Route.ScriptIndexPage(codeHash, Route.ScriptTransactions)},
    //     {name: "CODE", route: Route.ScriptIndexPage(codeHash, Route.ScriptCode)},
    //     {name: "EXECUTE", route: Route.ScriptIndexPage(codeHash, Route.ScriptExecute)},
    //     {name: "INTEGRATION", route: Route.ScriptIndexPage(codeHash, Route.ScriptIntegration)},
    //   |]
    //   currentRoute={Route.ScriptIndexPage(codeHash, hashtag)}>
    //   {switch (hashtag) {
    //    | ScriptTransactions =>
    //      <div className=Styles.tableLowerContainer>
    //        <Text
    //          value={(totalCount |> Format.iPretty) ++ " Request Transactions"}
    //          color=Colors.grayHeader
    //          size=Text.Lg
    //        />
    //        <VSpacing size=Spacing.lg />
    //        <TxsTable txs />
    //        <VSpacing size=Spacing.lg />
    //        {if (txs->Belt_List.size < totalCount) {
    //           <LoadMore onClick={_ => {setLimit(oldLimit => oldLimit + step)}} />;
    //         } else {
    //           React.null;
    //         }}
    //      </div>
    //    | ScriptCode =>
    //      switch (scriptOpt) {
    //      | Some(script) => <ScriptCode codeHash params={script.info.params} />
    //      | None => <div />
    //      }
    //    | ScriptExecute =>
    //      switch (scriptOpt) {
    //      | Some(script) => <ScriptExecute script />
    //      | None => <div />
    //      }
    //    | ScriptIntegration => <div> {"TODO2" |> React.string} </div>
    //    }}
    // </Tab>
    <VSpacing size=Spacing.xxl />
  </div>;
};
