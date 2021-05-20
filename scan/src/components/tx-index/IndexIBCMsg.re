module CreateClientMsg = {
  [@react.component]
  let make = (~info: TxSub.Msg.CreateClient.t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Chain ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.chainID} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Trusting Period" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.trustingPeriod} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Unboding Period" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.unbondingPeriod} />
      </Col>
    </Row>;
  };
};

module UpdateClientMsg = {
  [@react.component]
  let make = (~info: TxSub.Msg.UpdateClient.t) => {
    <Row>
      <Col col=Col.Six>
        <Heading value="Chain ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.chainID} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Client ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.clientID} />
      </Col>
    </Row>;
  };
};

module SubmitClientMisbehaviourMsg = {
  [@react.component]
  let make = (~info: TxSub.Msg.SubmitClientMisbehaviour.t) => {
    <Row>
      <Col col=Col.Six>
        <Heading value="Chain ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.chainID} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Client ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.clientID} />
      </Col>
    </Row>;
  };
};

module ConnectionOpenInitMsg = {
  [@react.component]
  let make = (~info: TxSub.Msg.ConnectionOpenInit.t) => {
    <Row>
      <Col col=Col.Six>
        <Heading value="Chain ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.chainID} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Client ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.clientID} />
      </Col>
    </Row>;
  };
};

module ConnectionOpenTryMsg = {
  [@react.component]
  let make = (~info: TxSub.Msg.ConnectionOpenTry.t) => {
    <Row>
      <Col col=Col.Six>
        <Heading value="Chain ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.chainID} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Client ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.clientID} />
      </Col>
    </Row>;
  };
};

module ConnectionOpenAckMsg = {
  [@react.component]
  let make = (~info: TxSub.Msg.ConnectionOpenAck.t) => {
    <Row>
      <Col col=Col.Six>
        <Heading value="Chain ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.chainID} />
      </Col>
    </Row>;
  };
};

module ConnectionOpenConfirmMsg = {
  [@react.component]
  let make = (~info: TxSub.Msg.ConnectionOpenConfirm.t) => {
    <Row>
      <Col col=Col.Six>
        <Heading value="Chain ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.chainID} />
      </Col>
    </Row>;
  };
};

module ChannelOpenInitMsg = {
  [@react.component]
  let make = (~info: TxSub.Msg.ChannelOpenInit.t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Chain ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.chainID} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Port ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.portID} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Channel ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.channelID} />
      </Col>
    </Row>;
  };
};

module ChannelOpenTryMsg = {
  [@react.component]
  let make = (~info: TxSub.Msg.ChannelOpenTry.t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Chain ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.chainID} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Port ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.portID} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Channel ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.channelID} />
      </Col>
    </Row>;
  };
};

module ChannelOpenAckMsg = {
  [@react.component]
  let make = (~info: TxSub.Msg.ChannelOpenAck.t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Chain ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.chainID} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Port ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.portID} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Channel ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.channelID} />
      </Col>
    </Row>;
  };
};

module ChannelOpenConfirmMsg = {
  [@react.component]
  let make = (~info: TxSub.Msg.ChannelOpenConfirm.t) => {
    // <Row>
    //   <Col col=Col.Six mb=24>
    //     <Heading value="Chain ID" size=Heading.H5 marginBottom=8 />
    //     <Text size=Text.Lg value={info.common.chainID} />
    //   </Col>
    //   <Col col=Col.Six mb=24>
    //     <Heading value="Port ID" size=Heading.H5 marginBottom=8 />
    //     <Text size=Text.Lg value={info.common.portID} />
    //   </Col>
    //   <Col col=Col.Six>
    //     <Heading value="Channel ID" size=Heading.H5 marginBottom=8 />
    //     <Text size=Text.Lg value={info.common.channelID} />
    //   </Col>
    // </Row>;
    React.null;
  };
};

module ChannelCloseInitMsg = {
  [@react.component]
  let make = (~info: TxSub.Msg.ChannelCloseInit.t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Chain ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.chainID} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Port ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.portID} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Channel ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.channelID} />
      </Col>
    </Row>;
  };
};

module ChannelCloseConfirmMsg = {
  [@react.component]
  let make = (~info: TxSub.Msg.ChannelCloseConfirm.t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Chain ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.chainID} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Port ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.portID} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Channel ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.channelID} />
      </Col>
    </Row>;
  };
};

module PacketMsg = {
  [@react.component]
  let make = (~info: TxSub.Msg.Packet.t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Data" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.data} breakAll=true />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Chain ID" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.chainID} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Sequence" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.sequence} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Source Port" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.sourcePort} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Source Channel" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.sourceChannel} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Destination Port" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.destinationPort} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Destination Channel" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.destinationChannel} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Timeout Height" size=Heading.H5 marginBottom=8 />
        <Text size=Text.Lg value={info.common.timeoutHeight} />
      </Col>
    </Row>;
  };
};
