module Styles = {
  open Css;

  let pageContainer = style([paddingTop(`px(40))]);

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
};

[@react.component]
let make = (~reqID, ~hashtag: Route.request_tab_t) => {
  let requestOpt = RequestHook.get(reqID);

  <div className=Styles.pageContainer>
    <Row justify=Row.Between>
      <Col>
        <div className=Styles.vFlex>
          <img src=Images.oracleScriptLogo className=Styles.logo />
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
               prefix="Last updated "
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
           <Col size=1.>
             <InfoHL header="OWNER" info={InfoHL.Address(request.requester, 430)} />
           </Col>
           <Col size=0.95>
             <InfoHL
               info={
                 InfoHL.OracleScript(
                   ID.OracleScript.ID(request.oracleScriptID),
                   request.oracleScriptName,
                 )
               }
               header="RELATED DATA SOURCES"
             />
           </Col>
         </Row>
         <VSpacing size=Spacing.xl />
       </>
     | None => React.null
     }}
  </div>;
};
