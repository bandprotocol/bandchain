package api

// #include "bindings.h"
//
// Error cGetCalldata(env_t *e, Span *calldata);
// Error cGetCalldata_cgo(env_t *e, Span *calldata) { return cGetCalldata(e, calldata); }
// Error cSetReturnData(env_t *e, Span data);
// Error cSetReturnData_cgo(env_t *e, Span data) { return cSetReturnData(e, data); }
// int64_t cGetAskCount(env_t *e);
// int64_t cGetAskCount_cgo(env_t *e) { return cGetAskCount(e); }
// int64_t cGetMinCount(env_t *e);
// int64_t cGetMinCount_cgo(env_t *e) { return cGetMinCount(e); }
// Error cGetAnsCount(env_t *e, int64_t *val);
// Error cGetAnsCount_cgo(env_t *e, int64_t *val) { return cGetAnsCount(e, val); }
// Error cAskExternalData(env_t *e, int64_t eid, int64_t did, Span data);
// Error cAskExternalData_cgo(env_t *e, int64_t eid, int64_t did, Span data) { return cAskExternalData(e, eid, did, data); }
// Error cGetExternalDataStatus(env_t *e, int64_t eid, int64_t vid, int64_t *status);
// Error cGetExternalDataStatus_cgo(env_t *e, int64_t eid, int64_t vid, int64_t *status) { return cGetExternalDataStatus(e, eid, vid, status); }
// Error cGetExternalData(env_t *e, int64_t eid, int64_t vid, Span *data);
// Error cGetExternalData_cgo(env_t *e, int64_t eid, int64_t vid, Span *data) { return cGetExternalData(e, eid, vid, data); }
import "C"
