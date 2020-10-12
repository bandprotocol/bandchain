module Styles = {
  open Css;
  let infoContainer =
    style([
      backgroundColor(Colors.white),
      boxShadow(
        Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, `num(0.08))),
      ),
      padding(`px(24)),
      height(`percent(100.)),
      position(`relative),
      Media.mobile([padding(`px(16))]),
    ]);
  let resultsInfoContainer = style([paddingTop(`px(17))]);

  let infoHeader =
    style([borderBottom(`px(1), `solid, Colors.gray9), paddingBottom(`px(16))]);

  let resultsInfoheader =
    style([paddingBottom(`px(13)), Media.mobile([paddingBottom(`px(16))])]);

  let tableContainer =
    style([
      backgroundColor(Colors.white),
      boxShadow(
        Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, `num(0.08))),
      ),
      Media.mobile([margin2(~h=`px(-12), ~v=`zero), padding2(~h=`px(16), ~v=`zero)]),
    ]);

  let tableHeader =
    style([
      padding2(~v=`px(16), ~h=`px(24)),
      Media.mobile([padding2(~v=`px(14), ~h=`px(12))]),
    ]);

  let statusLogo = style([width(`px(20))]);
  let resultContainer =
    style([
      margin3(~top=`px(40), ~h=`zero, ~bottom=`px(20)),
      selector("> div + div", [marginTop(`px(24))]),
    ]);
  let separatorLine =
    style([
      borderStyle(`none),
      backgroundColor(Colors.gray9),
      height(`px(1)),
      margin2(~v=`px(24), ~h=`auto),
    ]);
  let voteButton =
    fun
    | ProposalSub.Voting => style([visibility(`visible)])
    | Deposit
    | Passed
    | Rejected
    | Failed => style([visibility(`hidden)]);
};

module VoteButton = {
  [@react.component]
  let make = (~proposalID, ~proposalName) => {
    let trackingSub = TrackingSub.use();

    let (accountOpt, _) = React.useContext(AccountContext.context);
    let (_, dispatchModal) = React.useContext(ModalContext.context);

    let connect = chainID => dispatchModal(OpenModal(Connect(chainID)));
    let vote = () => Vote(proposalID, proposalName)->SubmitTx->OpenModal->dispatchModal;

    switch (accountOpt) {
    | Some(_) =>
      <Button px=20 py=5 style={CssHelper.flexBox()} onClick={_ => vote()}>
        <Icon name="fas fa-pencil" size=12 color=Colors.white />
        <HSpacing size=Spacing.sm />
        <Text value="Vote" nowrap=true block=true />
      </Button>
    | None =>
      switch (trackingSub) {
      | Data({chainID}) =>
        <Button px=20 py=5 style={CssHelper.flexBox()} onClick={_ => connect(chainID)}>
          <Icon name="fas fa-pencil" size=12 color=Colors.white />
          <HSpacing size=Spacing.sm />
          <Text value="Vote" nowrap=true block=true />
        </Button>
      | Error(err) =>
        // log for err details
        Js.Console.log(err);
        <Text value="chain id not found" />;
      | _ => <LoadingCensorBar width=90 height=26 />
      }
    };
  };
};

[@react.component]
let make = (~proposalID) => {
  let proposalSub = ProposalSub.get(proposalID);
  let voteStatByProposalIDSub = VoteSub.getVoteStatByProposalID(proposalID);
  let bondedTokenCountSub = ValidatorSub.getTotalBondedAmount();

  let allSub = Sub.all3(proposalSub, voteStatByProposalIDSub, bondedTokenCountSub);
  let isMobile = Media.isMobile();

  <Section pbSm=0>
    <div className=CssHelper.container>
      <Row.Grid marginBottom=40 marginBottomSm=16>
        <Col.Grid>
          <Heading value="Proposal" size=Heading.H4 marginBottom=40 marginBottomSm=24 />
        </Col.Grid>
        <Col.Grid col=Col.Eight mbSm=16>
          {switch (allSub) {
           | Data(({id, name}, _, _)) =>
             <div className={CssHelper.flexBox()}>
               <TypeID.Proposal id position=TypeID.Title />
               <HSpacing size=Spacing.sm />
               <Heading size=Heading.H3 value=name />
             </div>
           | _ => <LoadingCensorBar width=270 height=15 />
           }}
        </Col.Grid>
        <Col.Grid col=Col.Four>
          <div
            className={Css.merge([
              CssHelper.flexBox(~justify=`flexEnd, ()),
              CssHelper.flexBoxSm(~justify=`flexStart, ()),
            ])}>
            {switch (allSub) {
             | Data(({status}, _, _)) => <ProposalBadge status />
             | _ => <LoadingCensorBar width=100 height=15 radius=50 />
             }}
          </div>
        </Col.Grid>
      </Row.Grid>
      <Row.Grid marginBottom=24>
        <Col.Grid>
          <div className=Styles.infoContainer>
            <Row.Grid>
              <Col.Grid>
                <Heading
                  value="Information"
                  size=Heading.H4
                  style=Styles.infoHeader
                  marginBottom=24
                />
              </Col.Grid>
              <Col.Grid col=Col.Six mb=24>
                <Heading value="Proposer" size=Heading.H5 marginBottom=8 />
                {switch (allSub) {
                 | Data(({proposerAddressOpt}, _, _)) =>
                   switch (proposerAddressOpt) {
                   | Some(proposerAddress) =>
                     <AddressRender address=proposerAddress position=AddressRender.Subtitle />
                   | None => <Text value="Proposed on Wenchang" />
                   }
                 | _ => <LoadingCensorBar width=270 height=15 />
                 }}
              </Col.Grid>
              <Col.Grid col=Col.Six mb=24>
                <Heading value="Submit Time" size=Heading.H5 marginBottom=8 />
                {switch (allSub) {
                 | Data(({submitTime}, _, _)) => <Timestamp size=Text.Lg time=submitTime />
                 | _ => <LoadingCensorBar width={isMobile ? 120 : 270} height=15 />
                 }}
              </Col.Grid>
              <Col.Grid col=Col.Six mb=24>
                <Heading value="Proposal Type" size=Heading.H5 marginBottom=8 />
                {switch (allSub) {
                 | Data(({proposalType}, _, _)) =>
                   <Text value=proposalType size=Text.Lg block=true />
                 | _ => <LoadingCensorBar width=90 height=15 />
                 }}
              </Col.Grid>
            </Row.Grid>
            <Row.Grid>
              <Col.Grid>
                <Heading value="Description" size=Heading.H5 marginBottom=8 />
                {switch (allSub) {
                 | Data(({description}, _, _)) => <Markdown value=description />
                 | _ => <LoadingCensorBar width=270 height=15 />
                 }}
              </Col.Grid>
            </Row.Grid>
          </div>
        </Col.Grid>
      </Row.Grid>
      {switch (allSub) {
       | Data((
           {status, name, votingStartTime, votingEndTime},
           {
             total,
             totalYes,
             totalYesPercent,
             totalNo,
             totalNoPercent,
             totalNoWithVeto,
             totalNoWithVetoPercent,
             totalAbstain,
             totalAbstainPercent,
           },
           bondedToken,
         )) =>
         switch (status) {
         | Deposit => React.null
         | Voting
         | Passed
         | Rejected
         | Failed =>
           <>
             <Row.Grid>
               <Col.Grid col=Col.Six mb=24 mbSm=16>
                 <div className=Styles.infoContainer>
                   <Heading
                     value="Voting Overview"
                     size=Heading.H4
                     style=Styles.infoHeader
                     marginBottom=24
                   />
                   <Row.Grid marginTop=38 alignItems=Row.Center>
                     <Col.Grid col=Col.Seven>
                       <div className={CssHelper.flexBoxSm(~justify=`spaceAround, ())}>
                         {let turnoutPercent =
                            total /. (bondedToken |> Coin.getBandAmountFromCoin) *. 100.;
                          <TurnoutChart percent=turnoutPercent />}
                         {isMobile
                            ? <div>
                                <Heading value="Total" size=Heading.H5 marginBottom=8 />
                                <Text
                                  value={(total |> Format.fPretty(~digits=2)) ++ " BAND"}
                                  size=Text.Lg
                                  block=true
                                  color=Colors.gray6
                                />
                              </div>
                            : React.null}
                       </div>
                     </Col.Grid>
                     {isMobile
                        ? <Col.Grid> <hr className=Styles.separatorLine /> </Col.Grid> : React.null}
                     <Col.Grid col=Col.Five>
                       <Row.Grid>
                         {isMobile
                            ? React.null
                            : <Col.Grid mb=24>
                                <Heading value="Total" size=Heading.H5 marginBottom=8 />
                                <Text
                                  value={(total |> Format.fPretty(~digits=2)) ++ " BAND"}
                                  size=Text.Lg
                                  block=true
                                  color=Colors.gray6
                                />
                              </Col.Grid>}
                         <Col.Grid mb=24 mbSm=0 colSm=Col.Six>
                           <Heading value="Voting Start" size=Heading.H5 marginBottom=8 />
                           <Timestamp.Grid size=Text.Lg time=votingStartTime />
                         </Col.Grid>
                         <Col.Grid mb=24 mbSm=0 colSm=Col.Six>
                           <Heading value="Voting End" size=Heading.H5 marginBottom=8 />
                           <Timestamp.Grid size=Text.Lg time=votingEndTime />
                         </Col.Grid>
                       </Row.Grid>
                     </Col.Grid>
                   </Row.Grid>
                 </div>
               </Col.Grid>
               <Col.Grid col=Col.Six mb=24 mbSm=16>
                 <div className={Css.merge([Styles.infoContainer, Styles.resultsInfoContainer])}>
                   <div
                     className={Css.merge([
                       CssHelper.flexBox(~justify=`spaceBetween, ()),
                       Styles.infoHeader,
                       Styles.resultsInfoheader,
                       CssHelper.mb(~size=24, ()),
                     ])}>
                     <Heading value="Results" size=Heading.H4 />
                     {isMobile
                        ? React.null
                        : <div className={Styles.voteButton(status)}>
                            <VoteButton proposalID proposalName=name />
                          </div>}
                   </div>
                   <div className=Styles.resultContainer>
                     <>
                       <ProgressBar.Voting
                         label=VoteSub.Yes
                         amount=totalYes
                         percent=totalYesPercent
                       />
                       <ProgressBar.Voting
                         label=VoteSub.No
                         amount=totalNo
                         percent=totalNoPercent
                       />
                       <ProgressBar.Voting
                         label=VoteSub.NoWithVeto
                         amount=totalNoWithVeto
                         percent=totalNoWithVetoPercent
                       />
                       <ProgressBar.Voting
                         label=VoteSub.Abstain
                         amount=totalAbstain
                         percent=totalAbstainPercent
                       />
                     </>
                   </div>
                 </div>
               </Col.Grid>
             </Row.Grid>
             <Row.Grid marginBottom=24>
               <Col.Grid> <VoteBreakdownTable proposalID /> </Col.Grid>
             </Row.Grid>
           </>
         }
       | _ => React.null
       }}
      <Row.Grid marginBottom=24>
        <Col.Grid>
          <div className=Styles.infoContainer>
            <Row.Grid>
              <Col.Grid>
                <Heading value="Deposit" size=Heading.H4 style=Styles.infoHeader marginBottom=24 />
              </Col.Grid>
              <Col.Grid col=Col.Six mbSm=24>
                <Heading value="Deposit Status" size=Heading.H5 marginBottom=8 />
                {switch (proposalSub) {
                 | Data({totalDeposit, status}) =>
                   switch (status) {
                   | ProposalSub.Deposit => <ProgressBar.Deposit totalDeposit />
                   | _ =>
                     <div className={CssHelper.flexBox()}>
                       <img src=Images.success className=Styles.statusLogo />
                       <HSpacing size=Spacing.sm />
                       // TODO: remove hard-coded later
                       <Text value="Completed Min Deposit 1,000 BAND" />
                     </div>
                   }
                 | _ => <LoadingCensorBar width={isMobile ? 120 : 270} height=15 />
                 }}
              </Col.Grid>
              <Col.Grid col=Col.Six>
                <Heading value="Deposit End Time" size=Heading.H5 marginBottom=8 />
                {switch (proposalSub) {
                 | Data({depositEndTime}) => <Timestamp size=Text.Lg time=depositEndTime />
                 | _ => <LoadingCensorBar width=90 height=15 />
                 }}
              </Col.Grid>
            </Row.Grid>
          </div>
        </Col.Grid>
      </Row.Grid>
      <div className=Styles.tableContainer>
        <Row.Grid>
          <Col.Grid>
            <Heading value="Depositors" size=Heading.H4 style=Styles.tableHeader />
          </Col.Grid>
          <Col.Grid> <DepositorTable proposalID /> </Col.Grid>
        </Row.Grid>
      </div>
    </div>
  </Section>;
};
