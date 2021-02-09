type display_t =
  | Full
  | Mini;

module Sub = {
  let toString =
    fun
    | RequestSub.Success => "Success"
    | Failure => "Failure"
    | Pending => "Pending"
    | Expired => "Expired"
    | Unknown => "Unknown";

  [@react.component]
  let make = (~resolveStatus, ~display=Mini, ~style="") => {
    <div className={CssHelper.flexBox(~align=`center, ())}>
      {switch (resolveStatus) {
       | RequestSub.Success => <img src=Images.success className=style />
       | Failure => <img src=Images.fail className=style />
       | Pending => <img src=Images.pending className=style />
       | Expired => <img src=Images.expired className=style />
       | Unknown => <img src=Images.unknown className=style />
       }}
      {display == Full
         ? <>
             <HSpacing size=Spacing.sm />
             <Text value={resolveStatus |> toString} size=Text.Lg />
           </>
         : React.null}
    </div>;
  };
};

module Query = {
  let toString =
    fun
    | RequestQuery.Success => "Success"
    | Failure => "Failure"
    | Pending => "Pending"
    | Expired => "Expired"
    | Unknown => "Unknown";

  [@react.component]
  let make = (~resolveStatus, ~display=Mini, ~style="") => {
    <div className={CssHelper.flexBox(~align=`center, ())}>
      {switch (resolveStatus) {
       | RequestQuery.Success => <img src=Images.success className=style />
       | Failure => <img src=Images.fail className=style />
       | Pending => <img src=Images.pending className=style />
       | Expired => <img src=Images.expired className=style />
       | Unknown => <img src=Images.unknown className=style />
       }}
      {display == Full
         ? <>
             <HSpacing size=Spacing.sm />
             <Text value={resolveStatus |> toString} size=Text.Lg />
           </>
         : React.null}
    </div>;
  };
};
