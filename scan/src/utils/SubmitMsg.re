type t =
  | Send(option(Address.t))
  | Delegate(Address.t)
  | Undelegate(Address.t)
  | Redelegate(Address.t)
  | WithdrawReward(Address.t)
  | Reinvest(Address.t, float)
  | Vote(ID.Proposal.t, string);

let toString =
  fun
  | Send(_) => "Send Token"
  | Delegate(_) => "Delegate"
  | Undelegate(_) => "Undelegate"
  | Redelegate(_) => "Redelegate"
  | WithdrawReward(_) => "Withdraw Reward"
  | Reinvest(_) => "Reinvest"
  | Vote(_) => "Vote";

let gasLimit =
  fun
  | Send(_)
  | Delegate(_)
  | Undelegate(_)
  | Vote(_)
  | WithdrawReward(_)
  | Reinvest(_) => 200000
  | Redelegate(_) => 300000;
