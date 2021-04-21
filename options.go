package mindreminder

import "reflect"

type Option func(options *options)

type options struct {
	lazyUpdate       bool
	lazyUpdateFields []string
	metaTypes        map[string]reflect.Type
	objectTypes      map[string]reflect.Type
	// computeDiff      bool
}
