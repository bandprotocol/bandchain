module Styles = {
  open Css;

  let tableWrapper = style([padding2(~v=`px(20), ~h=`px(15))]);

  let withWidth = w => style([width(`px(w))]);

  let icon = style([width(`px(80)), height(`px(80))]);
  let iconWrapper =
    style([
      width(`percent(100.)),
      display(`flex),
      flexDirection(`column),
      alignItems(`center),
    ]);
};

type request_t = {
  id: int,
  requester: Address.t,
  age: MomentRe.Moment.t,
  blockHeight: int,
  txHash: Hash.t,
};

[@react.component]
let make = () => {
  let requests: list(request_t) = [
    {
      id: 6,
      requester: "e38475F47166d30A6e4E2E2C37e4B75E88Aa8b5B" |> Address.fromHex,
      age: MomentRe.momentNow(),
      blockHeight: 234554,
      txHash: Hash.fromHex("e7f3388a05a804fa99470aa90a18c60abb6b41b8f766e2096db5b1ad89154538"),
    },
    {
      id: 23,
      requester: "e38475F47166d30A6e4E2E2C37e4B75E88Aa8b5B" |> Address.fromHex,
      age:
        MomentRe.momentNow() |> MomentRe.Moment.subtract(~duration=MomentRe.duration(2., `hours)),
      blockHeight: 64563,
      txHash: Hash.fromHex("90cf054923b80b6cf18fceb5a930aea45a9726c450620c48a5626d79740542dd"),
    },
    {
      id: 162,
      requester: "e38475F47166d30A6e4E2E2C37e4B75E88Aa8b5B" |> Address.fromHex,
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
      <Text value="Requests" />
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
                     weight=Text.Semibold
                     color=Colors.gray6
                   />
                 </div>
               </Col>
               <Col size=2.64>
                 <Text
                   block=true
                   value="REQUESTER"
                   size=Text.Sm
                   weight=Text.Semibold
                   color=Colors.gray6
                 />
               </Col>
               <Col size=1.61>
                 <Text
                   block=true
                   value="AGE"
                   size=Text.Sm
                   weight=Text.Semibold
                   color=Colors.gray6
                 />
               </Col>
               <Col size=1.26>
                 <Text
                   block=true
                   value="BLOCK"
                   size=Text.Sm
                   weight=Text.Semibold
                   color=Colors.gray6
                 />
               </Col>
               <Col size=2.8>
                 <Text
                   block=true
                   value="TX HASH"
                   size=Text.Sm
                   weight=Text.Semibold
                   color=Colors.gray6
                 />
               </Col>
               <Col> <HSpacing size=Spacing.lg /> </Col>
             </Row>
           </THead>
           {requests
            ->Belt.List.map(({id, requester, age, blockHeight, txHash}) => {
                <TBody key={txHash |> Hash.toHex(~upper=true)}>
                  <Row>
                    <Col> <HSpacing size=Spacing.lg /> </Col>
                    <Col size=1.> <TypeID.Request id={ID.Request.ID(id)} /> </Col>
                    <Col size=2.64>
                      <div className={Styles.withWidth(220)}>
                        <AddressRender address=requester />
                      </div>
                    </Col>
                    <Col size=1.61> <TimeAgos time=age size=Text.Md weight=Text.Medium /> </Col>
                    <Col size=1.26> <TypeID.Block id={ID.Block.ID(blockHeight)} /> </Col>
                    <Col size=2.8>
                      <div className={Styles.withWidth(230)}>
                        <Text
                          block=true
                          value={txHash |> Hash.toHex(~upper=true)}
                          weight=Text.Medium
                          code=true
                          color=Colors.gray7
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
           <Text block=true value="NO REQUEST" weight=Text.Regular color=Colors.blue4 />
           <VSpacing size={`px(15)} />
         </div>}
  </div>;
};
