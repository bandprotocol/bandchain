[@react.component]
let make = (~resolveStatus, ~style="") => {
  switch (resolveStatus) {
  | RequestSub.Success => <img src=Images.success className=style />
  | Failure => <img src=Images.fail className=style />
  | Pending => <img src=Images.pending className=style />
  | Expired => <img src=Images.expired className=style />
  | Unknown => <img src=Images.unknown className=style />
  };
};
