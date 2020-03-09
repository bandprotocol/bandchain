module Styles = {
  open Css;

  let tableWrapper = style([padding2(~v=`px(20), ~h=`px(15))]);

  let fixWidth = style([width(`px(230))]);

  let icon = style([width(`px(80)), height(`px(80))]);
  let iconWrapper =
    style([
      width(`percent(100.)),
      display(`flex),
      flexDirection(`column),
      alignItems(`center),
    ]);
};

type oracle_script_t = {
  id: int,
  description: string,
};

type request_t = {
  id: int,
  oracleScript: oracle_script_t,
  age: MomentRe.Moment.t,
  blockHeight: int,
  txHash: Hash.t,
};

[@react.component]
let make = () => {
  let requests: list(request_t) = [
    {
      id: 6,
      oracleScript: {
        id: 895,
        description: "Mean Bitcoin Price",
      },
      age: MomentRe.momentNow(),
      blockHeight: 234554,
      txHash: Hash.fromHex("e7f3388a05a804fa99470aa90a18c60abb6b41b8f766e2096db5b1ad89154538"),
    },
    {
      id: 23,
      oracleScript: {
        id: 32,
        description: "Median Stellar Price",
      },
      age:
        MomentRe.momentNow() |> MomentRe.Moment.subtract(~duration=MomentRe.duration(2., `hours)),
      blockHeight: 64563,
      txHash: Hash.fromHex("90cf054923b80b6cf18fceb5a930aea45a9726c450620c48a5626d79740542dd"),
    },
    {
      id: 162,
      oracleScript: {
        id: 786,
        description: "Advance Algo for Crypto Price",
      },
      age:
        MomentRe.momentNow() |> MomentRe.Moment.subtract(~duration=MomentRe.duration(1., `days)),
      blockHeight: 3425,
      txHash: Hash.fromHex("d12f97901f466f6c2e9680798a7460413c538776cdd85372be601d7603f8de17"),
    },
  ];

  let numRequest = requests |> Belt_List.size;

  <div className=Styles.tableWrapper>
    <Row>
      <HSpacing size={`px(25)} />
      <Text value={numRequest |> string_of_int} weight=Text.Bold />
      <HSpacing size={`px(5)} />
      <Text value="Revisions" />
    </Row>
    <VSpacing size=Spacing.lg />
    {numRequest > 0
       ? <>
           <THead>
             <Row>
               <Col> <HSpacing size=Spacing.lg /> </Col>
               <Col size=1.>
                 <div className=TElement.Styles.hashContainer>
                   <Text
                     block=true
                     value="REQUEST"
                     size=Text.Sm
                     weight=Text.Bold
                     color=Colors.mediumLightGray
                   />
                 </div>
               </Col>
               <Col size=2.8>
                 <Text
                   block=true
                   value="ORACLE SCRIPT"
                   size=Text.Sm
                   weight=Text.Bold
                   color=Colors.mediumLightGray
                 />
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
               <Col size=2.7>
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
           {requests
            ->Belt.List.map(({id, oracleScript, age, blockHeight, txHash}) => {
                <TBody key={txHash |> Hash.toHex(~upper=true)}>
                  <Row>
                    <Col> <HSpacing size=Spacing.lg /> </Col>
                    <Col size=1.>
                      <Text
                        block=true
                        value={"#R" ++ (id |> string_of_int)}
                        weight=Text.Bold
                        code=true
                        color=Colors.brightOrange
                      />
                    </Col>
                    <Col size=2.8>
                      <Row>
                        <Text
                          block=true
                          value={"#O" ++ (oracleScript.id |> string_of_int)}
                          weight=Text.Bold
                          code=true
                          color=Colors.brightPink
                        />
                        <HSpacing size={`px(5)} />
                        <Text
                          block=true
                          value={oracleScript.description}
                          weight=Text.Medium
                          color=Colors.mediumGray
                        />
                      </Row>
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
                    <Col size=2.7>
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
           <img src=Images.noRequestIcon className=Styles.icon />
           <VSpacing size={`px(40)} />
           <Text block=true value="NO REQUEST" weight=Text.Regular color=Colors.brightLightBlue />
           <VSpacing size={`px(15)} />
         </div>}
  </div>;
};
