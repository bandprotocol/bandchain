module Styles = {
  open Css;
  let infoContainer =
    style([
      backgroundColor(Colors.white),
      boxShadow(
        Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, `num(0.08))),
      ),
      padding(`px(24)),
      position(`relative),
      Media.mobile([padding(`px(16))]),
    ]);
  let idContainer = {
    style([
      selector(
        "> h4",
        [marginLeft(`px(6)), Media.mobile([marginLeft(`zero), marginTop(`px(16))])],
      ),
    ]);
  };
  let badgeContainer = {
    style([Media.mobile([position(`absolute), right(`px(16)), top(`px(16))])]);
  };

  let noDataImage = style([width(`auto), height(`px(70)), marginBottom(`px(16))]);
};

module ProposalCard = {
  [@react.component]
  let make =
      (
        ~reserveIndex,
        ~proposalSub: ApolloHooks.Subscription.variant(ProposalSub.t),
        ~turnoutRate,
      ) => {
    let isMobile = Media.isMobile();
    <Col key={reserveIndex |> string_of_int} mb=24 mbSm=16>
      <div className=Styles.infoContainer>
        <Row marginBottom=18>
          <Col col=Col.Eight>
            <div
              className={Css.merge([
                CssHelper.flexBox(),
                CssHelper.flexBoxSm(~direction=`column, ~align=`flexStart, ()),
                Styles.idContainer,
              ])}>
              {switch (proposalSub) {
               | Data({id, name}) =>
                 <>
                   <TypeID.Proposal id position=TypeID.Subtitle />
                   <Heading size=Heading.H4 value=name />
                 </>
               | _ =>
                 isMobile
                   ? <>
                       <LoadingCensorBar width=50 height=15 mbSm=16 />
                       <LoadingCensorBar width=150 height=15 mbSm=16 />
                     </>
                   : <LoadingCensorBar width=270 height=15 />
               }}
            </div>
          </Col>
          <Col col=Col.Four>
            <div
              className={Css.merge([
                CssHelper.flexBox(~justify=`flexEnd, ()),
                Styles.badgeContainer,
              ])}>
              {switch (proposalSub) {
               | Data({status}) => <ProposalBadge status />
               | _ => <LoadingCensorBar width=100 height=15 radius=50 />
               }}
            </div>
          </Col>
        </Row>
        <Row marginBottom=24>
          <Col>
            {switch (proposalSub) {
             | Data({description}) => <Markdown value=description />
             | _ => <LoadingCensorBar width=270 height=15 />
             }}
          </Col>
        </Row>
        <Row>
          <Col col=Col.Four mbSm=16>
            <Heading value="Proposer" size=Heading.H5 marginBottom=8 />
            {switch (proposalSub) {
             | Data({proposerAddressOpt}) =>
               switch (proposerAddressOpt) {
               | Some(proposerAddress) =>
                 <AddressRender address=proposerAddress position=AddressRender.Subtitle />
               | None => <Text value="Proposed on Wenchang" />
               }
             | _ => <LoadingCensorBar width=270 height=15 />
             }}
          </Col>
          <Col col=Col.Four colSm=Col.Seven>
            <div className={CssHelper.mb(~size=8, ())}>
              {switch (proposalSub) {
               | Data({status}) =>
                 <Heading
                   value={
                     switch (status) {
                     | Deposit => "Deposit End Time"
                     | Voting
                     | Passed
                     | Rejected
                     | Failed => "Voting End Time"
                     }
                   }
                   size=Heading.H5
                 />
               | _ => <LoadingCensorBar width=100 height=15 />
               }}
            </div>
            {switch (proposalSub) {
             | Data({depositEndTime, votingEndTime, status}) =>
               <Timestamp
                 size=Text.Lg
                 time={
                   switch (status) {
                   | Deposit => depositEndTime
                   | Voting
                   | Passed
                   | Rejected
                   | Failed => votingEndTime
                   }
                 }
               />
             | _ => <LoadingCensorBar width={isMobile ? 120 : 270} height=15 />
             }}
          </Col>
          {switch (proposalSub) {
           | Data({status}) =>
             switch (status) {
             | Deposit => React.null
             | Voting
             | Passed
             | Rejected
             | Failed =>
               <Col col=Col.Four colSm=Col.Five>
                 <Heading value="Turnout" size=Heading.H5 marginBottom=8 />
                 <Text value={turnoutRate |> Format.fPercent(~digits=2)} size=Text.Lg />
               </Col>
             }
           | _ =>
             <Col col=Col.Four colSm=Col.Five>
               <LoadingCensorBar width=100 height=15 mb=8 />
               <LoadingCensorBar width=50 height=15 />
             </Col>
           }}
        </Row>
      </div>
    </Col>;
  };
};

[@react.component]
let make = () => {
  let pageSize = 10;
  let proposalsSub = ProposalSub.getList(~pageSize, ~page=1, ());
  let voteStatSub = VoteSub.getVoteStats();
  let bondedTokenCountSub = ValidatorSub.getTotalBondedAmount();

  let allSub = Sub.all3(proposalsSub, bondedTokenCountSub, voteStatSub);

  <Section>
    <div className=CssHelper.container id="proposalsSection">
      <Row alignItems=Row.Center marginBottom=40 marginBottomSm=24>
        <Col col=Col.Twelve> <Heading value="All Proposals" size=Heading.H2 /> </Col>
      </Row>
      <Row>
        {switch (allSub) {
         | Data((proposals, bondedTokenCount, voteStatSub)) =>
           proposals->Belt.Array.size > 0
             ? proposals
               ->Belt_Array.mapWithIndex((i, proposal) => {
                   let turnoutRate =
                     proposal.id == ID.Proposal.ID(1)
                       ? 71.59
                       : (
                           voteStatSub->Belt_MapInt.get(proposal.id |> ID.Proposal.toInt)
                           |> Belt_Option.getWithDefault(_, 0.)
                         )
                         /. (bondedTokenCount |> Coin.getBandAmountFromCoin)
                         *. 100.;
                   <ProposalCard
                     key={i |> string_of_int}
                     reserveIndex=i
                     proposalSub={Sub.resolve(proposal)}
                     turnoutRate
                   />;
                 })
               ->React.array
             : <EmptyContainer>
                 <img src=Images.noSource className=Styles.noDataImage />
                 <Heading
                   size=Heading.H4
                   value="No Proposal"
                   align=Heading.Center
                   weight=Heading.Regular
                   color=Colors.bandBlue
                 />
               </EmptyContainer>
         | _ =>
           Belt_Array.make(pageSize, ApolloHooks.Subscription.NoData)
           ->Belt_Array.mapWithIndex((i, noData) =>
               <ProposalCard
                 key={i |> string_of_int}
                 reserveIndex=i
                 proposalSub=noData
                 turnoutRate=0.
               />
             )
           ->React.array
         }}
      </Row>
    </div>
  </Section>;
};
