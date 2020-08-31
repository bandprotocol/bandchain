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
  let infoHeader =
    style([borderBottom(`px(1), `solid, Colors.gray9), paddingBottom(`px(16))]);
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
      <Row.Grid>
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
            <Row.Grid marginBottom=24>
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
    </div>
  </Section>;
};
