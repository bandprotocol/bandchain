type amount_t = {
  amount: string,
  denom: string,
};

type send_value_t = {
  oracleScriptID: string,
  calldata: string,
  requestedValidatorCount: string,
  sufficientValidatorCount: string,
  expiration: string,
  prepareGas: string,
  executeGas: string,
  sender: string,
};

type msg_t = {
  [@bs.as "type"]
  type_: string,
  value: send_value_t,
};

type fee_t = {
  amount: array(amount_t),
  gas: string,
};

type t = {
  msgs: array(msg_t),
  chain_id: string,
  fee: fee_t,
  memo: string,
  account_number: string,
  sequence: string,
};
