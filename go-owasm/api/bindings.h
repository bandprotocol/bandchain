#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

enum Error {
  NoError = 0,
  CompliationError = 1,
  RunError = 2,
  ParseError = 3,
  WriteBinaryError = 4,
  ResolveNamesError = 5,
  ValidateError = 6,
  SpanExceededCapacityError = 7,
  DeserializationError = 8,
  GasCounterInjectionError = 9,
  SerializationError = 10,
  GasLimitExceedError = 11,
  NoMemoryWasmError = 12,
  MinimumMemoryExceedError = 13,
  SetMaximumMemoryError = 14,
  StackHeightInstrumentationError = 15,
  CheckWasmImportsError = 16,
  CheckWasmExportsError = 17,
  InvalidSignatureFunctionError = 18,
  FunctionNotFoundError = 11,
  GasLimitExceedError = 12,
  NoMemoryWasmError = 13,
  MinimumMemoryExceedError = 14,
  SetMaximumMemoryError = 15,
  StackHeightInstrumentationError = 16,
  UnknownError = 255,
};
typedef int32_t Error;

enum GoResult {
  Ok = 0,
  SetReturnDataWrongPeriod = 1,
  AnsCountWrongPeriod = 2,
  AskExternalDataWrongPeriod = 3,
  GetExternalDataStatusWrongPeriod = 4,
  GetExternalDataWrongPeriod = 5,
  GetExternalDataFromUnreportedValidator = 6,
  SpanExceededCapacity = 7,
  /**
   * An error happened during normal operation of a Go callback
   */
  Other = 8,
};
typedef int32_t GoResult;

typedef struct {
  uint8_t *ptr;
  uintptr_t len;
  uintptr_t cap;
} Span;

typedef struct {
  uint8_t _private[0];
} env_t;

typedef struct {
  Span (*get_calldata)(env_t *);
  GoResult (*set_return_data)(env_t *, Span data);
  int64_t (*get_ask_count)(env_t *);
  int64_t (*get_min_count)(env_t *);
  GoResult (*get_ans_count)(env_t *, int64_t *);
  GoResult (*ask_external_data)(env_t *, int64_t eid, int64_t did, Span data);
  GoResult (*get_external_data_status)(env_t *, int64_t eid, int64_t vid,
                                       int64_t *status);
  GoResult (*get_external_data)(env_t *, int64_t eid, int64_t vid, Span *data);
} EnvDispatcher;

typedef struct {
  env_t *env;
  EnvDispatcher dis;
} Env;

Error do_compile(Span input, Span *output);

Error do_run(Span code, uint32_t gas_limit, bool is_prepare, Env env);

Error do_wat2wasm(Span input, Span *output);
