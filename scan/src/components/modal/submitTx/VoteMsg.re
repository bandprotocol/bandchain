module Styles = {
  open Css;
  let buttonGroup =
    style([
      margin2(~v=`px(30), ~h=`px(-12)),
      selector(
        "> button",
        [
          flexGrow(0.),
          flexShrink(0.),
          flexBasis(`calc((`sub, `percent(25.), `px(24)))),
          margin2(~v=`zero, ~h=`px(12)),
          disabled([backgroundColor(Colors.bandBlue), color(Colors.white)]),
        ],
      ),
    ]);
};

module VoteInput = {
  [@react.component]
  let make = (~setAnswerOpt, ~answerOpt) => {
    <div className={Css.merge([CssHelper.flexBox(), Styles.buttonGroup])}>
      <Button
        variant=Button.Outline
        px=15
        py=8
        onClick={_ => setAnswerOpt(_ => Some("Yes"))}
        disabled={answerOpt == Some("Yes")}>
        <Text size=Text.Md value="Yes" />
      </Button>
      <Button
        variant=Button.Outline
        px=15
        py=8
        onClick={_ => setAnswerOpt(_ => Some("No"))}
        disabled={answerOpt == Some("No")}>
        <Text size=Text.Md value="No" />
      </Button>
      <Button
        variant=Button.Outline
        px=15
        py=8
        onClick={_ => setAnswerOpt(_ => Some("NoWithVeto"))}
        disabled={answerOpt == Some("NoWithVeto")}>
        <Text size=Text.Md value="No with Veto" />
      </Button>
      <Button
        variant=Button.Outline
        px=15
        py=8
        onClick={_ => setAnswerOpt(_ => Some("Abstain"))}
        disabled={answerOpt == Some("Abstain")}>
        <Text size=Text.Md value="Abstain" />
      </Button>
    </div>;
  };
};

[@react.component]
let make = (~proposalID, ~proposalName, ~setMsgsOpt) => {
  let (answerOpt, setAnswerOpt) = React.useState(_ => None);
  React.useEffect1(
    _ => {
      let msgsOpt = {
        let%Opt answer = answerOpt;
        Some([|TxCreator.Vote(proposalID, answer)|]);
      };
      setMsgsOpt(_ => msgsOpt);
      None;
    },
    [|answerOpt|],
  );

  <>
    <div className={CssHelper.flexBox(~justify=`spaceBetween, ())}>
      <Text value="Vote To" size=Text.Lg spacing={Text.Em(0.03)} nowrap=true block=true />
      <div className={CssHelper.flexBox()}>
        <TypeID.Proposal id=proposalID position=TypeID.Subtitle />
        <HSpacing size=Spacing.sm />
        <Heading size=Heading.H5 value=proposalName />
      </div>
    </div>
    <VoteInput answerOpt setAnswerOpt />
  </>;
};
