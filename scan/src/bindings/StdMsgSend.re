type amount_t = {
  amount: string,
  denom: string,
};

type send_value_t = {
  amount: array(amount_t),
  from_address: string,
  to_address: string,
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
