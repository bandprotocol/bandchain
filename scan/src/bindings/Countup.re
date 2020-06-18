[@bs.deriving abstract]
type props = {
  start: float,
  [@bs.as "end"]
  end_: float,
  delay: int,
  decimals: int,
  duration: int,
  useEasing: bool,
  separator: string,
};

[@bs.deriving abstract]
type t = {
  countUp: float,
  update: float => unit,
};

[@bs.val] [@bs.module "react-countup"] external context: props => t = "useCountUp";
