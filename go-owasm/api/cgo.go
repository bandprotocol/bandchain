package api

// #include "bindings.h"
//
// GoResult cGetCalldata(env_t *e, Span *calldata);
// GoResult cGetCalldata_cgo(env_t *e, Span *calldata) { return cGetCalldata(e, calldata); }
// GoResult cSetReturnData(env_t *e, Span data);
// GoResult cSetReturnData_cgo(env_t *e, Span data) { return cSetReturnData(e, data); }
// int64_t cGetAskCount(env_t *e);
// int64_t cGetAskCount_cgo(env_t *e) { return cGetAskCount(e); }
// int64_t cGetMinCount(env_t *e);
// int64_t cGetMinCount_cgo(env_t *e) { return cGetMinCount(e); }
// GoResult cGetAnsCount(env_t *e, int64_t *val);
// GoResult cGetAnsCount_cgo(env_t *e, int64_t *val) { return cGetAnsCount(e, val); }
// GoResult cAskExternalData(env_t *e, int64_t eid, int64_t did, Span data);
// GoResult cAskExternalData_cgo(env_t *e, int64_t eid, int64_t did, Span data) { return cAskExternalData(e, eid, did, data); }
// GoResult cGetExternalDataStatus(env_t *e, int64_t eid, int64_t vid, int64_t *status);
// GoResult cGetExternalDataStatus_cgo(env_t *e, int64_t eid, int64_t vid, int64_t *status) { return cGetExternalDataStatus(e, eid, vid, status); }
// GoResult cGetExternalData(env_t *e, int64_t eid, int64_t vid, Span *data);
// GoResult cGetExternalData_cgo(env_t *e, int64_t eid, int64_t vid, Span *data) { return cGetExternalData(e, eid, vid, data); }
import "C"
