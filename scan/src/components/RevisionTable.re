module Styles = {
  open Css;

  let tableWrapper = style([padding2(~v=`px(20), ~h=`px(15))]);

  let fixWidth = style([width(`px(300))]);

  let icon = style([width(`px(80)), height(`px(80))]);
  let iconWrapper =
    style([
      width(`percent(100.)),
      display(`flex),
      flexDirection(`column),
      alignItems(`center),
    ]);
};

type revision_t = {
  name: string,
  age: MomentRe.Moment.t,
  blockHeight: int,
  txHash: Hash.t,
};

[@react.component]
let make = () => {
  let revisions: list(revision_t) = [
    {
      name: "Binance OpenAPI v1",
      age: MomentRe.momentNow(),
      blockHeight: 234554,
      txHash: Hash.fromHex("e7f3388a05a804fa99470aa90a18c60abb6b41b8f766e2096db5b1ad89154538"),
    },
    {
      name: "CoinMarketCap With Timestamp",
      age:
        MomentRe.momentNow() |> MomentRe.Moment.subtract(~duration=MomentRe.duration(2., `hours)),
      blockHeight: 64563,
      txHash: Hash.fromHex("90cf054923b80b6cf18fceb5a930aea45a9726c450620c48a5626d79740542dd"),
    },
    {
      name: "Median Crypto Price",
      age:
        MomentRe.momentNow() |> MomentRe.Moment.subtract(~duration=MomentRe.duration(1., `days)),
      blockHeight: 13425,
      txHash: Hash.fromHex("d12f97901f466f6c2e9680798a7460413c538776cdd85372be601d7603f8de17"),
    },
    {
      name: "Advance Premium Crypto Price",
      age:
        MomentRe.momentNow() |> MomentRe.Moment.subtract(~duration=MomentRe.duration(3., `days)),
      blockHeight: 2542,
      txHash: Hash.fromHex("3f75f78492711fbe2a3d97fe06304616bc994b6c297571fc883fd869a91478f3"),
    },
  ];

  let numRevision = revisions |> Belt_List.size;
  <div className=Styles.tableWrapper>
    <Row>
      <HSpacing size={`px(25)} />
      <Text value={numRevision |> string_of_int} weight=Text.Bold />
      <HSpacing size={`px(5)} />
      <Text value="Revisions" />
    </Row>
    <VSpacing size=Spacing.lg />
    {numRevision > 0
       ? <>
           <THead>
             <Row>
               <Col> <HSpacing size=Spacing.lg /> </Col>
               <Col size=3.>
                 <div className=TElement.Styles.hashContainer>
                   <Text
                     block=true
                     value="NAME"
                     size=Text.Sm
                     weight=Text.Bold
                     color=Colors.mediumLightGray
                   />
                 </div>
               </Col>
               <Col size=2.>
                 <Text
                   block=true
                   value="AGE"
                   size=Text.Sm
                   weight=Text.Bold
                   color=Colors.mediumLightGray
                 />
               </Col>
               <Col size=1.5>
                 <Text
                   block=true
                   value="BLOCK"
                   size=Text.Sm
                   weight=Text.Bold
                   color=Colors.mediumLightGray
                 />
               </Col>
               <Col size=3.5>
                 <Text
                   block=true
                   value="TX HASH"
                   size=Text.Sm
                   weight=Text.Bold
                   color=Colors.mediumLightGray
                 />
               </Col>
               <Col> <HSpacing size=Spacing.lg /> </Col>
             </Row>
           </THead>
           {revisions
            ->Belt.List.map(({name, age, blockHeight, txHash}) => {
                <TBody key={txHash |> Hash.toHex(~upper=true)}>
                  <Row>
                    <Col> <HSpacing size=Spacing.lg /> </Col>
                    <Col size=3.>
                      <Text block=true value=name weight=Text.Medium color=Colors.mediumGray />
                    </Col>
                    <Col size=2.> <TimeAgos time=age size=Text.Md weight=Text.Medium /> </Col>
                    <Col size=1.5>
                      <Text
                        block=true
                        value={"#B" ++ (blockHeight |> string_of_int)}
                        weight=Text.Bold
                        code=true
                        color=Colors.brightBlue
                      />
                    </Col>
                    <Col size=3.5>
                      <div className=Styles.fixWidth>
                        <Text
                          block=true
                          value={txHash |> Hash.toHex(~upper=true)}
                          weight=Text.Medium
                          code=true
                          color=Colors.mediumGray
                          ellipsis=true
                          nowrap=true
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
           <img src=Images.noRevisionIcon className=Styles.icon />
           <VSpacing size={`px(40)} />
           <Text block=true value="NO REVISION" weight=Text.Regular color=Colors.brightLightBlue />
           <VSpacing size={`px(15)} />
         </div>}
  </div>;
};
