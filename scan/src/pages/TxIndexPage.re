module Styles = {
  open Css;

  let titleContainer =
    style([
      Media.mobile([
        flexDirection(`columnReverse),
        alignItems(`flexStart),
        minHeight(`px(80)),
      ]),
    ]);

  let separatorLine =
    style([
      borderStyle(`none),
      backgroundColor(Colors.gray9),
      height(`px(1)),
      margin2(~v=`px(24), ~h=`auto),
    ]);

  let successLogo = style([width(`px(20)), marginRight(`px(10))]);

  let notfoundContainer =
    style([
      width(`percent(100.)),
      minHeight(`px(450)),
      display(`flex),
      flexDirection(`column),
      alignItems(`center),
      paddingLeft(`px(50)),
      paddingRight(`px(50)),
      justifyContent(`center),
      backgroundColor(Colors.white),
      borderRadius(`px(4)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), rgba(0, 0, 0, `num(0.1)))),
    ]);
  let notfoundLogo = style([width(`px(180)), marginRight(`px(10))]);
  let infoContainerFullwidth =
    style([
      Media.mobile([
        selector("> div", [flexBasis(`percent(100.))]),
        selector("> div + div", [marginTop(`px(15))]),
        selector("> div > div > div", [display(`block)]),
      ]),
    ]);
  let infoContainerHalfwidth =
    style([
      Media.mobile([
        selector(
          "> div",
          [flexGrow(0.), flexShrink(0.), flexBasis(`calc((`sub, `percent(50.), `px(20))))],
        ),
        selector("> div + div + div", [marginTop(`px(15))]),
        selector("> div *", [alignItems(`flexStart)]),
      ]),
    ]);

  let infoContainer =
    style([
      backgroundColor(Colors.white),
      boxShadow(
        Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, `num(0.08))),
      ),
      padding(`px(24)),
      Media.mobile([padding(`px(16))]),
    ]);
  let infoHeader =
    style([borderBottom(`px(1), `solid, Colors.gray9), paddingBottom(`px(16))]);
};

module TxNotFound = {
  [@react.component]
  let make = () => {
    <Section>
      <div className=CssHelper.container>
        <VSpacing size=Spacing.lg />
        <div className=Styles.notfoundContainer>
          <Col> <img src=Images.notFoundBg className=Styles.notfoundLogo /> </Col>
          <VSpacing size=Spacing.md />
          <Text
            value="Sorry, we are unable to retrieve information on this transaction hash."
            size=Text.Lg
            color=Colors.blueGray6
          />
          <VSpacing size=Spacing.lg />
          <Text
            value="Note: Transactions usually take 5-10 seconds to appear."
            size=Text.Lg
            color=Colors.blueGray6
          />
        </div>
      </div>
    </Section>;
  };
};

[@react.component]
let make = (~txHash) => {
  let isMobile = Media.isMobile();
  let txSub = TxSub.get(txHash);
  switch (txSub) {
  | Loading
  | Data(_) =>
    <Section>
      <div className=CssHelper.container>
        <Row marginBottom=40 marginBottomSm=16>
          <Col.Grid>
            <Heading value="Transaction" size=Heading.H4 marginBottom=40 marginBottomSm=24 />
            <div
              className={Css.merge([
                CssHelper.flexBox(~justify=`spaceBetween, ()),
                Styles.titleContainer,
              ])}>
              <div className={CssHelper.flexBox()}>
                {switch (txSub) {
                 | Data(_) =>
                   isMobile
                     ? <Text
                         value={txHash |> Hash.toHex(~upper=true)}
                         size=Text.Lg
                         weight=Text.Bold
                         nowrap=false
                         breakAll=true
                         code=true
                         color=Colors.gray7
                       />
                     : <>
                         <Text
                           value={txHash |> Hash.toHex(~upper=true)}
                           size=Text.Xxl
                           weight=Text.Bold
                           nowrap=true
                           code=true
                           color=Colors.gray7
                         />
                         <HSpacing size=Spacing.sm />
                         <CopyRender width=15 message={txHash |> Hash.toHex(~upper=true)} />
                       </>
                 | _ => <LoadingCensorBar width=270 height=15 />
                 }}
              </div>
              <div className={CssHelper.flexBox()}>
                {switch (txSub) {
                 | Data({success}) =>
                   <>
                     <img
                       src={success ? Images.success : Images.fail}
                       className=Styles.successLogo
                     />
                     <Text
                       value={success ? "Success" : "Failed"}
                       nowrap=true
                       size=Text.Lg
                       color=Colors.gray7
                       block=true
                     />
                   </>
                 | _ =>
                   <>
                     <LoadingCensorBar width=20 height=20 radius=20 />
                     <HSpacing size=Spacing.sm />
                     <LoadingCensorBar width=60 height=15 />
                   </>
                 }}
              </div>
            </div>
          </Col.Grid>
        </Row>
        {switch (txSub) {
         | Data({success, errMsg}) when !success =>
           <Row> <Col.Grid> <TxError.Full msg=errMsg /> </Col.Grid> </Row>
         | _ => React.null
         }}
        <Row marginBottom=24>
          <Col.Grid>
            <div className=Styles.infoContainer>
              <Heading
                value="Information"
                size=Heading.H4
                style=Styles.infoHeader
                marginBottom=24
              />
              <Row>
                <Col.Grid col=Col.Six mb=24 mbSm=24>
                  <Heading value="Block" size=Heading.H5 marginBottom=8 />
                  {switch (txSub) {
                   | Data({blockHeight}) =>
                     <TypeID.Block id=blockHeight position=TypeID.Subtitle />
                   | _ => <LoadingCensorBar width=75 height=15 />
                   }}
                </Col.Grid>
                <Col.Grid col=Col.Six mb=24 mbSm=24>
                  <Heading value="Sender" size=Heading.H5 marginBottom=8 />
                  {switch (txSub) {
                   | Data({sender}) =>
                     <AddressRender address=sender position=AddressRender.Subtitle />
                   | _ => <LoadingCensorBar width=280 height=15 />
                   }}
                </Col.Grid>
                <Col.Grid col=Col.Six mb=24 mbSm=24>
                  <Heading value="Timestamp" size=Heading.H5 marginBottom=8 />
                  {switch (txSub) {
                   | Data({timestamp}) =>
                     <div className={CssHelper.flexBox()}>
                       <Text
                         value={
                           timestamp
                           |> MomentRe.Moment.format(Config.timestampDisplayFormat)
                           |> String.uppercase_ascii
                         }
                         size=Text.Lg
                       />
                       <HSpacing size=Spacing.sm />
                       <TimeAgos time=timestamp prefix="(" suffix=")" color=Colors.gray7 />
                     </div>
                   | _ => <LoadingCensorBar width=280 height=15 />
                   }}
                </Col.Grid>
                <Col.Grid>
                  <Heading value="Memo" size=Heading.H5 marginBottom=8 />
                  {switch (txSub) {
                   | Data({memo}) =>
                     <p>
                       <Text
                         value=memo
                         weight=Text.Regular
                         size=Text.Lg
                         color=Colors.gray7
                         block=true
                       />
                     </p>
                   | _ => <LoadingCensorBar width=280 height=15 />
                   }}
                </Col.Grid>
              </Row>
              <hr className=Styles.separatorLine />
              <Row>
                <Col.Grid col=Col.Three colSm=Col.Six mbSm=24>
                  <Heading value="Gas Used" size=Heading.H5 marginBottom=8 />
                  {switch (txSub) {
                   | Data({gasUsed}) => <Text value={gasUsed |> Format.iPretty} size=Text.Lg />
                   | _ => <LoadingCensorBar width=75 height=15 />
                   }}
                </Col.Grid>
                <Col.Grid col=Col.Three colSm=Col.Six mbSm=24>
                  <Heading value="Gas Limit" size=Heading.H5 marginBottom=8 />
                  {switch (txSub) {
                   | Data({gasLimit}) => <Text value={gasLimit |> Format.iPretty} size=Text.Lg />
                   | _ => <LoadingCensorBar width=75 height=15 />
                   }}
                </Col.Grid>
                <Col.Grid col=Col.Three colSm=Col.Six>
                  <Heading value="Gas Price (UBAND)" size=Heading.H5 marginBottom=8 />
                  {switch (txSub) {
                   | Data({gasFee, gasLimit}) =>
                     <Text
                       value={
                         (gasFee |> Coin.getBandAmountFromCoins)
                         /. (gasLimit |> float_of_int)
                         *. 1e6
                         |> Format.fPretty
                       }
                       size=Text.Lg
                     />
                   | _ => <LoadingCensorBar width=75 height=15 />
                   }}
                </Col.Grid>
                <Col.Grid col=Col.Three colSm=Col.Six>
                  <Heading value="Fee (BAND)" size=Heading.H5 marginBottom=8 />
                  {switch (txSub) {
                   | Data({gasFee}) =>
                     <Text
                       value={gasFee |> Coin.getBandAmountFromCoins |> Format.fPretty}
                       size=Text.Lg
                     />
                   | _ => <LoadingCensorBar width=75 height=15 />
                   }}
                </Col.Grid>
              </Row>
            </div>
          </Col.Grid>
        </Row>
        <Row marginBottom=24>
          <Col.Grid>
            {switch (txSub) {
             | Data({messages}) =>
               let msgCount = messages |> Belt.List.length;
               <div className={CssHelper.flexBox()}>
                 <Text value={msgCount |> string_of_int} size=Text.Lg />
                 <HSpacing size=Spacing.md />
                 <Text value={msgCount > 1 ? "messages" : "message"} size=Text.Lg />
               </div>;

             | _ => <LoadingCensorBar width=100 height=20 />
             }}
          </Col.Grid>
        </Row>
        {switch (txSub) {
         | Data({messages}) => <TxIndexPageTable messages />
         | _ => <TxIndexPageTable.Loading />
         }}
      </div>
    </Section>
  | _ => <TxNotFound />
  };
};
