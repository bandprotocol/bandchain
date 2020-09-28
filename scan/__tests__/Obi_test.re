open Jest;
open Obi;
open Expect;

describe("Expect Obi to extract fields correctly", () => {
  test("should be able to extract fields from bytes correctly", () => {
    expect(
      Some([|
        {fieldName: "symbol", fieldType: "string"},
        {fieldName: "multiplier", fieldType: "u64"},
      |]),
    )
    |> toEqual(extractFields({j|{symbol:string,multiplier:u64}/{volume:u64}|j}, "input"))
  });

  test("should return None on invalid type", () => {
    expect(None)
    |> toEqual(extractFields({j|{symbol:string,multiplier:u64}/{volume:u64}|j}, "Input"))
  });
});

describe("Expect Obi encode correctly", () => {
  test("should be able to encode input (string, int) correctly", () => {
    expect(Some("0x00000003425443000000003b9aca00" |> JsBuffer.fromHex))
    |> toEqual(
         encode(
           {j|{symbol: string,multiplier: u64}/{price: u64,sources: [{ name: string, time: u64 }]}|j},
           "input",
           [|
             {fieldName: "symbol", fieldValue: "BTC"},
             {fieldName: "multiplier", fieldValue: "1000000000"},
           |],
         ),
       )
  });

  test("should be able to encode input (bytes) correctly", () => {
    expect(Some("0x0000000455555555" |> JsBuffer.fromHex))
    |> toEqual(
         encode(
           {j|{symbol: bytes}/{price: u64}|j},
           "input",
           [|{fieldName: "symbol", fieldValue: "0x55555555"}|],
         ),
       )
  });

  test("should be able to encode nested input correctly", () => {
    expect(Some("0x00000001000000020000000358585800000003595959" |> JsBuffer.fromHex))
    |> toEqual(
         encode(
           {j|{list: [{symbol: {name: [string]}}]}/{price: u64}|j},
           "input",
           [|{fieldName: "list", fieldValue: {j|[{"symbol": {"name": ["XXX", "YYY"]}}]|j}}|],
         ),
       )
  });

  test("should be able to encode output correctly", () => {
    expect(Some("0x0000000200000000000000780000000000000143" |> JsBuffer.fromHex))
    |> toEqual(
         encode(
           {j|{list: [{symbol: {name: [string]}}]}/{price: [u64]}|j},
           "output",
           [|{fieldName: "price", fieldValue: {j|[120, 323]|j}}|],
         ),
       )
  });

  test("should return None if invalid type", () => {
    expect(None)
    |> toEqual(
         encode(
           {j|{symbol: string,multiplier: u64}/{price: u64,sources: [{ name: string, time: u64 }]}|j},
           "nothing",
           [|
             {fieldName: "symbol", fieldValue: "BTC"},
             {fieldName: "multiplier", fieldValue: "1000000000"},
           |],
         ),
       )
  });

  test("should return None if invalid data", () => {
    expect(None)
    |> toEqual(
         encode(
           {j|{symbol: string,multiplier: u64}/{price: u64,sources: [{ name: string, time: u64 }]}|j},
           "nothing",
           [|{fieldName: "symbol", fieldValue: "BTC"}|],
         ),
       )
  });

  test("should return None if invalid input schema", () => {
    expect(None)
    |> toEqual(
         encode(
           {j|{symbol: string}/{price: u64,sources: [{ name: string, time: u64 }]}|j},
           "nothing",
           [|
             {fieldName: "symbol", fieldValue: "BTC"},
             {fieldName: "multiplier", fieldValue: "1000000000"},
           |],
         ),
       )
  });

  test("should return None if invalid output schema", () => {
    expect(None)
    |> toEqual(
         encode(
           {j|{symbol: string}|j},
           "nothing",
           [|{fieldName: "symbol", fieldValue: "BTC"}|],
         ),
       )
  });
});

describe("Expect Obi decode correctly", () => {
  test("should be able to decode from bytes correctly", () => {
    expect(
      Some([|
        {fieldName: "symbol", fieldValue: "\"BTC\""},
        {fieldName: "multiplier", fieldValue: "1000000000"},
      |]),
    )
    |> toEqual(
         decode(
           {j|{symbol: string,multiplier: u64}/{price: u64,sources: [{ name: string, time: u64 }]}|j},
           "input",
           "0x00000003425443000000003b9aca00" |> JsBuffer.fromHex,
         ),
       )
  });

  test("should be able to decode from bytes correctly (nested)", () => {
    expect(
      Some([|{fieldName: "list", fieldValue: "[{\"symbol\":{\"name\":[\"XXX\",\"YYY\"]}}]"}|]),
    )
    |> toEqual(
         decode(
           {j|{list: [{symbol: {name: [string]}}]}/{price: u64}|j},
           "input",
           "0x00000001000000020000000358585800000003595959" |> JsBuffer.fromHex,
         ),
       )
  });

  test("should be able to decode from bytes correctly (bytes)", () => {
    expect(Some([|{fieldName: "symbol", fieldValue: "0x55555555"}|]))
    |> toEqual(
         decode(
           {j|{symbol: bytes}/{price: u64}|j},
           "input",
           "0x0000000455555555" |> JsBuffer.fromHex,
         ),
       )
  });

  test("should return None if invalid schema", () => {
    expect(None)
    |> toEqual(
         decode(
           {j|{symbol: string}/{price: u64,sources: [{ name: string, time: u64 }]}|j},
           "input",
           "0x00000003425443000000003b9aca00" |> JsBuffer.fromHex,
         ),
       )
  });
});

describe("should be able to generate solidity correctly", () => {
  test("should be able to generate solidity", () => {
    expect(
      Some(
        {j|pragma solidity ^0.5.0;

import "./Obi.sol";

library ParamsDecoder {
    using Obi for Obi.Data;

    struct Params {
        string symbol;
        uint64 multiplier;
    }

    function decodeParams(bytes memory _data)
        internal
        pure
        returns (Params memory result)
    {
        Obi.Data memory data = Obi.from(_data);
        result.symbol = string(data.decodeBytes());
        result.multiplier = data.decodeU64();
    }
}

library ResultDecoder {
    using Obi for Obi.Data;

    struct Result {
        uint64 px;
    }

    function decodeResult(bytes memory _data)
        internal
        pure
        returns (Result memory result)
    {
        Obi.Data memory data = Obi.from(_data);
        result.px = data.decodeU64();
    }
}

|j},
      ),
    )
    |> toEqual(generateDecoderSolidity({j|{symbol:string,multiplier:u64}/{px:u64}|j}))
  })

  test("should be able to generate solidity when parameter is array", () => {
    expect(
      Some(
        {j|pragma solidity ^0.5.0;

import "./Obi.sol";

library ParamsDecoder {
    using Obi for Obi.Data;

    struct Params {
        string[] symbols;
        uint64 multiplier;
    }

    function decodeParams(bytes memory _data)
        internal
        pure
        returns (Params memory result)
    {
        Obi.Data memory data = Obi.from(_data);
        uint32 length = data.decodeU32();
        string[] memory symbols = new string[](length);
        for (uint256 i = 0; i < length; i++) {
          symbols[i] = string(data.decodeBytes());
        }
        result.symbols = symbols
        result.multiplier = data.decodeU64();
    }
}

library ResultDecoder {
    using Obi for Obi.Data;

    struct Result {
        uint64[] rates;
    }

    function decodeResult(bytes memory _data)
        internal
        pure
        returns (Result memory result)
    {
        Obi.Data memory data = Obi.from(_data);
        uint32 length = data.decodeU32();
        uint64[] memory rates = new uint64[](length);
        for (uint256 i = 0; i < length; i++) {
          rates[i] = data.decodeU64();
        }
        result.rates = rates
    }
}

|j},
      ),
    )
    |> toEqual(generateDecoderSolidity({j|{symbols:[string],multiplier:u64}/{rates:[u64]}|j}))
  })
});

describe("should be able to generate go code correctly", () => {
  // TODO: Change to real generated code once golang ParamsDecode is implemented
  test("should be able to generate go code 1", () => {
    expect(Some({j|"Code is not available."|j}))
    |> toEqual(
         generateDecoderGo("main", {j|{symbol:string,multiplier:u64}/{px:u64}|j}, Obi.Params),
       )
  });
  test("should be able to generate go code 2", () => {
    expect(
      Some(
        {j|package test

import "github.com/bandchain/chain/pkg/obi"

type Result struct {
	Px uint64
}

func DecodeResult(data []byte) (Result, error) {
	decoder := obi.NewObiDecoder(data)

	px, err := decoder.DecodeU64()
	if err != nil {
		return Result{}, err
	}

	if !decoder.Finished() {
		return Result{}, errors.New("Obi: bytes left when decode result")
	}

	return Result{
		Px: px
	}, nil
}|j},
      ),
    )
    |> toEqual(
         generateDecoderGo("test", {j|{symbol:string,multiplier:u64}/{px:u64}|j}, Obi.Result),
       )
  });
});

describe("should be able to generate encode go code correctly", () => {
  test("should be able to generate encode go code 1", () => {
    expect(
      Some(
        {j|package main

import "github.com/bandchain/chain/pkg/obi"

type Result struct {
	Symbol string
	Multiplier uint64
}

func(result *Result) EncodeResult() []byte {
	encoder := obi.NewObiEncoder()

	encoder.EncodeString(result.symbol)
	encoder.EncodeU64(result.multiplier)

	return encoder.GetEncodedData()
}|j},
      ),
    )
    |> toEqual(generateEncodeGo("main", {j|{symbol:string,multiplier:u64}/{px:u64}|j}, "input"))
  });
  test("should be able to generate encode go code 2", () => {
    expect(
      Some(
        {j|package test

import "github.com/bandchain/chain/pkg/obi"

type Result struct {
	Px uint64
}

func(result *Result) EncodeResult() []byte {
	encoder := obi.NewObiEncoder()

	encoder.EncodeU64(result.px)

	return encoder.GetEncodedData()
}|j},
      ),
    )
    |> toEqual(generateEncodeGo("test", {j|{symbol:string,multiplier:u64}/{px:u64}|j}, "output"))
  });
});

