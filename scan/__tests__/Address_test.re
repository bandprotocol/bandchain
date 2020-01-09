open Jest;
open Address;

describe("Expect Address to work correctly", () => {
  Expect.(
    test("should be able create address from hex", () =>
      expect("F23391B5DBF982E37FB7DADEA64AAE21CAE4C172"->fromHex)
      |> toEqual(Address("f23391b5dbf982e37fb7dadea64aae21cae4c172"))
    )
  );

  Expect.(
    test("should be able create address from hex with 0x prefix", () =>
      expect("0xF23391B5DBF982E37FB7DADEA64AAE21CAE4C172"->fromHex)
      |> toEqual(Address("f23391b5dbf982e37fb7dadea64aae21cae4c172"))
    )
  );

  Expect.(
    test("should be able create address from fromBech32", () =>
      expect("bandvaloper13zmknvkq2sj920spz90g4r9zjan8g58423y76e"->fromBech32)
      |> toEqual(Address("88b769b2c05424553e01115e8a8ca297667450f5"))
    )
  );

  Expect.(
    test("should be able to convert self to hex", () =>
      expect(Address("88b769b2c05424553e01115e8a8ca297667450f5")->toHex)
      |> toEqual("88b769b2c05424553e01115e8a8ca297667450f5")
    )
  );

  Expect.(
    test("should be able to convert self to toOperatorBech32", () =>
      expect(Address("88b769b2c05424553e01115e8a8ca297667450f5")->toOperatorBech32)
      |> toEqual("bandvaloper13zmknvkq2sj920spz90g4r9zjan8g58423y76e")
    )
  );

  Expect.(
    test("should be able to convert self to toBech32", () =>
      expect(Address("88b769b2c05424553e01115e8a8ca297667450f5")->toBech32)
      |> toEqual("band13zmknvkq2sj920spz90g4r9zjan8g584x8qalj")
    )
  );

  Expect.(
    test("should be able to convert toBech32 to hex directly", () =>
      expect("band13zmknvkq2sj920spz90g4r9zjan8g584x8qalj"->bech32ToHex)
      |> toEqual("88b769b2c05424553e01115e8a8ca297667450f5")
    )
  );

  Expect.(
    test("should be able to convert hex to bech32 directly", () =>
      expect("88b769b2c05424553e01115e8a8ca297667450f5"->hexToBech32)
      |> toEqual("band13zmknvkq2sj920spz90g4r9zjan8g584x8qalj")
    )
  );

  Expect.(
    test("should be able to convert hex to hexToOperatorBech32 directly", () =>
      expect("88b769b2c05424553e01115e8a8ca297667450f5"->hexToOperatorBech32)
      |> toEqual("bandvaloper13zmknvkq2sj920spz90g4r9zjan8g58423y76e")
    )
  );
});
