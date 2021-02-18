type client_t;

type reference_data_t = {
  pair: string,
  rate: float,
};

[@bs.module "@bandprotocol/bandchain.js"] [@bs.new]
external createClient: string => client_t = "Client";
[@bs.send]
external getReferenceData: (client_t, array(string)) => Js.Promise.t(array(reference_data_t)) =
  "getReferenceData";
