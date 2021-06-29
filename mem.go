//the reference: https://github.com/DmitriyVTitov/size/blob/master/size.go
package lru_cache

import (
	"reflect"
)

var (
	realtimeMem int64
)

func (l *LRUCache) exceedMaxMem() bool {
	cache := make(map[uintptr]bool)
	realtimeMem = sizeOf(reflect.Indirect(reflect.ValueOf(l)), cache)
	return realtimeMem > l.MaxM
}

// sizeOf returns the number of bytes the actual data represented by v occupies in memory.
// If there is an error, sizeOf returns -1.
func sizeOf(v reflect.Value, cache map[uintptr]bool) int64 {
	switch v.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		// return 0 if this node has been visited already (infinite recursion)
		if v.Kind() != reflect.Array && cache[v.Pointer()] {
			return 0
		}
		if v.Kind() != reflect.Array {
			cache[v.Pointer()] = true
		}
		sum := int64(0)
		for i := 0; i < v.Len(); i++ {
			s := sizeOf(v.Index(i), cache)
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum + int64(v.Type().Size())
	case reflect.Struct:
		sum := int64(0)
		for i, n := 0, v.NumField(); i < n; i++ {
			s := sizeOf(v.Field(i), cache)
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum
	case reflect.String:
		return int64(len(v.String()) + int(v.Type().Size()))
	case reflect.Ptr:
		// return Ptr size if this node has been visited already (infinite recursion)
		if cache[v.Pointer()] {
			return int64(v.Type().Size())
		}
		cache[v.Pointer()] = true
		if v.IsNil() {
			return int64(reflect.New(v.Type()).Type().Size())
		}
		s := sizeOf(reflect.Indirect(v), cache)
		if s < 0 {
			return -1
		}
		return s + int64(v.Type().Size())
	case reflect.Bool,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Int,
		reflect.Chan,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return int64(v.Type().Size())
	case reflect.Map:
		// return 0 if this node has been visited already (infinite recursion)
		if cache[v.Pointer()] {
			return 0
		}
		cache[v.Pointer()] = true
		sum := int64(0)
		keys := v.MapKeys()
		for i := range keys {
			val := v.MapIndex(keys[i])
			// calculate size of key and value separately
			sv := sizeOf(val, cache)
			if sv < 0 {
				return -1
			}
			sum += sv
			sk := sizeOf(keys[i], cache)
			if sk < 0 {
				return -1
			}
			sum += sk
		}
		// Include overhead due to unused map buckets.  10.79 comes
		// from https://golang.org/src/runtime/map.go.
		return sum + int64(v.Type().Size()) + int64(float64(len(keys))*10.79)
	case reflect.Interface:
		return sizeOf(v.Elem(), cache) + int64(v.Type().Size())
	}
	return -1
}
