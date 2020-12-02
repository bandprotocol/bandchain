module Styles = {
  open Css;

  let container = style([paddingBottom(`px(24))]);

  let buttonGroup =
    style([
      margin4(~top=`px(0), ~right=`px(-12), ~bottom=`px(23), ~left=`px(-12)),
      selector(
        "> button",
        [
          flexGrow(0.),
          flexShrink(0.),
          flexBasis(`calc((`sub, `percent(50.), `px(24)))),
          margin2(~v=`zero, ~h=`px(12)),
          disabled([backgroundColor(Colors.bandBlue), color(Colors.white)]),
          border(`px(1), `solid, Colors.gray9),
          color(Colors.gray7),
          fontWeight(`light),
        ],
      ),
    ]);
};

module VoteInput = {
  [@react.component]
  let make = (~setAnswerOpt, ~answerOpt) => {
    <>
      <div className={Css.merge([CssHelper.flexBox(), Styles.buttonGroup])}>
        <Button
          variant=Button.Outline
          px=15
          py=9
          onClick={_ => setAnswerOpt(_ => Some("Yes"))}
          disabled={answerOpt == Some("Yes")}>
          <Text size=Text.Lg value="Yes" />
        </Button>
        <Button
          variant=Button.Outline
          px=15
          py=9
          onClick={_ => setAnswerOpt(_ => Some("No"))}
          disabled={answerOpt == Some("No")}>
          <Text size=Text.Lg value="No" />
        </Button>
      </div>
      <div className={Css.merge([CssHelper.flexBox(), Styles.buttonGroup])}>
        <Button
          variant=Button.Outline
          px=15
          py=9
          onClick={_ => setAnswerOpt(_ => Some("NoWithVeto"))}
          disabled={answerOpt == Some("NoWithVeto")}>
          <Text size=Text.Lg value="No with Veto" />
        </Button>
        <Button
          variant=Button.Outline
          px=15
          py=9
          onClick={_ => setAnswerOpt(_ => Some("Abstain"))}
          disabled={answerOpt == Some("Abstain")}>
          <Text size=Text.Lg value="Abstain" />
        </Button>
      </div>
    </>;
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
    <div className=Styles.container>
      <Text value="Vote to" size=Text.Md weight=Text.Medium nowrap=true block=true />
      <VSpacing size=Spacing.sm />
      <div className={CssHelper.flexBox()}>
        <TypeID.Proposal id=proposalID position=TypeID.Subtitle />
        <HSpacing size=Spacing.sm />
        <Text value=proposalName size=Text.Lg nowrap=true block=true />
      </div>
    </div>
    <VoteInput answerOpt setAnswerOpt />
  </>;
};
