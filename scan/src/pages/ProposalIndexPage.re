module Styles = {
  open Css;
  let infoContainer =
    style([
      backgroundColor(Colors.white),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
      padding(`px(24)),
      height(`percent(100.)),
      position(`relative),
      Media.mobile([padding(`px(16))]),
    ]);
  let infoHeader =
    style([borderBottom(`px(1), `solid, Colors.gray9), paddingBottom(`px(16))]);

  let tableContainer =
    style([
      backgroundColor(Colors.white),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
      Media.mobile([margin2(~h=`px(-15), ~v=`zero), padding2(~h=`px(16), ~v=`zero)]),
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
};

[@react.component]
let make = (~proposalID) => {
  let proposalSub = ProposalSub.get(proposalID);
  let isMobile = Media.isMobile();

  <Section pbSm=0>
    <div className=CssHelper.container>
      <Row.Grid marginBottom=40 marginBottomSm=16>
        <Col.Grid>
          <Heading value="Proposal" size=Heading.H4 marginBottom=40 marginBottomSm=24 />
        </Col.Grid>
        <Col.Grid col=Col.Eight mbSm=16>
          {switch (proposalSub) {
           | Data({id, name}) =>
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
            {switch (proposalSub) {
             | Data({status}) => <ProposalBadge status />
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
                {switch (proposalSub) {
                 | Data({proposerAddress}) =>
                   <AddressRender address=proposerAddress position=AddressRender.Subtitle />
                 | _ => <LoadingCensorBar width=270 height=15 />
                 }}
              </Col.Grid>
              <Col.Grid col=Col.Six mb=24>
                <Heading value="Submit Time" size=Heading.H5 marginBottom=8 />
                {switch (proposalSub) {
                 | Data({submitTime}) => <Timestamp size=Text.Lg time=submitTime />
                 | _ => <LoadingCensorBar width={isMobile ? 120 : 270} height=15 />
                 }}
              </Col.Grid>
              <Col.Grid col=Col.Six mb=24>
                <Heading value="Proposal Type" size=Heading.H5 marginBottom=8 />
                {switch (proposalSub) {
                 | Data({proposalType}) => <Text value=proposalType size=Text.Lg block=true />
                 | _ => <LoadingCensorBar width=90 height=15 />
                 }}
              </Col.Grid>
            </Row.Grid>
            <Row.Grid>
              <Col.Grid>
                <Heading value="Description" size=Heading.H5 marginBottom=8 />
                {switch (proposalSub) {
                 | Data({description}) => <Markdown value=description />
                 | _ => <LoadingCensorBar width=270 height=15 />
                 }}
              </Col.Grid>
            </Row.Grid>
          </div>
        </Col.Grid>
      </Row.Grid>
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
                  <TurnoutChart percent=20. />
                  {isMobile
                     ? <div>
                         <Heading value="Total" size=Heading.H5 marginBottom=8 />
                         //TODO: wire up later
                         <Text
                           value={(4000000 |> Format.iPretty) ++ " BAND"}
                           size=Text.Lg
                           block=true
                           color=Colors.gray6
                         />
                       </div>
                     : React.null}
                </div>
              </Col.Grid>
              {isMobile ? <Col.Grid> <hr className=Styles.separatorLine /> </Col.Grid> : React.null}
              <Col.Grid col=Col.Five>
                <Row.Grid>
                  {isMobile
                     ? React.null
                     : <Col.Grid mb=24>
                         <Heading value="Total" size=Heading.H5 marginBottom=8 />
                         //TODO: wire up later
                         <Text
                           value={(4000000 |> Format.iPretty) ++ " BAND"}
                           size=Text.Lg
                           block=true
                           color=Colors.gray6
                         />
                       </Col.Grid>}
                  <Col.Grid mb=24 mbSm=0 colSm=Col.Six>
                    <Heading value="Voting Start" size=Heading.H5 marginBottom=8 />
                    {switch (proposalSub) {
                     | Data({votingStartTime}) =>
                       <Timestamp.Grid size=Text.Lg time=votingStartTime />
                     | _ =>
                       <>
                         <LoadingCensorBar width=70 height=15 />
                         <LoadingCensorBar width=80 height=15 mt=5 />
                       </>
                     }}
                  </Col.Grid>
                  <Col.Grid mb=24 mbSm=0 colSm=Col.Six>
                    <Heading value="Voting End" size=Heading.H5 marginBottom=8 />
                    {switch (proposalSub) {
                     | Data({votingEndTime}) => <Timestamp.Grid size=Text.Lg time=votingEndTime />
                     | _ =>
                       <>
                         <LoadingCensorBar width=70 height=15 />
                         <LoadingCensorBar width=80 height=15 mt=5 />
                       </>
                     }}
                  </Col.Grid>
                </Row.Grid>
              </Col.Grid>
            </Row.Grid>
          </div>
        </Col.Grid>
        <Col.Grid col=Col.Six mb=24 mbSm=16>
          <div className=Styles.infoContainer>
            <div
              className={Css.merge([
                CssHelper.flexBox(),
                Styles.infoHeader,
                CssHelper.mb(~size=24, ()),
              ])}>
              // TODO: will add the button on the voting issue
               <Heading value="Results" size=Heading.H4 /> </div>
            //TODO: will re-structure when the data is wired up.
            <div className=Styles.resultContainer>
              <ProgressBar.Voting label="Yes" amount=3600000 percent=90 />
              <ProgressBar.Voting label="No" amount=200000 percent=5 />
              <ProgressBar.Voting label="No with Veto" amount=120000 percent=3 />
              <ProgressBar.Voting label="Abstain" amount=80000 percent=2 />
            </div>
          </div>
        </Col.Grid>
      </Row.Grid>
      <Row.Grid marginBottom=24>
        <Col.Grid> <VoteBreakdownTable proposalID /> </Col.Grid>
      </Row.Grid>
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
