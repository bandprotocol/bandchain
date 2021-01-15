module SendMsg = {
  [@react.component]
  let make = (~send: TxSub.Msg.Send.t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="From" size=Heading.H5 marginBottom=8 />
        <AddressRender address={send.fromAddress} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="To" size=Heading.H5 marginBottom=8 />
        <AddressRender address={send.toAddress} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Amount" size=Heading.H5 marginBottom=8 />
        <AmountRender coins={send.amount} pos=AmountRender.TxIndex />
      </Col>
    </Row>;
  };
};

module DelegateMsg = {
  [@react.component]
  let make = (~delegation: TxSub.Msg.Delegate.success_t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Delegator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={delegation.delegatorAddress} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Validator" size=Heading.H5 marginBottom=8 />
        <ValidatorMonikerLink
          validatorAddress={delegation.validatorAddress}
          moniker={delegation.moniker}
          identity={delegation.identity}
        />
      </Col>
      <Col col=Col.Six>
        <Heading value="Amount" size=Heading.H5 marginBottom=8 />
        <AmountRender coins=[delegation.amount] pos=AmountRender.TxIndex />
      </Col>
    </Row>;
  };
};

module DelegateFailMsg = {
  [@react.component]
  let make = (~delegation: TxSub.Msg.Delegate.fail_t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Delegator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={delegation.delegatorAddress} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Validator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={delegation.validatorAddress} accountType=`validator />
      </Col>
      <Col col=Col.Six>
        <Heading value="Amount" size=Heading.H5 marginBottom=8 />
        <AmountRender coins=[delegation.amount] pos=AmountRender.TxIndex />
      </Col>
    </Row>;
  };
};

module UndelegateMsg = {
  [@react.component]
  let make = (~undelegation: TxSub.Msg.Undelegate.success_t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Delegator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={undelegation.delegatorAddress} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Validator" size=Heading.H5 marginBottom=8 />
        <ValidatorMonikerLink
          validatorAddress={undelegation.validatorAddress}
          moniker={undelegation.moniker}
          identity={undelegation.identity}
        />
      </Col>
      <Col col=Col.Six>
        <Heading value="Amount" size=Heading.H5 marginBottom=8 />
        <AmountRender coins=[undelegation.amount] pos=AmountRender.TxIndex />
      </Col>
    </Row>;
  };
};

module UndelegateFailMsg = {
  [@react.component]
  let make = (~undelegation: TxSub.Msg.Undelegate.fail_t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Delegator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={undelegation.delegatorAddress} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Validator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={undelegation.validatorAddress} accountType=`validator />
      </Col>
      <Col col=Col.Six>
        <Heading value="Amount" size=Heading.H5 marginBottom=8 />
        <AmountRender coins=[undelegation.amount] pos=AmountRender.TxIndex />
      </Col>
    </Row>;
  };
};

module RedelegateMsg = {
  [@react.component]
  let make = (~redelegation: TxSub.Msg.Redelegate.success_t) => {
    <Row>
      <Col mb=24>
        <Heading value="Delegator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={redelegation.delegatorAddress} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Source Validator" size=Heading.H5 marginBottom=8 />
        <ValidatorMonikerLink
          validatorAddress={redelegation.validatorSourceAddress}
          moniker={redelegation.monikerSource}
          identity={redelegation.identitySource}
        />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Destination Validator" size=Heading.H5 marginBottom=8 />
        <ValidatorMonikerLink
          validatorAddress={redelegation.validatorDestinationAddress}
          moniker={redelegation.monikerDestination}
          identity={redelegation.identityDestination}
        />
      </Col>
      <Col>
        <Heading value="Amount" size=Heading.H5 marginBottom=8 />
        <AmountRender coins=[redelegation.amount] pos=AmountRender.TxIndex />
      </Col>
    </Row>;
  };
};

module RedelegateFailMsg = {
  [@react.component]
  let make = (~redelegation: TxSub.Msg.Redelegate.fail_t) => {
    <Row>
      <Col mb=24>
        <Heading value="Delegator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={redelegation.delegatorAddress} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Validator Source Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={redelegation.validatorSourceAddress} accountType=`validator />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Validator Destination Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={redelegation.validatorDestinationAddress} accountType=`validator />
      </Col>
      <Col>
        <Heading value="Amount" size=Heading.H5 marginBottom=8 />
        <AmountRender coins=[redelegation.amount] pos=AmountRender.TxIndex />
      </Col>
    </Row>;
  };
};

module WithdrawRewardMsg = {
  [@react.component]
  let make = (~withdrawal: TxSub.Msg.WithdrawReward.success_t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Delegator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={withdrawal.delegatorAddress} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Validator" size=Heading.H5 marginBottom=8 />
        <ValidatorMonikerLink
          validatorAddress={withdrawal.validatorAddress}
          moniker={withdrawal.moniker}
          identity={withdrawal.identity}
        />
      </Col>
      <Col col=Col.Six>
        <Heading value="Amount" size=Heading.H5 marginBottom=8 />
        <AmountRender coins={withdrawal.amount} pos=AmountRender.TxIndex />
      </Col>
    </Row>;
  };
};

module WithdrawRewardFailMsg = {
  [@react.component]
  let make = (~withdrawal: TxSub.Msg.WithdrawReward.fail_t) => {
    <Row>
      <Col col=Col.Six mbSm=24>
        <Heading value="Delegator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={withdrawal.delegatorAddress} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Validator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={withdrawal.validatorAddress} accountType=`validator />
      </Col>
    </Row>;
  };
};

module WithdrawComissionMsg = {
  [@react.component]
  let make = (~withdrawal: TxSub.Msg.WithdrawCommission.success_t) => {
    <Row>
      <Col col=Col.Six mbSm=24>
        <Heading value="Validator" size=Heading.H5 marginBottom=8 />
        <ValidatorMonikerLink
          validatorAddress={withdrawal.validatorAddress}
          moniker={withdrawal.moniker}
          identity={withdrawal.identity}
        />
      </Col>
      <Col col=Col.Six>
        <Heading value="Amount" size=Heading.H5 marginBottom=8 />
        <AmountRender coins={withdrawal.amount} pos=AmountRender.TxIndex />
      </Col>
    </Row>;
  };
};

module WithdrawComissionFailMsg = {
  [@react.component]
  let make = (~withdrawal: TxSub.Msg.WithdrawCommission.fail_t) => {
    <Row>
      <Col col=Col.Six>
        <Heading value="Validator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender address={withdrawal.validatorAddress} accountType=`validator />
      </Col>
    </Row>;
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
        margin3(~top=`px(8), ~h=`auto, ~bottom=`px(24)),
      ]);
  };

  [@react.component]
  let make = (~tx: TxSub.Msg.MultiSend.t) => {
    let isMobile = Media.isMobile();
    <>
      <Row>
        <Col col=Col.Six> <Heading value="From" size=Heading.H5 marginBottom=8 /> </Col>
        {isMobile
           ? React.null
           : <Col col=Col.Six> <Heading value="Amount" size=Heading.H5 marginBottom=8 /> </Col>}
        {tx.inputs
         ->Belt_List.mapWithIndex((idx, input) =>
             <React.Fragment key={(idx |> string_of_int) ++ (input.address |> Address.toBech32)}>
               <Col col=Col.Six mb=16 mbSm=8> <AddressRender address={input.address} /> </Col>
               <Col col=Col.Six mb=16 mbSm=12>
                 <AmountRender coins={input.coins} pos=AmountRender.TxIndex />
               </Col>
             </React.Fragment>
           )
         ->Belt_List.toArray
         ->React.array}
      </Row>
      <hr className=Styles.separatorLine />
      <Heading value="To" size=Heading.H5 marginBottom=8 />
      <Row>
        {tx.outputs
         ->Belt_List.mapWithIndex((idx, output) =>
             <React.Fragment key={(idx |> string_of_int) ++ (output.address |> Address.toBech32)}>
               <Col col=Col.Six mb=16 mbSm=8> <AddressRender address={output.address} /> </Col>
               <Col col=Col.Six mb=16 mbSm=12>
                 <AmountRender coins={output.coins} pos=AmountRender.TxIndex />
               </Col>
             </React.Fragment>
           )
         ->Belt_List.toArray
         ->React.array}
      </Row>
    </>;
  };
};
