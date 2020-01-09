open Jest;
open Sha256;

describe("Expect Sha256 to work correctly", () => {
  Expect.(
    test("should be able to hash int array in to int array format", () =>
      expect("F23391B5DBF982E37FB7DADEA64AAE21CAE4C172"->JsBuffer.hexToArray->digest)
      |> toEqual([|
           225,
           204,
           133,
           245,
           148,
           141,
           52,
           200,
           118,
           206,
           234,
           166,
           14,
           0,
           79,
           109,
           96,
           243,
           68,
           190,
           59,
           36,
           214,
           129,
           178,
           190,
           249,
           18,
           145,
           152,
           182,
           11,
         |])
    )
  );

  Expect.(
    test("should be able to hash int array in to int array format", () =>
      expect("F23391B5DBF982E37FB7DADEA64AAE21CAE4C172"->JsBuffer.hexToArray->hexDigest)
      |> toBe("e1cc85f5948d34c876ceeaa60e004f6d60f344be3b24d681b2bef9129198b60b")
    )
  );
});
