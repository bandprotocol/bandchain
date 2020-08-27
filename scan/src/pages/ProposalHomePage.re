module Styles = {
  open Css;
  let infoContainer =
    style([
      backgroundColor(Colors.white),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
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
  let badge = color => {
    style([backgroundColor(color), padding2(~v=`px(3), ~h=`px(10)), borderRadius(`px(50))]);
  };
};

let getBadgeText =
  fun
  | ProposalSub.Deposit => "Deposit Period"
  | Voting => "Voting Period"
  | Passed => "Passed"
  | Rejected => "Rejected"
  | Failed => "Failed";

let getBadgeColor =
  fun
  | ProposalSub.Deposit
  | Voting => Colors.bandBlue
  | Passed => Colors.green4
  | Rejected
  | Failed => Colors.red4;

module ProposalCard = {
  [@react.component]
  let make = (~reserveIndex, ~proposalSub: ApolloHooks.Subscription.variant(ProposalSub.t)) => {
    let isMobile = Media.isMobile();
    <Col.Grid key={reserveIndex |> string_of_int} mb=24 mbSm=16>
      <div className=Styles.infoContainer>
        <Row.Grid marginBottom=18>
          <Col.Grid col=Col.Eight>
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
          </Col.Grid>
          <Col.Grid col=Col.Four>
            <div
              className={Css.merge([
                CssHelper.flexBox(~justify=`flexEnd, ()),
                Styles.badgeContainer,
              ])}>
              {switch (proposalSub) {
               | Data({status}) =>
                 <div
                   className={Css.merge([
                     Styles.badge(getBadgeColor(status)),
                     CssHelper.flexBox(~justify=`center, ()),
                   ])}>
                   <Text value={getBadgeText(status)} color=Colors.white />
                 </div>
               | _ => <LoadingCensorBar width=100 height=15 radius=50 />
               }}
            </div>
          </Col.Grid>
        </Row.Grid>
        <Row.Grid marginBottom=24>
          <Col.Grid>
            {switch (proposalSub) {
             | Data({description}) => <Text value=description size=Text.Lg block=true />
             | _ => <LoadingCensorBar width=270 height=15 />
             }}
          </Col.Grid>
        </Row.Grid>
        <Row.Grid>
          <Col.Grid col=Col.Four mbSm=16>
            <Heading value="Proposer" size=Heading.H5 marginBottom=8 />
            {switch (proposalSub) {
             | Data({proposerAddress}) =>
               <AddressRender address=proposerAddress position=AddressRender.Subtitle />
             | _ => <LoadingCensorBar width=270 height=15 />
             }}
          </Col.Grid>
          <Col.Grid col=Col.Four colSm=Col.Seven>
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
          </Col.Grid>
          {switch (proposalSub) {
           | Data({status, turnout}) =>
             switch (status) {
             | Deposit => React.null
             | Voting
             | Passed
             | Rejected
             | Failed =>
               <Col.Grid col=Col.Four colSm=Col.Five>
                 <Heading value="Turnout" size=Heading.H5 marginBottom=8 />
                 <Text value={turnout |> Format.fPercent(~digits=2)} size=Text.Lg />
               </Col.Grid>
             }
           | _ =>
             <Col.Grid col=Col.Four colSm=Col.Five>
               <LoadingCensorBar width=100 height=15 mb=8 />
               <LoadingCensorBar width=50 height=15 />
             </Col.Grid>
           }}
        </Row.Grid>
      </div>
    </Col.Grid>;
  };
};

[@react.component]
let make = () => {
  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 10;
  let proposalsSub = ProposalSub.getList(~pageSize, ~page, ());
  let proposalsCountSub = ProposalSub.count();
  let allSub = Sub.all2(proposalsSub, proposalsCountSub);
  let isMobile = Media.isMobile();

  <Section>
    <div className=CssHelper.container>
      <Row.Grid alignItems=Row.Center marginBottom=40 marginBottomSm=24>
        <Col.Grid col=Col.Twelve>
          <Heading value="All Proposals" size=Heading.H2 marginBottom=40 marginBottomSm=24 />
          {switch (allSub) {
           | Data((_, proposalsCount)) =>
             <Heading value={proposalsCount->string_of_int ++ " In total"} size=Heading.H3 />
           | _ => <LoadingCensorBar width=65 height=21 />
           }}
        </Col.Grid>
      </Row.Grid>
      <Row.Grid>
        {switch (allSub) {
         | Data((proposals, proposalsCount)) =>
           let pageCount = Page.getPageCount(proposalsCount, pageSize);
           <>
             {proposals
              ->Belt_Array.mapWithIndex((i, e) =>
                  <ProposalCard
                    key={i |> string_of_int}
                    reserveIndex=i
                    proposalSub={Sub.resolve(e)}
                  />
                )
              ->React.array}
             {isMobile
                ? React.null
                : <Pagination
                    currentPage=page
                    pageCount
                    onPageChange={newPage => setPage(_ => newPage)}
                  />}
           </>;
         | _ =>
           Belt_Array.make(10, ApolloHooks.Subscription.NoData)
           ->Belt_Array.mapWithIndex((i, noData) =>
               <ProposalCard key={i |> string_of_int} reserveIndex=i proposalSub=noData />
             )
           ->React.array
         }}
      </Row.Grid>
    </div>
  </Section>;
};
