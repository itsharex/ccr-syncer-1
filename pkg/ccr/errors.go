package ccr

import "github.com/selectdb/ccr_syncer/pkg/xerror"

var (
	errBackendNotFound = xerror.XNew(xerror.Meta, "backend not found")
)