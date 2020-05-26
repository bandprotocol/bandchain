type modal_t =
  | Connect(string)
  | SubmitTx(SubmitMsg.t);

type t = {
  canExit: bool,
  modal: modal_t,
};

type a =
  | OpenModal(modal_t)
  | CloseModal
  | EnableExit
  | DisableExit;

let reducer = state =>
  fun
  | OpenModal(m) => Some({canExit: true, modal: m})
  | CloseModal => None
  | EnableExit => {
      switch (state) {
      | Some({modal}) => Some({canExit: true, modal})
      | None => None
      };
    }
  | DisableExit => {
      switch (state) {
      | Some({modal}) => Some({canExit: false, modal})
      | None => None
      };
    };

let context = React.createContext(ContextHelper.default: (option(t), a => unit));

[@react.component]
let make = (~children) => {
  React.createElement(
    React.Context.provider(context),
    {"value": React.useReducer(reducer, None), "children": children},
  );
};
