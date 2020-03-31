open Jest;
open JsBuffer;
open Expect;

describe("Expect JsBuffer to work correctly", () => {
  test("should be able to convert upper case hex to lower case", () =>
    expect(fromHex("F23391B5DBF982E37FB7DADEA64AAE21CAE4C172") |> toHex)
    |> toBe("f23391b5dbf982e37fb7dadea64aae21cae4c172")
  );

  test("should be able to convert upper case hex with 0x to lower case without 0x", () =>
    expect(fromHex("0xF23391B5DBF982E37FB7DADEA64AAE21CAE4C172"))
    |> toEqual(fromHex("f23391b5dbf982e37fb7dadea64aae21cae4c172"))
  );

  test("should be able to convert upper case hex with 0X to lower case without 0X", () =>
    expect(fromHex("0XF23391B5DBF982E37FB7DADEA64AAE21CAE4C172"))
    |> toEqual(fromHex("f23391b5dbf982e37fb7dadea64aae21cae4c172"))
  );

  test("should be able to get hexString with 0x prefix", () =>
    expect(fromHex("f23391b5dbf982e37fb7dadea64aae21cae4c172") |> toHex(~with0x=true))
    |> toBe("0xf23391b5dbf982e37fb7dadea64aae21cae4c172")
  );

  test("should be able to convert buffer to base64", () =>
    expect(
      fromHex("0x0235b618ab0f0e9f48b1af32f56b78d955c432279893714737a937035024b83c58") |> toBase64,
    )
    |> toBe("AjW2GKsPDp9Isa8y9Wt42VXEMieYk3FHN6k3A1AkuDxY")
  );

  test("should be able to convert from array", () =>
    expect(
      from([|
        242,
        51,
        145,
        181,
        219,
        249,
        130,
        227,
        127,
        183,
        218,
        222,
        166,
        74,
        174,
        33,
        202,
        228,
        193,
        114,
      |]),
    )
    |> toEqual(fromHex("f23391b5dbf982e37fb7dadea64aae21cae4c172"))
  );

  test("should be able to convert from base64", () =>
    expect(fromBase64("AjW2GKsPDp9Isa8y9Wt42VXEMieYk3FHN6k3A1AkuDxY"))
    |> toEqual(fromHex("0235b618ab0f0e9f48b1af32f56b78d955c432279893714737a937035024b83c58"))
  );

  test("should be able to convert base64 to hex directly", () =>
    expect("AjW2GKsPDp9Isa8y9Wt42VXEMieYk3FHN6k3A1AkuDxY" |> base64ToHex)
    |> toBe("0235b618ab0f0e9f48b1af32f56b78d955c432279893714737a937035024b83c58")
  );

  test("should be able to convert hex to base64 directly", () =>
    expect("0235b618ab0f0e9f48b1af32f56b78d955c432279893714737a937035024b83c58" |> hexToBase64)
    |> toBe("AjW2GKsPDp9Isa8y9Wt42VXEMieYk3FHN6k3A1AkuDxY")
  );

  test("should be able to convert int array to hex directly", () =>
    expect(
      [|
        242,
        51,
        145,
        181,
        219,
        249,
        130,
        227,
        127,
        183,
        218,
        222,
        166,
        74,
        174,
        33,
        202,
        228,
        193,
        114,
      |]
      |> arrayToHex,
    )
    |> toBe("f23391b5dbf982e37fb7dadea64aae21cae4c172")
  );

  test("should be able to convert hex to int array directly", () =>
    expect("f23391b5dbf982e37fb7dadea64aae21cae4c172" |> hexToArray)
    |> toEqual([|
         242,
         51,
         145,
         181,
         219,
         249,
         130,
         227,
         127,
         183,
         218,
         222,
         166,
         74,
         174,
         33,
         202,
         228,
         193,
         114,
       |])
  );

  test("should be able to convert int array to base64 directly", () =>
    expect(
      [|
        2,
        53,
        182,
        24,
        171,
        15,
        14,
        159,
        72,
        177,
        175,
        50,
        245,
        107,
        120,
        217,
        85,
        196,
        50,
        39,
        152,
        147,
        113,
        71,
        55,
        169,
        55,
        3,
        80,
        36,
        184,
        60,
        88,
      |]
      |> arrayToBase64,
    )
    |> toBe("AjW2GKsPDp9Isa8y9Wt42VXEMieYk3FHN6k3A1AkuDxY")
  );

  test("should be able to convert base64 to int array directly", () =>
    expect("AjW2GKsPDp9Isa8y9Wt42VXEMieYk3FHN6k3A1AkuDxY" |> base64ToArray)
    |> toEqual([|
         2,
         53,
         182,
         24,
         171,
         15,
         14,
         159,
         72,
         177,
         175,
         50,
         245,
         107,
         120,
         217,
         85,
         196,
         50,
         39,
         152,
         147,
         113,
         71,
         55,
         169,
         55,
         3,
         80,
         36,
         184,
         60,
         88,
       |])
  );
});
