open Css;

[@react.component]
let make = (~dir) => {
  switch (dir) {
  | "left" => <div className={style([marginLeft(`auto)])} />
  | "right" => <div className={style([marginRight(`auto)])} />
  | "top" => <div className={style([marginTop(`auto)])} />
  | "bottom" => <div className={style([marginBottom(`auto)])} />
  | _ => React.null
  };
};
