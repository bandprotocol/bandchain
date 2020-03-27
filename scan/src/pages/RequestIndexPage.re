module Styles = {
  open Css;

  let pageContainer = style([paddingTop(`px(20))]);

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let logo = style([width(`px(50)), marginRight(`px(10))]);
  let headerContainer = style([lineHeight(`px(25))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.gray7),
    ]);

  let lowerPannel =
    style([
      width(`percent(100.)),
      height(`px(540)),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      backgroundColor(Colors.white),
      boxShadows([
        Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(4), Css.rgba(0, 0, 0, 0.1)),
        Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(12), Css.rgba(0, 0, 0, 0.03)),
      ]),
      borderRadius(`px(10)),
    ]);
};

[@react.component]
let make = (~reqID, ~hashtag: Route.request_tab_t) => {
  let requestOpt = RequestHook.get(reqID);
  let requestValidators =
    switch (React.useContext(GlobalContext.context), requestOpt) {
    | (Some(info), Some(request)) =>
      info.validators
      ->Belt_List.keep(validator =>
          request.requestedValidators
          ->Belt_List.has(validator.operatorAddress, (a, b) => a->Address.isEqual(b))
        )
    | _ => []
    };

  <div className=Styles.pageContainer>
    <Row justify=Row.Between>
      <Col>
        <div className=Styles.vFlex>
          <img src=Images.requestLogo className=Styles.logo />
          <Text
            value="DATA REQUEST"
            weight=Text.Medium
            size=Text.Md
            spacing={Text.Em(0.06)}
            height={Text.Px(15)}
            nowrap=true
            color=Colors.gray7
            block=true
          />
          <div className=Styles.seperatedLine />
          {switch (requestOpt) {
           | Some(request) =>
             <TimeAgos
               time={request.requestedAtTime}
               size=Text.Md
               weight=Text.Thin
               spacing={Text.Em(0.06)}
               height={Text.Px(18)}
               upper=true
             />

           | None =>
             <Text
               value="???"
               size=Text.Md
               weight=Text.Thin
               spacing={Text.Em(0.06)}
               height={Text.Px(18)}
             />
           }}
        </div>
      </Col>
    </Row>
    {switch (requestOpt) {
     | Some(request) =>
       <>
         <VSpacing size=Spacing.xl />
         <div className=Styles.vFlex>
           <TypeID.Request id={ID.Request.ID(request.id)} position=TypeID.Title />
         </div>
         <VSpacing size=Spacing.xl />
         <Row>
           <Col size=2.8>
             <InfoHL
               info={
                 InfoHL.OracleScript(
                   ID.OracleScript.ID(request.oracleScriptID),
                   request.oracleScriptName,
                 )
               }
               header="ORACLE SCRIPT"
             />
           </Col>
           <Col size=3.2>
             <InfoHL header="SENDER" info={InfoHL.Address(request.requester, 280)} />
           </Col>
           <Col size=4.0>
             <InfoHL header="TX HASH" info={InfoHL.TxHash(request.txHash, 385)} />
           </Col>
         </Row>
         <VSpacing size=Spacing.xl />
         <Row>
           <Col>
             <InfoHL info={InfoHL.Validators(requestValidators)} header="REQUEST TO VALIDATORS" />
           </Col>
         </Row>
       </>
     | None => React.null
     }}
    <VSpacing size=Spacing.xl />
    <div className=Styles.lowerPannel> {"TODO" |> React.string} </div>
  </div>;
};
