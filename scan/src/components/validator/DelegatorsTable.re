module Styles = {
  open Css;

  let vFlex = align => style([display(`flex), flexDirection(`row), alignItems(align)]);

  let tableWrapper = style([padding2(~v=`px(20), ~h=`px(15))]);

  let icon = style([width(`px(80)), height(`px(80))]);
  let iconWrapper =
    style([
      width(`percent(100.)),
      display(`flex),
      flexDirection(`column),
      alignItems(`center),
    ]);

  let withWidth = w => style([width(`px(w))]);

  let fillLeft = style([marginLeft(`auto)]);
};

[@react.component]
let make = () => {
  // TODO: Mock to use
  let (page, setPage) = React.useState(_ => 1);

  let pageCount = 1;

  let delegators = [
    ("band1j9vk75jjty02elhwqqjehaspfslaem8p0utr4q", 12.0, 123123.0),
    ("band1j9vk75jjty02elhwqqjehaspfslaem8p0utr4q", 88.0, 883312.0),
  ];

  <div className=Styles.tableWrapper>
    <Row>
      <HSpacing size={`px(25)} />
      <Text value={delegators |> Belt_List.length |> string_of_int} weight=Text.Bold />
      <HSpacing size={`px(5)} />
      <Text value="Delegators" />
    </Row>
    <VSpacing size=Spacing.lg />
    {delegators->Belt_List.length > 0
       ? <>
           <THead>
             <Row>
               <Col> <HSpacing size=Spacing.lg /> </Col>
               <Col size=1.4>
                 <Text
                   block=true
                   value="DELEGATOR"
                   size=Text.Sm
                   weight=Text.Semibold
                   color=Colors.gray6
                   spacing={Text.Em(0.05)}
                 />
               </Col>
               <Col size=1.45>
                 <div className={Styles.vFlex(`flexEnd)}>
                   <div className=Styles.fillLeft />
                   <Text
                     block=true
                     value="SHARE (%)"
                     size=Text.Sm
                     weight=Text.Semibold
                     color=Colors.gray6
                     spacing={Text.Em(0.05)}
                   />
                 </div>
               </Col>
               <Col size=1.45>
                 <div className={Styles.vFlex(`flexEnd)}>
                   <div className=Styles.fillLeft />
                   <Text
                     block=true
                     value="AMOUNT (BAND)"
                     size=Text.Sm
                     weight=Text.Semibold
                     color=Colors.gray6
                     spacing={Text.Em(0.05)}
                   />
                 </div>
               </Col>
               <Col> <HSpacing size=Spacing.lg /> </Col>
             </Row>
           </THead>
           {delegators
            ->Belt.List.map(((address, share, amount)) => {
                <TBody>
                  <Row>
                    <Col> <HSpacing size=Spacing.lg /> </Col>
                    <Col size=1.4> <AddressRender address={address |> Address.fromBech32} /> </Col>
                    <Col size=1.30>
                      <div className={Styles.vFlex(`flexEnd)}>
                        <div className=Styles.fillLeft />
                        <Text
                          block=true
                          value={share |> Format.fPretty}
                          size=Text.Md
                          weight=Text.Regular
                          color=Colors.gray7
                          spacing={Text.Em(0.05)}
                          code=true
                        />
                      </div>
                    </Col>
                    <Col size=1.45>
                      <div className={Styles.vFlex(`flexEnd)}>
                        <div className=Styles.fillLeft />
                        <Text
                          block=true
                          value={amount |> Format.fPretty}
                          size=Text.Md
                          weight=Text.Regular
                          color=Colors.gray7
                          spacing={Text.Em(0.05)}
                          code=true
                        />
                      </div>
                    </Col>
                    <Col> <HSpacing size=Spacing.lg /> </Col>
                  </Row>
                </TBody>
              })
            ->Array.of_list
            ->React.array}
         </>
       : <div className=Styles.iconWrapper>
           <VSpacing size={`px(30)} />
           <img src=Images.noRequestIcon className=Styles.icon />
           <VSpacing size={`px(40)} />
           <Text block=true value="NO DELEGATORS" weight=Text.Regular color=Colors.blue4 />
           <VSpacing size={`px(15)} />
         </div>}
    <VSpacing size=Spacing.xl />
    <VSpacing size=Spacing.sm />
    <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />
  </div>;
};
