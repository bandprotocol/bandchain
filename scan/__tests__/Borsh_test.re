open Jest;
open Borsh;
open Expect;

describe("Expect Borsh to encode correctly", () => {
  test("should be able to encode from bytes correctly", () => {
    expect(Some("0x03000000425443320000000000000064" |> JsBuffer.fromHex))
    |> toEqual(
         encode(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
           [|("symbol", "BTC"), ("multiplier", "50"), ("what", "100")|],
         ),
       )
  });

  test("should be able to encode from bytes correctly2", () => {
    expect(Some("0x03000000425443900100000000000064" |> JsBuffer.fromHex))
    |> toEqual(
         encode(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
           [|("symbol", "BTC"), ("multiplier", "400"), ("what", "100")|],
         ),
       )
  });

  test("should return None if invalid schema", () => {
    expect(None)
    |> toEqual(
         encode(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
           [|("symbol", "band"), ("nulti", "400"), ("what", "100")|],
         ),
       )
  });

  test("should return None if invalid value", () => {
    expect(None)
    |> toEqual(
         encode(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
           [|("symbol", "band"), ("multiplier", "400"), ("what", "100")|],
         ),
       )
  });
});

describe("Expect Borsh to decode correctly", () => {
  test("should be able to decode from bytes correctly", () => {
    expect(Some([|("symbol", "BTC"), ("multiplier", "50"), ("what", "100")|]))
    |> toEqual(
         decode(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
           "0x03000000425443320000000000000064" |> JsBuffer.fromHex,
         ),
       )
  });

  test("should be able to decode from bytes correctly2", () => {
    expect(Some([|("symbol", "band"), ("multiplier", "400"), ("what", "100")|]))
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
