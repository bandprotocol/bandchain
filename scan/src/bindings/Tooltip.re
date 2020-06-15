[@bs.obj]
external makeProps:
  (
    ~children: React.element,
    ~title: React.element,
    ~placement: string,
    ~arrow: bool,
    ~leaveDelay: int,
    unit
  ) =>
  _ =
  "";

[@bs.module "@material-ui/core"]
external make:
  React.component({
    .
    "children": React.element,
    "title": React.element,
    "placement": string,
    "arrow": bool,
    "leaveDelay": int,
  }) =
  "Tooltip";
