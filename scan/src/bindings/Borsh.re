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
    switch (varType |> ChangeCase.camelCase) {
    | "string" => Some(String)
    | "u64" => Some(U64)
    | "u32" => Some(U32)
    | "u8" => Some(U8)
    | _ => None
    };

  let%Opt varType' = v;
  Some({name, varType: varType'});
};

let declareSolidity = ({name, varType}) => {
  switch (varType) {
  | String => {j|string $name;|j}
  | U64 => {j|uint64 $name;|j}
  | U32 => {j|uint32 $name;|j}
  | U8 => {j|uint8 $name;|j}
  };
};

let assignSolidity = ({name, varType}) => {
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
      fields |> Belt_Array.map(_, declareSolidity) |> Js.Array.joinWith(indent),
      fields |> Belt_Array.map(_, assignSolidity) |> Js.Array.joinWith(indent),
    ),
  );
};

let declareGo = ({name, varType}) => {
  let capitalizedName = name |> ChangeCase.pascalCase;
  switch (varType) {
  | String => {j|$capitalizedName string|j}
  | U64 => {j|$capitalizedName uint64|j}
  | U32 => {j|$capitalizedName uint32|j}
  | U8 => {j|$capitalizedName uint8|j}
  };
};

let assignGo = ({name, varType}) => {
  switch (varType) {
  | String => {j|$name, err := decoder.DecodeString()
	if err != nil {
		return Result{}, err
	}|j}
  | U64 => {j|$name, err := decoder.DecodeU64()
	if err != nil {
		return Result{}, err
	}|j}
  | U32 => {j|$name, err := decoder.DecodeU32()
	if err != nil {
		return Result{}, err
	}|j}
  | U8 => {j|$name, err := decoder.DecodeU8()
	if err != nil {
		return Result{}, err
	}|j}
  };
};

let resultGo = ({name}) => {
  let capitalizedName = name |> ChangeCase.pascalCase;
  {j|$capitalizedName: $name|j};
};

let generateGo = (packageName, schema, name) => {
  let template = (structs, functions, results) => {j|package $packageName

type Result struct {
\t$structs
}

func DecodeResult(data []byte) (Result, error) {
\tdecoder := NewBorshDecoder(data)

\t$functions

\tif !decoder.Finished() {
\t\treturn Result{}, errors.New("Borsh: bytes left when decode result")
\t}

\treturn Result{
\t\t$results
\t}, nil
}|j};

  let%Opt fieldsPair = extractFields(schema, name);
  let%Opt fields = fieldsPair |> Belt_Array.map(_, parse) |> optionsAll;
  Some(
    template(
      fields |> Belt_Array.map(_, declareGo) |> Js.Array.joinWith("\n\t"),
      fields |> Belt_Array.map(_, assignGo) |> Js.Array.joinWith("\n\t"),
      fields |> Belt_Array.map(_, resultGo) |> Js.Array.joinWith("\n\t\t"),
    ),
  );
};
