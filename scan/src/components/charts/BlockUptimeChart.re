module Styles = {
  open Css;
  let chartContainer =
    style([
      width(`percent(100.)),
      maxWidth(`px(400)),
      margin2(~v=`zero, ~h=`auto),
      Media.smallMobile([maxWidth(`px(260))]),
    ]);
  let statusLabel = style([height(`px(8)), width(`px(8))]);

  let blockContainer =
    style([
      flexGrow(0.),
      flexShrink(0.),
      flexBasis(`calc((`sub, `percent(5.), `px(2)))),
      margin(`px(1)),
      height(`px(16)),
      display(`block),
      Media.smallMobile([height(`px(10))]),
    ]);
  let blockBase = style([width(`percent(100.)), height(`percent(100.))]);
  let status =
    fun
    | ValidatorSub.Missed => style([backgroundColor(Colors.blue10)])
    | Proposed => style([backgroundColor(Colors.blue9)])
    | Signed => style([backgroundColor(Colors.bandBlue)]);

  let labelBox =
    style([
      margin2(~v=`zero, ~h=`px(-12)),
      selector(
        "> div",
        [
          flexGrow(0.),
          flexShrink(0.),
          flexBasis(`calc((`sub, `percent(33.33), `px(24)))),
          margin2(~v=`zero, ~h=`px(12)),
        ],
      ),
      Media.smallMobile([
        margin2(~v=`zero, ~h=`px(-5)),
        selector(
          "> div",
          [
            flexBasis(`calc((`sub, `percent(33.33), `px(10)))),
            margin2(~v=`zero, ~h=`px(5)),
          ],
        ),
      ]),
    ]);
};

module UptimeBlock = {
  [@react.component]
  let make = (~status, ~height) => {
    let blockRoute = height |> ID.Block.getRoute;
    <CTooltip
      width=90
      tooltipPlacement=CTooltip.Top
      tooltipPlacementSm=CTooltip.BottomLeft
      mobile=false
      align=`center
      pd=10
      tooltipText={height |> ID.Block.toString}
      styles=Styles.blockContainer>
      <Link route=blockRoute className=Styles.blockContainer>
        <div className={Css.merge([Styles.blockBase, Styles.status(status)])} />
      </Link>
    </CTooltip>;
  };
};

[@react.component]
let make = (~consensusAddress) => {
  let getUptimeSub = ValidatorSub.getBlockUptimeByValidator(consensusAddress);
  <>
    <Row marginBottom=24>
      <Col.Grid>
        <div className={Css.merge([CssHelper.flexBox(), Styles.chartContainer])}>
          {switch (getUptimeSub) {
           | Data({validatorVotes}) =>
             validatorVotes
             ->Belt.Array.map(({blockHeight, status}) =>
                 <UptimeBlock key={blockHeight |> ID.Block.toString} status height=blockHeight />
               )
             ->React.array
           | _ => <LoadingCensorBar fullWidth=true height=90 />
           }}
        </div>
      </Col.Grid>
    </Row>
    <Row>
      <Col.Grid>
        <div className={Css.merge([CssHelper.flexBox(), Styles.labelBox])}>
          <div className={CssHelper.flexBox(~justify=`spaceBetween, ())}>
            <div className={CssHelper.flexBox()}>
              <div
                className={Css.merge([Styles.status(ValidatorSub.Proposed), Styles.statusLabel])}
              />
              <HSpacing size=Spacing.sm />
              <Text block=true value="Proposed" weight=Text.Semibold />
            </div>
            {switch (getUptimeSub) {
             | Data({proposedCount}) => <Text block=true value={proposedCount |> string_of_int} />
             | _ => <LoadingCensorBar width=20 height=14 />
             }}
          </div>
          <div className={CssHelper.flexBox(~justify=`spaceBetween, ())}>
            <div className={CssHelper.flexBox()}>
              <div
                className={Css.merge([Styles.status(ValidatorSub.Signed), Styles.statusLabel])}
              />
              <HSpacing size=Spacing.sm />
              <Text block=true value="Signed" weight=Text.Semibold />
            </div>
            {switch (getUptimeSub) {
             | Data({signedCount}) => <Text block=true value={signedCount |> string_of_int} />
             | _ => <LoadingCensorBar width=20 height=14 />
             }}
          </div>
          <div className={CssHelper.flexBox(~justify=`spaceBetween, ())}>
            <div className={CssHelper.flexBox()}>
              <div
                className={Css.merge([Styles.status(ValidatorSub.Missed), Styles.statusLabel])}
              />
              <HSpacing size=Spacing.sm />
              <Text block=true value="Missed" weight=Text.Semibold />
            </div>
            {switch (getUptimeSub) {
             | Data({missedCount}) => <Text block=true value={missedCount |> string_of_int} />
             | _ => <LoadingCensorBar width=20 height=14 />
             }}
          </div>
        </div>
      </Col.Grid>
    </Row>
  </>;
};
