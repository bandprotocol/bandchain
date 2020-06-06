open Jest;
open Obi;
open Expect;

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
