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

let create = (fromAddress, toAddress, sendAmount, accountNumber, sequence) => {
  {
    msgs: [|
      {
        type_: "cosmos-sdk/MsgSend",
        value: {
          amount: [|{amount: sendAmount |> string_of_int, denom: "uband"}|],
          from_address: fromAddress |> Address.toBech32,
          to_address: toAddress |> Address.toBech32,
        },
      },
    |],
    chain_id: "banchain",
    fee: {
      amount: [|{amount: "5000", denom: "uband"}|],
      gas: "200000",
    },
    memo: "",
    account_number: accountNumber,
    sequence,
  };
};
