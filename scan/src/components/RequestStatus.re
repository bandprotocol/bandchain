
type t =
  | Pending
  | Success
  | Failure
  | Expired
  | Unknown;

type display_t =
  | Full
  | Mini;

let fromJsonString = json => {
  let status = json |> Js.Json.decodeString |> Belt_Option.getExn;
  switch (status) {
  | "Open" => Pending
  | "Success" => Success
  | "Failure" => Failure
  | "Expired" => Expired
  | _ => Unknown
  };
};

let fromInt =
  fun
  | 0 => Pending
  | 1 => Success
  | 2 => Failure
  | 3 => Expired
  | _ => Unknown;

let toString =
  fun
  | Success => "Success"
  | Failure => "Failure"
  | Pending => "Pending"
  | Expired => "Expired"
  | Unknown => "Unknown";

[@react.component]
let make = (~resolveStatus, ~display=Mini, ~style="") => {
  <div className={CssHelper.flexBox(~align=`center, ())}>
    {switch (resolveStatus) {
     | Success => <img src=Images.success className=style />
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
