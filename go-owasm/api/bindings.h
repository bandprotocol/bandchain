#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

enum Error {
  Ok = 0,
  CompileFail = 1,
  RunFail = 2,
  ParseFail = 3,
  Nul = 4,
  NonUtf8Result = 5,
};
typedef int32_t Error;

typedef struct {
  uint8_t *ptr;
  uintptr_t len;
  uintptr_t cap;
} Span;

typedef struct {
  uint8_t _private[0];
} env_t;

typedef struct {
  Span (*get_calldata)(env_t*);
  void (*set_return_data)(env_t*, Span data);
  int64_t (*get_ask_count)(env_t*);
  int64_t (*get_min_count)(env_t*);
  int64_t (*get_ans_count)(env_t*);
  void (*ask_external_data)(env_t*, int64_t eid, int64_t did, Span data);
  int64_t (*get_external_data_status)(env_t*, int64_t eid, int64_t vid);
  Span (*get_external_data)(env_t*, int64_t eid, int64_t vid);
} EnvDispatcher;

typedef struct {
  env_t *env;
  EnvDispatcher dis;
} Env;

Error do_compile(Span input, Span *output);

Error do_run(Span code, bool is_prepare, Env env);

Error do_wat2wasm(Span input, Span *output);
