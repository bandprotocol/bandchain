open Jest;
open Borsh;
open Expect;

describe("Expect Borsh to extract fields correctly", () => {
  test("should be able to extract fields from bytes correctly", () => {
    expect(
      Some([|
        {fieldName: "symbol", fieldType: "string"},
        {fieldName: "multiplier", fieldType: "u64"},
        {fieldName: "what", fieldType: "u8"},
      |]),
    )
    |> toEqual(
         extractFields(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
         ),
       )
  });

  test("should return None on invalid schema", () => {
    expect(None)
    |> toEqual(
         extractFields(
           {j|{"Input2": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
         ),
       )
  });
});

describe("Expect Borsh to encode correctly", () => {
  test("should be able to encode from bytes correctly", () => {
    expect(Some("0x03000000425443320000000000000064" |> JsBuffer.fromHex))
    |> toEqual(
         encode(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
           [|
             {fieldName: "symbol", fieldValue: "BTC"},
             {fieldName: "multiplier", fieldValue: "50"},
             {fieldName: "what", fieldValue: "100"},
           |],
         ),
       )
  });

  test("should be able to encode from bytes correctly2", () => {
    expect(Some("0x03000000425443900100000000000064" |> JsBuffer.fromHex))
    |> toEqual(
         encode(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
           [|
             {fieldName: "symbol", fieldValue: "BTC"},
             {fieldName: "multiplier", fieldValue: "400"},
             {fieldName: "what", fieldValue: "100"},
           |],
         ),
       )
  });

  test("should return None if invalid schema", () => {
    expect(None)
    |> toEqual(
         encode(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
           [|
             {fieldName: "symbol", fieldValue: "band"},
             {fieldName: "nulti", fieldValue: "400"},
             {fieldName: "what", fieldValue: "100"},
           |],
         ),
       )
  });

  test("should return None if invalid value", () => {
    expect(None)
    |> toEqual(
         encode(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
           [|
             {fieldName: "symbol", fieldValue: "band"},
             {fieldName: "multiplier", fieldValue: "400"},
             {fieldName: "what", fieldValue: "100"},
           |],
         ),
       )
  });
});

describe("Expect Borsh to decode correctly", () => {
  test("should be able to decode from bytes correctly", () => {
    expect(
      Some([|
        {fieldName: "symbol", fieldValue: "BTC"},
        {fieldName: "multiplier", fieldValue: "50"},
        {fieldName: "what", fieldValue: "100"},
      |]),
    )
    |> toEqual(
         decode(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
           "0x03000000425443320000000000000064" |> JsBuffer.fromHex,
         ),
       )
  });

  test("should be able to decode from bytes correctly2", () => {
    expect(
      Some([|
        {fieldName: "symbol", fieldValue: "band"},
        {fieldName: "multiplier", fieldValue: "400"},
        {fieldName: "what", fieldValue: "100"},
      |]),
    )
    |> toEqual(
         decode(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
           "0x0400000062616e64900100000000000064" |> JsBuffer.fromHex,
         ),
       )
  });

  test("should return None if invalid schema", () => {
    expect(None)
    |> toEqual(
         decode(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input2",
           "0x03000000425443320000000000000064" |> JsBuffer.fromHex,
         ),
       )
  });

  test("should return None if invalid bytes", () => {
    expect(None)
    |> toEqual(
         decode(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
           "0x030000004254433200000000000064" |> JsBuffer.fromHex,
         ),
       )
  });
});

describe("should be able to generate solidity correctly", () => {
  test("should be able to generate solidity", () => {
    expect(
      Some(
        {j|
pragma solidity ^0.5.0;

import "./Borsh.sol";

library ResultDecoder {
    using Borsh for Borsh.Data;

    struct Result {
        string symbol;
        uint64 multiplier;
        uint8 what;
    }

    function decodeResult(bytes memory _data)
        internal
        pure
        returns (Result memory result)
    {
        Borsh.Data memory data = Borsh.from(_data);
        result.symbol = string(data.decodeBytes());
        result.multiplier = data.decodeU64();
        result.what = data.decodeU8();
    }
}
|j},
      ),
    )
    |> toEqual(
         generateSolidity(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
         ),
       )
  });

  test("should be able to generate solidity 2", () => {
    expect(
      Some(
        {j|
pragma solidity ^0.5.0;

import "./Borsh.sol";

library ResultDecoder {
    using Borsh for Borsh.Data;

    struct Result {
        uint64 px;
    }

    function decodeResult(bytes memory _data)
        internal
        pure
        returns (Result memory result)
    {
        Borsh.Data memory data = Borsh.from(_data);
        result.px = data.decodeU64();
    }
}
|j},
      ),
    )
    |> toEqual(
         generateSolidity(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }","Output": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"px\\", \\"u64\\"] ] }"}|j},
           "Output",
         ),
       )
  });

  test("should return None if invalid class (solidity)", () => {
    expect(None)
    |> toEqual(
         generateSolidity(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input2",
         ),
       )
  });

  test("should return None if invalid type (solidity)", () => {
    expect(None)
    |> toEqual(
         generateSolidity(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"bytes\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
         ),
       )
  });
});

describe("should be able to generate go code correctly", () => {
  test("should be able to generate go code 1", () => {
    expect(
      Some(
        {j|package main

type Result struct {
	Symbol string
	Multiplier uint64
	What uint8
}

func DecodeResult(data []byte) (Result, error) {
	decoder := NewBorshDecoder(data)

	symbol, err := decoder.DecodeString()
	if err != nil {
		return Result{}, err
	}
	multiplier, err := decoder.DecodeU64()
	if err != nil {
		return Result{}, err
	}
	what, err := decoder.DecodeU8()
	if err != nil {
		return Result{}, err
	}

	if !decoder.Finished() {
		return Result{}, errors.New("Borsh: bytes left when decode result")
	}

	return Result{
		Symbol: symbol
		Multiplier: multiplier
		What: what
	}, nil
}|j},
      ),
    )
    |> toEqual(
         generateGo(
           "main",
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
         ),
       )
  });
  test("should be able to generate go code 2", () => {
    expect(
      Some(
        {j|package test

type Result struct {
	Px uint64
}

func DecodeResult(data []byte) (Result, error) {
	decoder := NewBorshDecoder(data)

	px, err := decoder.DecodeU64()
	if err != nil {
		return Result{}, err
	}

	if !decoder.Finished() {
		return Result{}, errors.New("Borsh: bytes left when decode result")
	}

	return Result{
		Px: px
	}, nil
}|j},
      ),
    )
    |> toEqual(
         generateGo(
           "test",
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }","Output": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"px\\", \\"u64\\"] ] }"}|j},
           "Output",
         ),
       )
  });
  test("should return None if invalid class (go)", () => {
    expect(None)
    |> toEqual(
         generateGo(
           "main",
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input2",
         ),
       )
  });
  test("should return None if invalid type (go)", () => {
    expect(None)
    |> toEqual(
         generateGo(
           "main",
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"bytes\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
         ),
       )
  });
});
