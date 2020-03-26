open Jest;
open Borsh;
open Expect;

describe("Expect Borsh to decode correctly", () => {
  test("should be able to decode from bytes correctly", () => {
    expect(Some([|("symbol", "BTC"), ("multiplier", "50"), ("what", "100")|]))
    |> toEqual(
         Borsh.decode(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
           "0x03000000425443320000000000000064" |> JsBuffer.fromHex,
         ),
       )
  });

  test("should be able to decode from bytes correctly2", () => {
    expect(Some([|("symbol", "band"), ("multiplier", "400"), ("what", "100")|]))
    |> toEqual(
         Borsh.decode(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
           "0x0400000062616e64900100000000000064" |> JsBuffer.fromHex,
         ),
       )
  });

  test("should return None if invalid schema", () => {
    expect(None)
    |> toEqual(
         Borsh.decode(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input2",
           "0x03000000425443320000000000000064" |> JsBuffer.fromHex,
         ),
       )
  });

  test("should return None if invalid bytes", () => {
    expect(None)
    |> toEqual(
         Borsh.decode(
           {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
           "Input",
           "0x030000004254433200000000000064" |> JsBuffer.fromHex,
         ),
       )
  });
});
