module Styles = {
  open Css;

  let noDataImage =
    style([
      width(`auto),
      height(`px(40)),
      marginBottom(`px(16)),
      Media.mobile([marginBottom(`px(8))]),
    ]);
  let container =
    style([
      boxShadow(
        Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, `num(0.08))),
      ),
    ]);
};

let renderBody = (reserveIndex, requestSub: ApolloHooks.Subscription.variant(RequestSub.t)) => {
  <TBody.Grid
    key={
      switch (requestSub) {
      | Data({id}) => id |> ID.Request.toString
      | _ => reserveIndex |> string_of_int
      }
    }
    paddingH={`px(24)}>
    <Row.Grid alignItems=Row.Center>
      <Col.Grid col=Col.Three>
        {switch (requestSub) {
         | Data({id}) => <TypeID.Request id />
         | _ => <LoadingCensorBar width=60 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Six>
        {switch (requestSub) {
         | Data({oracleScript: {oracleScriptID, name}}) =>
           <div className={CssHelper.flexBox(~wrap=`nowrap, ())}>
             <TypeID.OracleScript id=oracleScriptID />
             <HSpacing size=Spacing.sm />
             <Text value=name ellipsis=true />
           </div>
         | _ => <LoadingCensorBar width=150 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Three>
        <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
          {switch (requestSub) {
           | Data({resolveStatus, requestedValidators, reports}) =>
             let reportedCount = reports->Array.length;
             let requestedCount = requestedValidators->Array.length;

             <div className={CssHelper.flexBox()}>
               <Text value={j|$reportedCount of $requestedCount|j} />
               <HSpacing size=Spacing.sm />
               <RequestStatus resolveStatus />
             </div>;
           | _ => <LoadingCensorBar width=70 height=15 />
           }}
        </div>
      </Col.Grid>
    </Row.Grid>
  </TBody.Grid>;
};

let renderBodyMobile = (reserveIndex, requestSub: ApolloHooks.Subscription.variant(RequestSub.t)) => {
  switch (requestSub) {
  | Data({id, oracleScript: {oracleScriptID, name}, resolveStatus, requestedValidators, reports}) =>
    let reportedCount = reports->Array.length;
    let requestedCount = requestedValidators->Array.length;
    <MobileCard
      values=InfoMobileCard.[
        ("Request ID", RequestID(id)),
        ("Oracle Script", OracleScript(oracleScriptID, name)),
        ("Report Status", Text({j|$reportedCount of $requestedCount|j})),
      ]
      key={id |> ID.Request.toString}
      idx={id |> ID.Request.toString}
      requestStatus=resolveStatus
    />;
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("Request ID", Loading(60)),
        ("Oracle Script", Loading(150)),
        ("Report Status", Loading(60)),
      ]
      key={reserveIndex |> string_of_int}
      idx={reserveIndex |> string_of_int}
    />
  };
};

[@react.component]
let make = () => {
  let isMobile = Media.isMobile();
  let requestCount = 5;
  let requestsSub = RequestSub.getList(~page=1, ~pageSize=requestCount, ());

  <>
    <div
      className={CssHelper.flexBox(~justify=`spaceBetween, ~align=`flexEnd, ())}
      id="latestRequestsSectionHeader">
      <div>
        <Text
          value="Latest Requests"
          size=Text.Lg
          block=true
          color=Colors.gray7
          weight=Text.Medium
        />
        <VSpacing size={`px(4)} />
        {switch (requestsSub) {
         | ApolloHooks.Subscription.Data(requests) =>
           <Text
             value={
               requests
               ->Belt.Array.get(0)
               ->Belt.Option.mapWithDefault(0, ({id}) => id |> ID.Request.toInt)
               ->Format.iPretty
             }
             size=Text.Lg
             color=Colors.gray7
             weight=Text.Medium
           />
         | _ => <LoadingCensorBar width=90 height=18 />
         }}
      </div>
      <Link className={CssHelper.flexBox(~align=`center, ())} route=Route.RequestHomePage>
        <Text value="All Requests" color=Colors.bandBlue weight=Text.Medium />
        <HSpacing size=Spacing.md />
        <Icon name="fal fa-angle-right" color=Colors.bandBlue />
      </Link>
    </div>
    <VSpacing size={`px(16)} />
    {isMobile
       ? React.null
       : <THead.Grid height=30>
           <Row.Grid alignItems=Row.Center>
             <Col.Grid col=Col.Three>
               <Text
                 block=true
                 value="Request ID"
                 size=Text.Sm
                 weight=Text.Semibold
                 color=Colors.gray7
               />
             </Col.Grid>
             <Col.Grid col=Col.Six>
               <Text
                 block=true
                 value="Oracle Script"
                 size=Text.Sm
                 weight=Text.Semibold
                 color=Colors.gray7
               />
             </Col.Grid>
             <Col.Grid col=Col.Three>
               <Text
                 block=true
                 value="Report Status"
                 size=Text.Sm
                 weight=Text.Semibold
                 color=Colors.gray7
                 align=Text.Right
               />
             </Col.Grid>
           </Row.Grid>
         </THead.Grid>}
    {switch (requestsSub) {
     | Data(requests) =>
       requests->Belt.Array.length > 0
         ? requests
           ->Belt_Array.mapWithIndex((i, e) =>
               isMobile ? renderBodyMobile(i, Sub.resolve(e)) : renderBody(i, Sub.resolve(e))
             )
           ->React.array
         : <EmptyContainer height={`calc((`sub, `percent(100.), `px(86)))} boxShadow=true>
             <img src=Images.noSource className=Styles.noDataImage />
             <Heading
               size=Heading.H4
               value="No Request"
               align=Heading.Center
               weight=Heading.Regular
               color=Colors.bandBlue
             />
           </EmptyContainer>
     | _ =>
       Belt_Array.make(requestCount, ApolloHooks.Subscription.NoData)
       ->Belt_Array.mapWithIndex((i, noData) =>
           isMobile ? renderBodyMobile(i, noData) : renderBody(i, noData)
         )
       ->React.array
     }}
  </>;
};
