module SubmitProposalMsg = {
  [@react.component]
  let make = (~proposal: TxSub.Msg.SubmitProposal.success_t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Proposer" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={proposal.proposer} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Proposal ID" size=Heading.H5 marginBottom=8 />
        <TypeID.Proposal position=TypeID.Subtitle id={proposal.proposalID} />
      </Col>
      <Col col=Col.Six mbSm=24>
        <Heading value="Title" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={proposal.title} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Deposit Amount" size=Heading.H5 marginBottom=8 />
        <AmountRender coins={proposal.initialDeposit} pos=AmountRender.TxIndex />
      </Col>
    </Row>;
  };
};

module SubmitProposalFailMsg = {
  [@react.component]
  let make = (~proposal: TxSub.Msg.SubmitProposal.fail_t) => {
    <Row>
      <Col mb=24>
        <Heading value="Proposer" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={proposal.proposer} />
      </Col>
      <Col col=Col.Six mbSm=24>
        <Heading value="Title" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={proposal.title} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Deposit Amount" size=Heading.H5 marginBottom=8 />
        <AmountRender coins={proposal.initialDeposit} pos=AmountRender.TxIndex />
      </Col>
    </Row>;
  };
};

module DepositMsg = {
  [@react.component]
  let make = (~deposit: TxSub.Msg.Deposit.success_t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Depositor" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={deposit.depositor} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Proposal ID" size=Heading.H5 marginBottom=8 />
        <TypeID.Proposal position=TypeID.Subtitle id={deposit.proposalID} />
      </Col>
      <Col col=Col.Six mbSm=24>
        <Heading value="Title" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={deposit.title} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Amount" size=Heading.H5 marginBottom=8 />
        <AmountRender coins={deposit.amount} pos=AmountRender.TxIndex />
      </Col>
    </Row>;
  };
};

module DepositFailMsg = {
  [@react.component]
  let make = (~deposit: TxSub.Msg.Deposit.fail_t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Depositor" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={deposit.depositor} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Proposal ID" size=Heading.H5 marginBottom=8 />
        <TypeID.Proposal position=TypeID.Subtitle id={deposit.proposalID} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Amount" size=Heading.H5 marginBottom=8 />
        <AmountRender coins={deposit.amount} pos=AmountRender.TxIndex />
      </Col>
    </Row>;
  };
};

module VoteMsg = {
  [@react.component]
  let make = (~vote: TxSub.Msg.Vote.success_t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Voter" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={vote.voterAddress} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Proposal ID" size=Heading.H5 marginBottom=8 />
        <TypeID.Proposal position=TypeID.Subtitle id={vote.proposalID} />
      </Col>
      <Col col=Col.Six mbSm=24>
        <Heading value="Title" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={vote.title} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Option" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={vote.option} />
      </Col>
    </Row>;
  };
};

module VoteFailMsg = {
  [@react.component]
  let make = (~vote: TxSub.Msg.Vote.fail_t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Voter" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={vote.voterAddress} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Proposal ID" size=Heading.H5 marginBottom=8 />
        <TypeID.Proposal position=TypeID.Subtitle id={vote.proposalID} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Option" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={vote.option} />
      </Col>
    </Row>;
  };
};
