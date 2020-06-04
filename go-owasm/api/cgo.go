package api

// #include "bindings.h"
//
// Span cGetCalldata(env_t *e);
// Span cGetCalldata_cgo(env_t *e) { return cGetCalldata(e); }
// void cSetReturnData(env_t *e, Span data);
// void cSetReturnData_cgo(env_t *e, Span data) { return cSetReturnData(e, data); }
// int64_t cGetAskCount(env_t *e);
// int64_t cGetAskCount_cgo(env_t *e) { return cGetAskCount(e); }
// int64_t cGetMinCount(env_t *e);
// int64_t cGetMinCount_cgo(env_t *e) { return cGetMinCount(e); }
// int64_t cGetAnsCount(env_t *e);
// int64_t cGetAnsCount_cgo(env_t *e) { return cGetAnsCount(e); }
// void cAskExternalData(env_t *e, int64_t eid, int64_t did);
// void cAskExternalData_cgo(env_t *e, int64_t eid, int64_t did) { return cAskExternalData(e, eid, did); }
// int64_t cGetExternalDataStatus(env_t *e, int64_t eid, int64_t vid);
// int64_t cGetExternalDataStatus_cgo(env_t *e, int64_t eid, int64_t vid) { return cGetExternalDataStatus(e, eid, vid); }
// Span cGetExternalData(env_t *e, int64_t eid, int64_t vid);
// Span cGetExternalData_cgo(env_t *e, int64_t eid, int64_t vid) { return cGetExternalData(e, eid, vid); }
import "C"
