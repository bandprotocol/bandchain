module Styles = {
  open Css;

  let timeContainer = style([display(`inlineFlex)]);
};

[@react.component]
let make =
    (
      ~time,
      ~prefix="",
      ~suffix="",
      ~size=Text.Sm,
      ~weight=Text.Regular,
      ~spacing=Text.Unset,
      ~color=Colors.gray7,
      ~code=false,
      ~upper=false,
      ~height=Text.Px(16),
    ) => {
  <div className=Styles.timeContainer>
    {prefix != ""
       ? <>
           <Text value=prefix size weight spacing color code nowrap=true />
           <HSpacing size=Spacing.sm />
         </>
       : React.null}
    <Text
      value={time |> MomentRe.Moment.format("YYYY-MM-DD HH:mm:ss.SSS")}
      size
      weight
      spacing
      color
      code
      nowrap=true
      height
      block=true
    />
    {suffix != ""
       ? <>
           <HSpacing size=Spacing.sm />
           <Text value=suffix size weight spacing color code nowrap=true />
         </>
       : React.null}
  </div>;
};

module Grid = {
  [@react.component]
  let make =
      (
        ~time,
        ~prefix="",
        ~suffix="",
        ~size=Text.Sm,
        ~weight=Text.Regular,
        ~spacing=Text.Unset,
        ~color=Colors.gray7,
        ~code=false,
        ~upper=false,
        ~height=Text.Px(16),
        ~textAlign=Text.Left,
      ) => {
    <div>
      {prefix != ""
         ? <>
             <Text value=prefix size weight spacing color code nowrap=true />
             <HSpacing size=Spacing.sm />
           </>
         : React.null}
      <div>
        <Text
          value={time |> MomentRe.Moment.format("YYYY-MM-DD")}
          size
          weight
          spacing
          color
          code
          nowrap=true
          height
          block=true
          align=textAlign
        />
      </div>
      <div>
        <Text
          value={time |> MomentRe.Moment.format("HH:mm:ss.SSS")}
          size
          weight
          spacing
          color=Colors.gray6
          code
          nowrap=true
          height
          block=true
          align=textAlign
        />
      </div>
      {suffix != ""
         ? <>
             <HSpacing size=Spacing.sm />
             <Text value=suffix size weight spacing color code nowrap=true />
           </>
         : React.null}
    </div>;
  };
};
