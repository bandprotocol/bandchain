module Styles = {
  open Css;

  let container =
    style([
      backgroundColor(Colors.white),
      boxShadow(
        Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, `num(0.08))),
      ),
      Media.mobile([margin2(~h=`px(-12), ~v=`zero)]),
    ]);

  let header =
    style([
      backgroundColor(Colors.white),
      padding2(~v=`zero, ~h=`px(24)),
      borderBottom(`px(1), `solid, Colors.gray4),
      selector("> div + div", [marginLeft(`px(40))]),
      Media.mobile([overflow(`auto), padding2(~v=`px(1), ~h=`px(15))]),
    ]);

  let buttonContainer = active =>
    style([
      height(`px(40)),
      display(`inlineFlex),
      justifyContent(`center),
      alignItems(`center),
      cursor(`pointer),
      padding3(~top=`px(24), ~h=`zero, ~bottom=`px(20)),
      borderBottom(`px(4), `solid, active ? Colors.bandBlue : Colors.white),
      Media.mobile([whiteSpace(`nowrap)]),
    ]);

  let childrenContainer =
    style([
      backgroundColor(Colors.blueGray1),
      Media.mobile([padding2(~h=`px(16), ~v=`zero)]),
    ]);

  let tableWrapper = style([Media.mobile([padding2(~v=`px(16), ~h=`zero)])]);
  let noDataImage = style([width(`auto), height(`px(70)), marginBottom(`px(16))]);
};

let renderBody = (reserveIndex, voteSub: ApolloHooks.Subscription.variant(VoteSub.t)) => {
  <TBody.Grid
    key={
      switch (voteSub) {
      | Data({voter}) => voter |> Address.toBech32
      | _ => reserveIndex |> string_of_int
      }
    }
    paddingH={`px(24)}>
    <Row.Grid alignItems=Row.Center>
      <Col.Grid col=Col.Five>
        {switch (voteSub) {
         | Data({voter, validator}) =>
           switch (validator) {
           | Some({moniker, operatorAddress, identity}) =>
             <ValidatorMonikerLink
               validatorAddress=operatorAddress
               moniker
               identity
               width={`percent(100.)}
               avatarWidth=20
             />
           | None => <AddressRender address=voter />
           }
         | _ => <LoadingCensorBar width=200 height=20 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Four>
        {switch (voteSub) {
         | Data({txHashOpt}) =>
           switch (txHashOpt) {
           | Some(txHash) => <TxLink txHash width=200 />
           | None => <Text value="Voted on Wenchang" />
           }
         | _ => <LoadingCensorBar width=200 height=20 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Three>
        <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
          {switch (voteSub) {
           | Data({timestampOpt}) =>
             switch (timestampOpt) {
             | Some(timestamp) =>
               <Timestamp.Grid
                 time=timestamp
                 size=Text.Md
                 weight=Text.Regular
                 textAlign=Text.Right
               />
             | None => <Text value="Created on Wenchang" />
             }
           | _ => <LoadingCensorBar width=80 height=15 />
           }}
        </div>
      </Col.Grid>
    </Row.Grid>
  </TBody.Grid>;
};

let renderBodyMobile = (reserveIndex, voteSub: ApolloHooks.Subscription.variant(VoteSub.t)) => {
  switch (voteSub) {
  | Data({voter, txHashOpt, timestampOpt, validator}) =>
    let key_ = voter |> Address.toBech32;

    <MobileCard
      values=InfoMobileCard.[
        (
          "Voter",
          {switch (validator) {
           | Some({operatorAddress, moniker, identity}) =>
             Validator(operatorAddress, moniker, identity)
           | None => Address(voter, 200, `account)
           }},
        ),
        (
          "TX Hash",
          switch (txHashOpt) {
          | Some(txHash) => TxHash(txHash, 200)
          | None => Text("Voted on Wenchang")
          },
        ),
        (
          "Timestamp",
          switch (timestampOpt) {
          | Some(timestamp) => Timestamp(timestamp)
          | None => Text("Created on Wenchang")
          },
        ),
      ]
      key=key_
      idx=key_
    />;
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("Voter", Loading(230)),
        ("TX Hash", Loading(100)),
        ("Timestamp", Loading(100)),
      ]
      key={reserveIndex |> string_of_int}
      idx={reserveIndex |> string_of_int}
    />
  };
};

module TabButton = {
  [@react.component]
  let make = (~tab, ~active, ~setTab) => {
    let tabString = tab |> VoteSub.toString(~withSpace=true);

    <div className={Styles.buttonContainer(active)} onClick={_ => setTab(_ => tab)}>
      <Text
        value=tabString
        weight={active ? Text.Semibold : Text.Regular}
        size=Text.Lg
        color=Colors.gray6
      />
    </div>;
  };
};

[@react.component]
let make = (~proposalID) => {
  let isMobile = Media.isMobile();
  let (currentTab, setTab) = React.useState(_ => VoteSub.Yes);
  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 5;
  let votesSub = VoteSub.getList(proposalID, currentTab, ~pageSize, ~page, ());
  let voteCountSub = VoteSub.count(proposalID, currentTab);

  <div className=Styles.container>
    <div className={Css.merge([Styles.header, CssHelper.flexBox(~wrap=`nowrap, ())])}>
      {[|VoteSub.Yes, No, NoWithVeto, Abstain|]
       ->Belt.Array.map(tab =>
           <TabButton key={tab |> VoteSub.toString} tab setTab active={tab == currentTab} />
         )
       ->React.array}
    </div>
    <div className=Styles.childrenContainer>
      <div className=Styles.tableWrapper>
        {isMobile
           ? <Row.Grid marginBottom=16>
               <Col.Grid>
                 {switch (voteCountSub) {
                  | Data(voteCount) =>
                    <div className={CssHelper.flexBox()}>
                      <Text
                        block=true
                        value={voteCount |> string_of_int}
                        weight=Text.Semibold
                        color=Colors.gray7
                      />
                      <HSpacing size=Spacing.xs />
                      <Text block=true value="Voters" weight=Text.Semibold color=Colors.gray7 />
                    </div>
                  | _ => <LoadingCensorBar width=100 height=15 />
                  }}
               </Col.Grid>
             </Row.Grid>
           : <THead.Grid>
               <Row.Grid alignItems=Row.Center>
                 <Col.Grid col=Col.Five>
                   {switch (voteCountSub) {
                    | Data(voteCount) =>
                      <div className={CssHelper.flexBox()}>
                        <Text
                          block=true
                          value={voteCount |> string_of_int}
                          weight=Text.Semibold
                          color=Colors.gray7
                        />
                        <HSpacing size=Spacing.xs />
                        <Text block=true value="Voters" weight=Text.Semibold color=Colors.gray7 />
                      </div>
                    | _ => <LoadingCensorBar width=100 height=15 />
                    }}
                 </Col.Grid>
                 <Col.Grid col=Col.Four>
                   <Text block=true value="TX Hash" weight=Text.Semibold color=Colors.gray7 />
                 </Col.Grid>
                 <Col.Grid col=Col.Three>
                   <Text
                     block=true
                     value="Timestamp"
                     weight=Text.Semibold
                     color=Colors.gray7
                     align=Text.Right
                   />
                 </Col.Grid>
               </Row.Grid>
             </THead.Grid>}
        {switch (votesSub) {
         | Data(votes) =>
           votes->Belt.Array.size > 0
             ? votes
               ->Belt_Array.mapWithIndex((i, e) =>
                   isMobile
                     ? renderBodyMobile(i, Sub.resolve(e)) : renderBody(i, Sub.resolve(e))
                 )
               ->React.array
             : <EmptyContainer height={`px(250)}>
                 <img src=Images.noAccount className=Styles.noDataImage />
                 <Heading
                   size=Heading.H4
                   value="No Voters"
                   align=Heading.Center
                   weight=Heading.Regular
                   color=Colors.bandBlue
                 />
               </EmptyContainer>
         | _ =>
           Belt_Array.make(pageSize, ApolloHooks.Subscription.NoData)
           ->Belt_Array.mapWithIndex((i, noData) =>
               isMobile ? renderBodyMobile(i, noData) : renderBody(i, noData)
             )
           ->React.array
         }}
        {switch (voteCountSub) {
         | Data(voteCount) =>
           let pageCount = Page.getPageCount(voteCount, pageSize);
           <Pagination
             currentPage=page
             pageCount
             onPageChange={newPage => setPage(_ => newPage)}
           />;
         | _ => React.null
         }}
      </div>
    </div>
  </div>;
};
