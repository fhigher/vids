package server

import (
	"context"
	"reflect"

	"github.com/filecoin-project/go-jsonrpc/auth"
	"golang.org/x/xerrors"
)

const (
	PermRead  auth.Permission = "read" // default
	PermWrite auth.Permission = "write"
	PermSign  auth.Permission = "sign"  
	PermAdmin auth.Permission = "admin" 
)

var ErrNotSupported = xerrors.New("method not supported")

var AllPermissions = []auth.Permission{PermRead, PermWrite, PermSign, PermAdmin}
var DefaultPerms = []auth.Permission{PermRead}

var _internalField = "Internal"

// Reference：https://github.com/filecoin-project/lotus/blob/master/api/proxy_util.go

// GetInternalStructs extracts all pointers to 'Internal' sub-structs from the provided pointer to a proxy struct
func GetInternalStructs(in interface{}) []interface{} {
	return getInternalStructs(reflect.ValueOf(in).Elem())
}

func getInternalStructs(rv reflect.Value) []interface{} {
	var out []interface{}

	internal := rv.FieldByName(_internalField)
	ii := internal.Addr().Interface()
	out = append(out, ii)

	for i := 0; i < rv.NumField(); i++ {
		if rv.Type().Field(i).Name == _internalField {
			continue
		}

		sub := getInternalStructs(rv.Field(i))

		out = append(out, sub...)
	}

	return out
}

func permissionedProxies(in, out interface{}) {
	outs := GetInternalStructs(out)
	for _, o := range outs {
		log.Debugf("api: %+v", o)
		auth.PermissionedProxy(AllPermissions, DefaultPerms, in, o)
	}
}

func PermissionedAPI(api VidsServerAPI) VidsServerAPI {
	var out VidsServerStruct
	permissionedProxies(api, &out)
	log.Debugf("version: %+v", out.Version(context.Background()))
	return &out
}

type VidsServerStruct struct {
	Internal struct {
		Version func(p0 context.Context) string `perm:"read"`
	}
}

func (v *VidsServerStruct) Version(ctx context.Context) string {
	if v.Internal.Version == nil {
		return ""
	}

	return v.Internal.Version(ctx)
}

var _ VidsServerAPI = new(VidsServerStruct)
