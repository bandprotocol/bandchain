[@bs.obj] external makeProps: (~children: React.element, ~className: string, unit) => _ = "";

[@bs.module "react-highlight"]
external make:
  React.component({
    .
    "children": React.element,
    "className": string,
  }) =
  "default";
