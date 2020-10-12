module SendMsg = {
  [@react.component]
  let make = (~send: TxSub.Msg.Send.t) => {
    <Row.Grid>
      <Col.Grid col=Col.Six mb=24>
        <Heading value="From" size=Heading.H5 marginBottom=8 />
        <AddressRender address={send.fromAddress} />
      </Col.Grid>
      <Col.Grid col=Col.Six mb=24>
        <Heading value="To" size=Heading.H5 marginBottom=8 />
        <AddressRender address={send.toAddress} />
      </Col.Grid>
      <Col.Grid col=Col.Six>
        <Heading value="Amount" size=Heading.H5 marginBottom=8 />
        <AmountRender coins={send.amount} pos=AmountRender.TxIndex />
      </Col.Grid>
    </Row.Grid>;
  };
};

module DelegateMsg = {
  [@react.component]
  let make = (~delegation: TxSub.Msg.Delegate.t) => {
    <Row.Grid>
      <Col.Grid col=Col.Six mb=24>
        <Heading value="Delegator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={delegation.delegatorAddress} />
      </Col.Grid>
      <Col.Grid col=Col.Six mb=24>
        <Heading value="Validator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={delegation.validatorAddress} accountType=`validator />
      </Col.Grid>
      <Col.Grid col=Col.Six>
        <Heading value="Amount" size=Heading.H5 marginBottom=8 />
        <AmountRender coins=[delegation.amount] pos=AmountRender.TxIndex />
      </Col.Grid>
    </Row.Grid>;
  };
};

module UndelegateMsg = {
  [@react.component]
  let make = (~delegation: TxSub.Msg.Undelegate.t) => {
    <Row.Grid>
      <Col.Grid col=Col.Six mb=24>
        <Heading value="Delegator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={delegation.delegatorAddress} />
      </Col.Grid>
      <Col.Grid col=Col.Six mb=24>
        <Heading value="Validator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={delegation.validatorAddress} accountType=`validator />
      </Col.Grid>
      <Col.Grid col=Col.Six>
        <Heading value="Amount" size=Heading.H5 marginBottom=8 />
        <AmountRender coins=[delegation.amount] pos=AmountRender.TxIndex />
      </Col.Grid>
    </Row.Grid>;
  };
};

module RedelegateMsg = {
  [@react.component]
  let make = (~delegation: TxSub.Msg.Redelegate.t) => {
    <Row.Grid>
      <Col.Grid col=Col.Six mb=24>
        <Heading value="Delegator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={delegation.delegatorAddress} />
      </Col.Grid>
      <Col.Grid col=Col.Six mb=24>
        <Heading value="Validator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={delegation.validatorSourceAddress} accountType=`validator />
      </Col.Grid>
      <Col.Grid col=Col.Six mbSm=24>
        <Heading value="Validator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={delegation.validatorDestinationAddress} accountType=`validator />
      </Col.Grid>
      <Col.Grid col=Col.Six>
        <Heading value="Amount" size=Heading.H5 marginBottom=8 />
        <AmountRender coins=[delegation.amount] pos=AmountRender.TxIndex />
      </Col.Grid>
    </Row.Grid>;
  };
};

module WithdrawRewardMsg = {
  [@react.component]
  let make = (~withdrawal: TxSub.Msg.WithdrawReward.success_t) => {
    <Row.Grid>
      <Col.Grid col=Col.Six mb=24>
        <Heading value="Delegator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={withdrawal.delegatorAddress} />
      </Col.Grid>
      <Col.Grid col=Col.Six mb=24>
        <Heading value="Validator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={withdrawal.validatorAddress} accountType=`validator />
      </Col.Grid>
      <Col.Grid col=Col.Six>
        <Heading value="Amount" size=Heading.H5 marginBottom=8 />
        <AmountRender coins={withdrawal.amount} pos=AmountRender.TxIndex />
      </Col.Grid>
    </Row.Grid>;
  };
};

module WithdrawRewardFailMsg = {
  [@react.component]
  let make = (~withdrawal: TxSub.Msg.WithdrawReward.fail_t) => {
    <Row.Grid>
      <Col.Grid col=Col.Six mb=24>
        <Heading value="Delegator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={withdrawal.delegatorAddress} />
      </Col.Grid>
      <Col.Grid col=Col.Six mb=24>
        <Heading value="Validator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={withdrawal.validatorAddress} accountType=`validator />
      </Col.Grid>
    </Row.Grid>;
  };
};

module WithdrawComissionMsg = {
  [@react.component]
  let make = (~withdrawal: TxSub.Msg.WithdrawCommission.success_t) => {
    <Row.Grid>
      <Col.Grid col=Col.Six mbSm=24>
        <Heading value="Validator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={withdrawal.validatorAddress} accountType=`validator />
      </Col.Grid>
      <Col.Grid col=Col.Six>
        <Heading value="Amount" size=Heading.H5 marginBottom=8 />
        <AmountRender coins={withdrawal.amount} pos=AmountRender.TxIndex />
      </Col.Grid>
    </Row.Grid>;
  };
};

module WithdrawComissionFailMsg = {
  [@react.component]
  let make = (~withdrawal: TxSub.Msg.WithdrawCommission.fail_t) => {
    <Row.Grid>
      <Col.Grid col=Col.Six mbSm=24>
        <Heading value="Validator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={withdrawal.validatorAddress} accountType=`validator />
      </Col.Grid>
    </Row.Grid>;
  };
};

module MultisendMsg = {
  module Styles = {
    open Css;
    let separatorLine =
      style([
        borderStyle(`none),
        backgroundColor(Colors.gray9),
        height(`px(1)),
        margin3(~top=`zero, ~h=`auto, ~bottom=`px(24)),
      ]);
  };

  [@react.component]
  let make = (~tx: TxSub.Msg.MultiSend.t) => {
    let isMobile = Media.isMobile();
    <>
      <Row.Grid>
        <Col.Grid col=Col.Six> <Heading value="From" size=Heading.H5 marginBottom=8 /> </Col.Grid>
        {isMobile
           ? React.null
           : <Col.Grid col=Col.Six>
               <Heading value="Amount" size=Heading.H5 marginBottom=8 />
             </Col.Grid>}
        {tx.inputs
         ->Belt_List.mapWithIndex((idx, input) =>
             <React.Fragment key={(idx |> string_of_int) ++ (input.address |> Address.toBech32)}>
               <Col.Grid col=Col.Six mb=24 mbSm=8>
                 <AddressRender address={input.address} />
               </Col.Grid>
               <Col.Grid col=Col.Six mb=24 mbSm=12>
                 <AmountRender coins={input.coins} pos=AmountRender.TxIndex />
               </Col.Grid>
             </React.Fragment>
           )
         ->Belt_List.toArray
         ->React.array}
      </Row.Grid>
      <hr className=Styles.separatorLine />
      <Heading value="To" size=Heading.H5 marginBottom=8 />
      <Row.Grid>
        {tx.outputs
         ->Belt_List.mapWithIndex((idx, output) =>
             <React.Fragment key={(idx |> string_of_int) ++ (output.address |> Address.toBech32)}>
               <Col.Grid col=Col.Six mb=24 mbSm=8>
                 <AddressRender address={output.address} />
               </Col.Grid>
               <Col.Grid col=Col.Six mb=24 mbSm=12>
                 <AmountRender coins={output.coins} pos=AmountRender.TxIndex />
               </Col.Grid>
             </React.Fragment>
           )
         ->Belt_List.toArray
         ->React.array}
      </Row.Grid>
    </>;
  };
};
