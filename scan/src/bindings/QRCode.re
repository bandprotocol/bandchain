[@bs.obj] external makeProps: (~value: string, ~size: int, unit) => _ = "";

[@bs.module]
external make:
  React.component({
    .
    "value": string,
    "size": int,
  }) =
  "qrcode.react";
