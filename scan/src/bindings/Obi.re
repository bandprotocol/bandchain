type field_key_type_t = {
  fieldName: string,
  fieldType: string,
};

type field_key_value_t = {
  fieldName: string,
  fieldValue: string,
};

type data_type_t =
  | Params
  | Result;

let dataTypeToString =
  fun
  | Params => "Params"
  | Result => "Result";

let dataTypeToSchemaField =
  fun
  | Params => "input"
  | Result => "output";

let extractFields: (string, string) => option(array(field_key_type_t)) = [%bs.raw
  {|
  function(schema, t) {
    try {
      const normalizedSchema = schema.replace(/\s+/g, '')
      const tokens = normalizedSchema.split('/')
      let val
      if (t === 'input') {
        val = tokens[0]
      } else if (t === 'output') {
        val = tokens[1]
      } else {
        return undefined
      }
      let specs = val.slice(1, val.length - 1).split(',')
      return specs.map((spec) => {
        let x = spec.split(':')
        return {fieldName: x[0], fieldType: x[1]}
      })
    } catch {
      return undefined
    }
  }
|}
];

let encode: (string, string, array(field_key_value_t)) => option(JsBuffer.t) = [%bs.raw
  {|
function(schema, t, data) {
  const { Obi } = require('obi.js')

  try {
    const obi = new Obi(schema)

    let payload = {}
    for (let x of data) {
      let value = x.fieldValue
      if (x.fieldValue.startsWith("[") || x.fieldValue.startsWith("{")) {
        value = JSON.parse(x.fieldValue)
      } else if (x.fieldValue.startsWith("0x")) {
        value = Buffer.from(x.fieldValue.slice(2), "hex")
      }
      payload = {...payload, [x.fieldName]: value}
    }

    let buf
    if (t === 'input') {
      buf = obi.encodeInput(payload)
    } else if (t === 'output') {
      buf = obi.encodeOutput(payload)
    } else {
      return undefined
    }

    return Buffer.from(buf)
  } catch(err) {
    console.error(`Error encode ${t}`, err)
    return undefined
  }
}
  |}
];

let decode: (string, string, JsBuffer.t) => option(array(field_key_value_t)) = [%bs.raw
  {|
function(schema, t, data) {
  const { Obi } = require('obi.js')

  function stringify(data) {
    if (Array.isArray(data)) {
      return "[" + [...data].map(stringify).join(",") + "]"
    } else if (typeof(data) === "bigint") {
      return data.toString()
    } else if (Buffer.isBuffer(data)) {
      return "0x" + data.toString('hex')
    } else if (typeof(data) === "object") {
      return "{" + Object.entries(data).map(([k,v]) => JSON.stringify(k)+ ":" + stringify(v)).join(",") + "}"
    } else {
      return JSON.stringify(data)
    }
  }

  try {
    const obi = new Obi(schema)
    let rawResult
    if (t === 'input') {
      rawResult = obi.decodeInput(data)
    } else if (t === 'output') {
      rawResult = obi.decodeOutput(data)
    } else {
      return undefined
    }

    let result = []
    for (let x of Object.entries(rawResult)) {
      let value = stringify(x[1])
      result = [...result, {fieldName: x[0], fieldValue: value}]
    }

    return result
  } catch(err) {
    console.error(`Error decode ${t}`, err)
    return undefined
  }
}
  |}
];

type variable_t =
  | String
  | Bytes
  | U64
  | U32
  | U8;

type field_t = {
  name: string,
  varType: variable_t,
};

let parse = ({fieldName, fieldType}) => {
  let v =
    switch (fieldType |> ChangeCase.camelCase) {
    | "string" => Some(String)
    | "bytes" => Some(Bytes)
    | "u64" => Some(U64)
    | "u32" => Some(U32)
    | "u8" => Some(U8)
    | _ => None
    };

  let%Opt varType' = v;
  Some({name: fieldName, varType: varType'});
};

let declareSolidity = ({name, varType}) => {
  switch (varType) {
  | String => {j|string $name;|j}
  | Bytes => {j|bytes $name;|j}
  | U64 => {j|uint64 $name;|j}
  | U32 => {j|uint32 $name;|j}
  | U8 => {j|uint8 $name;|j}
  };
};

let assignSolidity = ({name, varType}) => {
  switch (varType) {
  | String => {j|result.$name = string(data.decodeBytes());|j}
  | Bytes => {j|result.$name = data.decodeBytes();|j}
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

let generateDecodeLibSolidity = (schema, dataType) => {
  let dataTypeString = dataType |> dataTypeToString;
  let name = dataType |> dataTypeToSchemaField;
  let template = (structs, functions) => {j|library $(dataTypeString)Decoder {
    using Obi for Obi.Data;

    struct $dataTypeString {
        $structs
    }

    function decode$dataTypeString(bytes memory _data)
        internal
        pure
        returns ($dataTypeString memory result)
    {
        Obi.Data memory data = Obi.from(_data);
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

let generateDecoderSolidity = schema => {
  let template = {j|pragma solidity ^0.5.0;

import "./Obi.sol";

|j};
  let paramsCodeOpt = generateDecodeLibSolidity(schema, Params);
  let resultCodeOpt = generateDecodeLibSolidity(schema, Result);
  let%Opt paramsCode = paramsCodeOpt;
  let%Opt resultCode = resultCodeOpt;
  Some(template ++ paramsCode ++ resultCode);
};

let declareGo = ({name, varType}) => {
  let capitalizedName = name |> ChangeCase.pascalCase;
  switch (varType) {
  | String => {j|$capitalizedName string|j}
  | Bytes => {j|$capitalizedName bytes|j}
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
  | Bytes => {j|$name, err := decoder.DecodeBytes()
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

// TODO: Implement input/params decoding
let generateDecoderGo = (packageName, schema, dataType) => {
  switch (dataType) {
  | Params => Some({j|"Code is not available."|j})
  | Result =>
    let name = dataType |> dataTypeToSchemaField;
    let template = (structs, functions, results) => {j|package $packageName

import "github.com/bandchain/chain/pkg/obi"

type Result struct {
\t$structs
}

func DecodeResult(data []byte) (Result, error) {
\tdecoder := obi.NewObiDecoder(data)

\t$functions

\tif !decoder.Finished() {
\t\treturn Result{}, errors.New("Obi: bytes left when decode result")
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
};

let encodeStructGo = ({name, varType}) => {
  switch (varType) {
  | U8 => {j|encoder.EncodeU8(result.$name)|j}
  | U32 => {j|encoder.EncodeU32(result.$name)|j}
  | U64 => {j|encoder.EncodeU64(result.$name)|j}
  | String => {j|encoder.EncodeString(result.$name)|j}
  | Bytes => {j|encoder.EncodeBytes(result.$name)|j}
  };
};

let generateEncodeGo = (packageName, schema, name) => {
  let template = (structs, functions) => {j|package $packageName

import "github.com/bandchain/chain/pkg/obi"

type Result struct {
\t$structs
}

func(result *Result) EncodeResult() []byte {
\tencoder := obi.NewObiEncoder()

\t$functions

\treturn encoder.GetEncodedData()
}|j};

  let%Opt fieldsPair = extractFields(schema, name);
  let%Opt fields = fieldsPair |> Belt_Array.map(_, parse) |> optionsAll;
  Some(
    template(
      fields |> Belt_Array.map(_, declareGo) |> Js.Array.joinWith("\n\t"),
      fields |> Belt_Array.map(_, encodeStructGo) |> Js.Array.joinWith("\n\t"),
    ),
  );
};
