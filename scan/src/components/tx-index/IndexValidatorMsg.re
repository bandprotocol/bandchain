module CreateValidatorMsg = {
  [@react.component]
  let make = (~validator: TxSub.Msg.CreateValidator.t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Moniker" size=Heading.H5 marginBottom=8 />
        <ValidatorMonikerLink
          validatorAddress={validator.validatorAddress}
          moniker={validator.moniker}
          identity={validator.identity}
          width={`percent(100.)}
          avatarWidth=20
          size=Text.Lg
        />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Identity" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={validator.identity} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Commission Rate" size=Heading.H5 marginBottom=8 />
        <Text
          size=Text.Lg
          value={
            (validator.commissionRate *. 100.)->Js.Float.toFixedWithPrecision(~digits=4) ++ "%"
          }
        />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Commission Max Rate" size=Heading.H5 marginBottom=8 />
        <Text
          size=Text.Lg
          value={
            (validator.commissionMaxRate *. 100.)->Js.Float.toFixedWithPrecision(~digits=4) ++ "%"
          }
        />
      </Col>
      <Col mb=24>
        <Heading value="Commission Max Change" size=Heading.H5 marginBottom=8 />
        <Text
          size=Text.Lg
          value={
            (validator.commissionMaxChange *. 100.)->Js.Float.toFixedWithPrecision(~digits=4)
            ++ "%"
          }
        />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Delegator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={validator.delegatorAddress} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Validator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender
          position=AddressRender.Subtitle
          address={validator.validatorAddress}
          accountType=`validator
        />
      </Col>
      <Col mb=24>
        <Heading value="Public Key" size=Heading.H5 marginBottom=8 />
        <PubKeyRender pubKey={validator.publicKey} alignLeft=true position=PubKeyRender.Subtitle />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Min Self Delegation" size=Heading.H5 marginBottom=8 />
        <AmountRender coins=[validator.minSelfDelegation] pos=AmountRender.TxIndex />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Self Delegation" size=Heading.H5 marginBottom=8 />
        <AmountRender coins=[validator.selfDelegation] pos=AmountRender.TxIndex />
      </Col>
      <Col mb=24>
        <Heading value="Details" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={validator.details} />
      </Col>
      <Col>
        <Heading value="Website" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={validator.website} />
      </Col>
    </Row>;
  };
};

module EditValidatorMsg = {
  [@react.component]
  let make = (~validator: BandScan.TxSub.Msg.EditValidator.t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Moniker" size=Heading.H5 marginBottom=8 />
        <Text
          value={validator.moniker == Config.doNotModify ? "Unchanged" : validator.moniker}
          size=Text.Lg
        />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Identity" size=Heading.H5 marginBottom=8 />
        <Text
          size=Text.Lg
          value={validator.identity == Config.doNotModify ? "Unchanged" : validator.identity}
        />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Commission Rate" size=Heading.H5 marginBottom=8 />
        <Text
          size=Text.Lg
          value={
            switch (validator.commissionRate) {
            | Some(rate) => (rate *. 100.)->Js.Float.toFixedWithPrecision(~digits=4) ++ "%"
            | None => "Unchanged"
            }
          }
        />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Validator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender
          position=AddressRender.Subtitle
          address={validator.sender}
          accountType=`validator
        />
      </Col>
      <Col mb=24>
        <Heading value="Min Self Delegation" size=Heading.H5 marginBottom=8 />
        {switch (validator.minSelfDelegation) {
         | Some(minSelfDelegation') =>
           <AmountRender coins=[minSelfDelegation'] pos=AmountRender.TxIndex />
         | None => <Text value="Unchanged" size=Text.Lg />
         }}
      </Col>
      <Col>
        <Heading value="Details" size=Heading.H5 marginBottom=8 />
        <Text
          size=Text.Lg
          value={validator.details == Config.doNotModify ? "Unchanged" : validator.details}
        />
      </Col>
    </Row>;
  };
};

module UnjailMsg = {
  [@react.component]
  let make = (~unjail: TxSub.Msg.Unjail.t) => {
    <Row>
      <Col col=Col.Six>
        <Heading value="Validator" size=Heading.H5 marginBottom=8 />
        <AddressRender
          position=AddressRender.Subtitle
          address={unjail.address}
          accountType=`validator
        />
      </Col>
    </Row>;
  };
};

module AddReporterMsg = {
  [@react.component]
  let make = (~address: TxSub.Msg.AddReporter.success_t) => {
    <Row>
      <Col col=Col.Six mbSm=24>
        <Heading value="Validator" size=Heading.H5 marginBottom=8 />
        <AddressRender
          position=AddressRender.Subtitle
          address={address.validator}
          accountType=`validator
        />
      </Col>
      <Col col=Col.Six>
        <Heading value="Reporter Address" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={address.reporter} />
      </Col>
    </Row>;
  };
};

module AddReporterFailMsg = {
  [@react.component]
  let make = (~address: TxSub.Msg.AddReporter.fail_t) => {
    <Row>
      <Col col=Col.Six mbSm=24>
        <Heading value="Validator" size=Heading.H5 marginBottom=8 />
        <AddressRender
          position=AddressRender.Subtitle
          address={address.validator}
          accountType=`validator
        />
      </Col>
      <Col col=Col.Six>
        <Heading value="Reporter Address" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={address.reporter} />
      </Col>
    </Row>;
  };
};

module RemoveReporterMsg = {
  [@react.component]
  let make = (~address: TxSub.Msg.RemoveReporter.success_t) => {
    <Row>
      <Col col=Col.Six mbSm=24>
        <Heading value="Validator" size=Heading.H5 marginBottom=8 />
        <AddressRender
          position=AddressRender.Subtitle
          address={address.validator}
          accountType=`validator
        />
      </Col>
      <Col col=Col.Six>
        <Heading value="Reporter Address" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={address.reporter} />
      </Col>
    </Row>;
  };
};

module RemoveReporterFailMsg = {
  [@react.component]
  let make = (~address: TxSub.Msg.RemoveReporter.fail_t) => {
    <Row>
      <Col col=Col.Six mbSm=24>
        <Heading value="Validator" size=Heading.H5 marginBottom=8 />
        <AddressRender
          position=AddressRender.Subtitle
          address={address.validator}
          accountType=`validator
        />
      </Col>
      <Col col=Col.Six>
        <Heading value="Reporter Address" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={address.reporter} />
      </Col>
    </Row>;
  };
};

module ActivateMsg = {
  [@react.component]
  let make = (~activate: TxSub.Msg.Activate.t) => {
    <Row>
      <Col col=Col.Six>
        <Heading value="Validator" size=Heading.H5 marginBottom=8 />
        <AddressRender
          position=AddressRender.Subtitle
          address={activate.validatorAddress}
          accountType=`validator
        />
      </Col>
    </Row>;
  };
};

module SetWithdrawAddressMsg = {
  [@react.component]
  let make = (~set: TxSub.Msg.SetWithdrawAddress.t) => {
    <Row>
      <Col col=Col.Six mbSm=24>
        <Heading value="Delegator Address" size=Heading.H5 marginBottom=8 />
        <AddressRender
          position=AddressRender.Subtitle
          address={set.delegatorAddress}
          accountType=`validator
        />
      </Col>
      <Col col=Col.Six>
        <Heading value="Withdraw Address" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={set.withdrawAddress} />
      </Col>
    </Row>;
  };
};
