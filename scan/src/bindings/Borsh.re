let extractFields: (string, string) => option(array((string, string))) = [%bs.raw
  {|
  function(_schema, cls) {
    try {
      return JSON.parse(JSON.parse(_schema)[cls])["fields"]
    } catch(err) {
      return undefined
    }
  }
|}
];

let decode: (string, string, JsBuffer.t) => option(array((string, string))) = [%bs.raw
  {|
function(_schema, cls, data) {
  const borsh = require('borsh')

  function gen(members) {
    return function(...args) {
      for (let i in members) {
        this[members[i]] = args[i]
      }
    }
  }

  try {
    let schema = JSON.parse(_schema)
    let schemaMap = new Map()
    for (let className in schema) {
      let t = JSON.parse(schema[className])
      if (t.kind == "struct") {
        window[className] = gen(t.fields.map(x => x[0]))
        t.fields = t.fields.map(x => {
          x[1] = x[1].replace(/^\w/, c => c.toUpperCase());
          return x
        })
        schemaMap.set(window[className], t)
      }
    }
    let model = window[cls]
    let newValue = borsh.deserialize(schemaMap, model, data)
    return schemaMap.get(model).fields.map(([fieldName, _]) => {
      return [fieldName, newValue[fieldName].toString()]
    });
  } catch(err) {
    return undefined
  }
}
|}
];

let encode: (string, string, array((string, string))) => option(JsBuffer.t) = [%bs.raw
  {|
function(_schema, cls, data) {
  const borsh = require('borsh')

  function gen(members) {
    return function(args) {
      for (let i in members) {
        this[members[i]] = args[members[i]]
      }
    }
  }

  try {
    let schema = JSON.parse(_schema)
    let schemaMap = new Map()
    for (let className in schema) {
      let t = JSON.parse(schema[className])
      if (t.kind == "struct") {
        window[className] = gen(t.fields.map(x => x[0]))
        t.fields = t.fields.map(x => {
          x[1] = x[1].replace(/^\w/, c => c.toUpperCase());
          return x
        })
        schemaMap.set(window[className], t)
      }
    }

    let rawValue = {}
    let specs = schemaMap.get(window[cls])

    if (!specs.fields) return undefined

    for (let i in data) {
      let isFound = specs.fields.some(x => x[0] == data[i][0])
      if (!isFound) return undefined

      rawValue[data[i][0]] = data[i][1]
    }
    let value = new window[cls](rawValue)

    if (specs.fields.length !== data.length) return undefined

    let buf = borsh.serialize(schemaMap, value)
    return Buffer.from(buf)
  } catch(err) {
    return undefined
  }
}
|}
];

type variable_t =
  | String
  | U64
  | U32
  | U8;

type field_t = {
  name: string,
  varType: variable_t,
};

let parse = ((name, varType)) => {
  let v =
    switch (varType |> String.lowercase_ascii) {
    | "string" => Some(String)
    | "u64" => Some(U64)
    | "u32" => Some(U32)
    | "u8" => Some(U8)
    | _ => None
    };

  let%Opt varType' = v;
  Some({name, varType: varType'});
};

let declare = ({name, varType}) => {
  switch (varType) {
  | String => {j|string $name;|j}
  | U64 => {j|uint64 $name;|j}
  | U32 => {j|uint32 $name;|j}
  | U8 => {j|uint8 $name;|j}
  };
};

let assign = ({name, varType}) => {
  switch (varType) {
  | String => {j|result.$name = string(data.decodeBytes());|j}
  | U64 => {j|result.$name = data.decodeU64();|j}
  | U32 => {j|result.$name = data.decodeU32();|j}
  | U8 => {j|result.$name = data.decodeU8();|j}
  };
};

let optionsAll = options =>
  options
  |> Belt_Array.reduce(_, Some([||]), (acc, obj) => {
       switch (acc, obj) {
       | (Some(acc'), Some(obj')) => Some(acc' |> Js.Array.concat([|obj'|]))
       | (_, _) => None
       }
     });

let generateSolidity = (schema, name) => {
  let template = (structs, functions) => {j|
pragma solidity ^0.5.0;

import "./Borsh.sol";

library ResultDecoder {
    using Borsh for Borsh.Data;

    struct Result {
        $structs
    }

    function decodeResult(bytes memory _data)
        internal
        pure
        returns (Result memory result)
    {
        Borsh.Data memory data = Borsh.from(_data);
        $functions
    }
}
|j};

  let%Opt fieldsPairs = extractFields(schema, name);
  let%Opt fields = fieldsPairs |> Belt_Array.map(_, parse) |> optionsAll;
  let indent = "\n        ";
  Some(
    template(
      fields |> Belt_Array.map(_, declare) |> Js.Array.joinWith(indent),
      fields |> Belt_Array.map(_, assign) |> Js.Array.joinWith(indent),
    ),
  );
};
