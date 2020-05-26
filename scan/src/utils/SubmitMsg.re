type t =
  | Send(option(Address.t))
  | Delegate(Address.t)
  | Undelegate(Address.t)
  | Redelegate(Address.t)
  | WithdrawReward(Address.t);

let toString =
  fun
  | Send(_) => "Send"
  | Delegate(_) => "Delegate"
  | Undelegate(_) => "Undelegate"
  | Redelegate(_) => "Redelegate"
  | WithdrawReward(_) => "Withdraw Reward";

let gasLimit =
  fun
  | Send(_) => 70000
  | Delegate(_) => 160000
  | Undelegate(_) => 180000
  | Redelegate(_) => 210000
  | WithdrawReward(_) => 110000;
