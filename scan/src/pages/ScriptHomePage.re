module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let pageContainer = style([paddingTop(`px(50)), minHeight(`px(500))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.grayHeader),
    ]);

  let textContainer = style([paddingLeft(Spacing.lg), display(`flex)]);

  let proposerBox = style([maxWidth(`px(270)), display(`flex), flexDirection(`column)]);
};

[@react.component]
let make = () => {
  let scriptsOpt = ScriptHook.getScriptList(~limit=100000, ~pollInterval=3000, ());
  let scripts = scriptsOpt->Belt.Option.getWithDefault([]);
  let totalScript = scripts->Belt.List.length->string_of_int;

  <div className=Styles.pageContainer>
    <Row>
      <Col>
        <div className=Styles.vFlex>
          <Text
            value="ALL DATA ORACLE SCRIPTS"
            weight=Text.Bold
            size=Text.Xl
            nowrap=true
            color=Colors.grayHeader
          />
          <div className=Styles.seperatedLine />
          <Text value={j|$totalScript in total|j} />
        </div>
      </Col>
    </Row>
    <VSpacing size=Spacing.xl />
    <THead>
      <Row>
        <Col> <HSpacing size=Spacing.xl /> </Col>
        <Col size=1.1>
          <div className=TElement.Styles.hashContainer>
            <Text block=true value="NAME" size=Text.Sm weight=Text.Bold color=Colors.grayText />
          </div>
        </Col>
        <Col size=1.1>
          <Text
            block=true
            value="SCRIPT HASH"
            size=Text.Sm
            weight=Text.Bold
            color=Colors.grayText
          />
        </Col>
        <Col size=0.65>
          <Text
            block=true
            value="CREATED AT"
            size=Text.Sm
            weight=Text.Bold
            color=Colors.grayText
          />
        </Col>
        <Col size=1.1>
          <Text block=true value="CREATOR" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=0.5>
          <div className=TElement.Styles.feeContainer>
            <Text
              block=true
              value="QUERY FEE"
              size=Text.Sm
              weight=Text.Bold
              color=Colors.grayText
            />
          </div>
        </Col>
      </Row>
    </THead>
    {scripts
     ->Belt.List.map(({info, txHash, createdAtTime}) => {
         <div
           onClick={_ =>
             Route.redirect(
               Route.ScriptIndexPage(info.codeHash, Route.ScriptTransactions),
             )
           }>
           <TBody key={txHash |> Hash.toHex}>
             <Row>
               <Col> <HSpacing size=Spacing.xl /> </Col>
               <Col size=1.1> <TElement elementType={info.name->TElement.Name} /> </Col>
               <Col size=1.1> <TElement elementType={info.codeHash->TElement.Hash} /> </Col>
               <Col size=0.65> <TElement elementType={createdAtTime->TElement.Timestamp} /> </Col>
               <Col size=1.1> <TElement elementType={info.creator->TElement.Address} /> </Col>
               <Col size=0.5> <TElement elementType={0.->TElement.Fee} /> </Col>
             </Row>
           </TBody>
         </div>
       })
     ->Array.of_list
     ->React.array}
    <VSpacing size=Spacing.lg />
    <VSpacing size=Spacing.xl />
    <VSpacing size=Spacing.xl />
  </div>;
};
