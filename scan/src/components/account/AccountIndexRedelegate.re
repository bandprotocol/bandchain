module Styles = {
  open Css;

  let tableLowerContainer = style([padding(`px(10))]);

  let hFlex = style([display(`flex)]);

  let alignRight = style([display(`flex), justifyContent(`flexEnd)]);
  let alignLeft = style([display(`flex), justifyContent(`flexStart)]);
};

let renderBody =
    ({srcValidator, dstValidator, completionTime, amount}: RedelegateSub.redelegate_list_t) => {
  <TBody
    key={
      (srcValidator.operatorAddress |> Address.toBech32)
      ++ (dstValidator.operatorAddress |> Address.toBech32)
      ++ (completionTime |> MomentRe.Moment.toISOString)
      ++ (amount |> Coin.getBandAmountFromCoin |> Js.Float.toString)
    }
    minHeight=50>
    <Row>
      <Col> <HSpacing size=Spacing.lg /> </Col>
      <Col size=1.>
        <ValidatorMonikerLink
          validatorAddress={srcValidator.operatorAddress}
          moniker={srcValidator.moniker}
          identity={srcValidator.identity}
          width={`px(200)}
        />
      </Col>
      <Col size=1.>
        <div className=Styles.alignLeft>
          <ValidatorMonikerLink
            validatorAddress={dstValidator.operatorAddress}
            moniker={dstValidator.moniker}
            identity={dstValidator.identity}
            width={`px(200)}
          />
        </div>
      </Col>
      <Col size=0.6>
        <div className=Styles.alignRight>
          <Text value={amount |> Coin.getBandAmountFromCoin |> Format.fPretty} code=true />
        </div>
      </Col>
      <Col size=1.>
        <div className=Styles.alignRight>
          <Text
            value={
              completionTime
              |> MomentRe.Moment.format(Config.timestampDisplayFormat)
              |> String.uppercase_ascii
            }
            code=true
          />
        </div>
      </Col>
      <Col> <HSpacing size=Spacing.lg /> </Col>
    </Row>
  </TBody>;
};

let renderBodyMobile =
    ({amount, completionTime, dstValidator, srcValidator}: RedelegateSub.redelegate_list_t) => {
  let key_ =
    (srcValidator.operatorAddress |> Address.toBech32)
    ++ (dstValidator.operatorAddress |> Address.toBech32)
    ++ (completionTime |> MomentRe.Moment.toISOString)
    ++ (amount |> Coin.getBandAmountFromCoin |> Js.Float.toString);

  <MobileCard
    values=InfoMobileCard.[
      (
        "SOURCE\nVALIDATOR",
        Validator(srcValidator.operatorAddress, srcValidator.moniker, srcValidator.identity),
      ),
      (
        "DESTINATION\nVALIDATOR",
        Validator(dstValidator.operatorAddress, dstValidator.moniker, dstValidator.identity),
      ),
      ("AMOUNT\n(BAND)", Coin({value: [amount], hasDenom: false})),
      ("REDELEGATE\nCOMPLETE AT", Timestamp(completionTime)),
    ]
    key=key_
    idx=key_
  />;
};

[@react.component]
let make = (~address) =>
  {
    let isMobile = Media.isMobile();
    let currentTime =
      React.useContext(TimeContext.context) |> MomentRe.Moment.format(Config.timestampUseFormat);

    let (page, setPage) = React.useState(_ => 1);
    let pageSize = 10;

    let redelegateCountSub = RedelegateSub.getRedelegateCountByDelegator(address, currentTime);
    let redelegateListSub =
      RedelegateSub.getRedelegationByDelegator(address, currentTime, ~pageSize, ~page, ());

    let%Sub redelegateCount = redelegateCountSub;
    let%Sub redelegateList = redelegateListSub;

    let pageCount = Page.getPageCount(redelegateCount, pageSize);
    <div className=Styles.tableLowerContainer>
      <VSpacing size=Spacing.md />
      <div className=Styles.hFlex>
        <HSpacing size=Spacing.lg />
        <Text value={redelegateCount |> string_of_int} weight=Text.Semibold />
        <HSpacing size=Spacing.xs />
        <Text value="Redelegate Entries" />
      </div>
      <VSpacing size=Spacing.lg />
      <>
        {isMobile
           ? React.null
           : <THead>
               <Row>
                 <Col> <HSpacing size=Spacing.lg /> </Col>
                 <Col size=1.>
                   <Text
                     block=true
                     value="SOURCE VALIDATOR"
                     size=Text.Sm
                     weight=Text.Bold
                     spacing={Text.Em(0.05)}
                     color=Colors.gray6
                   />
                 </Col>
                 <Col size=1.>
                   <div className=Styles.alignLeft>
                     <Text
                       block=true
                       value="DESTINATION VALIDATOR"
                       size=Text.Sm
                       weight=Text.Bold
                       spacing={Text.Em(0.05)}
                       color=Colors.gray6
                     />
                   </div>
                 </Col>
                 <Col size=0.6>
                   <div className=Styles.alignRight>
                     <Text
                       block=true
                       value="AMOUNT (BAND)"
                       size=Text.Sm
                       spacing={Text.Em(0.05)}
                       weight=Text.Bold
                       color=Colors.gray6
                     />
                   </div>
                 </Col>
                 <Col size=1.>
                   <div className=Styles.alignRight>
                     <Text
                       block=true
                       value="REDELEGATE COMPLETE AT"
                       size=Text.Sm
                       spacing={Text.Em(0.05)}
                       weight=Text.Bold
                       color=Colors.gray6
                     />
                   </div>
                 </Col>
                 <Col> <HSpacing size=Spacing.lg /> </Col>
               </Row>
             </THead>}
        {redelegateList
         ->Belt.Array.map(redelegateEntry =>
             isMobile ? renderBodyMobile(redelegateEntry) : renderBody(redelegateEntry)
           )
         ->React.array}
        <VSpacing size=Spacing.lg />
        <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />
      </>
    </div>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
