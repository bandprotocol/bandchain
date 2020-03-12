open Jest;
open ExecutableParser;

open Expect;

describe("Expect Parser to work correctly", () => {
  test("test getElement in list", () =>
    expect(["hello", "world"] |> getElementInList(_, 1)) |> toEqual("world")
  );
  test("test check regex success", () =>
    expect("benz=$1" |> checker) |> toEqual(true)
  );
  test("test check regex fail", () =>
    expect("bun=!1" |> checker) |> toEqual(false)
  );
  test("test check splitEqual success", () =>
    expect("benz=$1" |> splitToPair) |> toEqual(("benz", 1))
  );
  test("test check splitEqual variables value more than max_int", () =>
    expect("benz=$1333333333333333333" |> splitToPair) |> toEqual(("benz", (-1)))
  );

  test("test get 1 variable", () =>
    expect(
      {f|#!/bin/sh

symbol=$1

# Cryptocurrency price endpoint: https://www.coingecko.com/api/documentations/v3
url=\"https://api.coingecko.com/api/v3/simple/price?ids=$symbol&vs_currencies=usd\"

# Performs data fetching and parses the result
curl -s -X GET $url -H \"accept: application/json\" | jq -er ".[\"$symbol\"].usd\"
|f}
      |> getVariables
      |> Belt_Option.getExn,
    )
    |> toEqual(["symbol"])
  );

  test("test get 1 variable", () =>
    expect(
      {f|"#!/bin/bash

calldata=$1

echo $calldata

"|f} |> getVariables |> Belt_Option.getExn,
    )
    |> toEqual(["calldata"])
  );

  test("test get many variables", () =>
    expect(
      {f|"#!/bin/bash
symbol=$2
calldata=$1

echo $calldata
"|f}
      |> getVariables
      |> Belt_Option.getExn,
    )
    |> toEqual(["calldata", "symbol"])
  );

  test("test get many variables", () =>
    expect({f|"#!/bin/bash
symbol=$23
calldata=$1

echo $calldata
"|f} |> getVariables)
    |> toEqual(None)
  );
  test("test no variablie in script", () =>
    expect({f|"#!/bin/bash

echo $calldata
"|f} |> getVariables) |> toEqual(None)
  );
  test("test parseExecutableScript", () =>
    expect(
      "IyEvYmluL3NoCgpzeW1ib2w9JDEKCiMgQ3J5cHRvY3VycmVuY3kgcHJpY2UgZW5kcG9pbnQ6IGh0dHBzOi8vd3d3LmNvaW5nZWNrby5jb20vYXBpL2RvY3VtZW50YXRpb25zL3YzCnVybD0iaHR0cHM6Ly9hcGkuY29pbmdlY2tvLmNvbS9hcGkvdjMvc2ltcGxlL3ByaWNlP2lkcz0kc3ltYm9sJnZzX2N1cnJlbmNpZXM9dXNkIgoKIyBQZXJmb3JtcyBkYXRhIGZldGNoaW5nIGFuZCBwYXJzZXMgdGhlIHJlc3VsdApjdXJsIC1zIC1YIEdFVCAkdXJsIC1IICJhY2NlcHQ6IGFwcGxpY2F0aW9uL2pzb24iIHwganEgLWVyICIuW1wiJHN5bWJvbFwiXS51c2QiCg"
      |> JsBuffer.fromBase64
      |> parseExecutableScript
      |> Belt_Option.getExn,
    )
    |> toEqual(["symbol"])
  );

  test("test parseExecutableScript return None", () =>
    expect(
      "IyEvYmluL3NoCgpzeW1ib2w9JDEyCgojIENyeXB0b2N1cnJlbmN5IHByaWNlIGVuZHBvaW50OiBodHRwczovL3d3dy5jb2luZ2Vja28uY29tL2FwaS9kb2N1bWVudGF0aW9ucy92Mwp1cmw9Imh0dHBzOi8vYXBpLmNvaW5nZWNrby5jb20vYXBpL3YzL3NpbXBsZS9wcmljZT9pZHM9JHN5bWJvbCZ2c19jdXJyZW5jaWVzPXVzZCIKCiMgUGVyZm9ybXMgZGF0YSBmZXRjaGluZyBhbmQgcGFyc2VzIHRoZSByZXN1bHQKY3VybCAtcyAtWCBHRVQgJHVybCAtSCAiYWNjZXB0OiBhcHBsaWNhdGlvbi9qc29uIiB8IGpxIC1lciAiLltcIiRzeW1ib2xcIl0udXNkIgo="
      |> JsBuffer.fromBase64
      |> parseExecutableScript,
    )
    |> toEqual(None)
  );
});
