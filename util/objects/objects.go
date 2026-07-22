package objects

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

func ToString(o any) string {
	switch v := o.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", o)
	}
}

func Equal(o1, o2 any) bool {
	if o1 == nil && o2 == nil {
		return true
	} else if o1 == nil || o2 == nil {
		return false
	}
	t1 := reflect.TypeOf(o1)
	t2 := reflect.TypeOf(o2)
	if t1 != t2 {
		return false
	}

	if eq1, ok1 := o1.(EqualsProvider); ok1 {
		eq2, ok2 := o2.(EqualsProvider)
		return ok2 && eq1.Equals(eq2)
	}

	if s1, ok1 := o1.(string); ok1 {
		s2, ok2 := o2.(string)
		return ok2 && s1 == s2
	}

	// any/interface{} is two machine words: [type, data]
	pair1 := (*[2]uintptr)(unsafe.Pointer(&o1))
	pair2 := (*[2]uintptr)(unsafe.Pointer(&o2))
	data1 := pair1[1]
	data2 := pair2[1]

	ptr1 := unsafe.Pointer(data1)
	ptr2 := unsafe.Pointer(data2)

	b1 := unsafe.Slice((*byte)(ptr1), t1.Size())
	b2 := unsafe.Slice((*byte)(ptr2), t2.Size())

	return bytes.Equal(b1, b2)
}

func HashCode(o any) int {
	if o == nil {
		return 0
	} else if h, ok := o.(HashCodeProvider); ok {
		return h.HashCode()
	} else if s, ok := o.(string); ok {
		return int(Fnv1a64String(s))
	}
	t := reflect.TypeOf(o)
	size := t.Size()

	// any/interface{} is two machine words: [type, data]
	pair := (*[2]uintptr)(unsafe.Pointer(&o))
	data := pair[1]

	ptr := unsafe.Pointer(data)
	bytes := unsafe.Slice((*byte)(ptr), size)

	return int(Fnv1a64(bytes))
}

func Fnv1a64String(s string) uint64 {
	var h uint64 = 1469598103934665603
	const prime uint64 = 1099511628211

	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= prime
	}
	return h
}

func Fnv1a64(objects ...[]byte) uint64 {
	var h uint64 = 1469598103934665603
	const prime uint64 = 1099511628211

	for _, object := range objects {
		for _, b := range object {
			h ^= uint64(b)
			h *= prime
		}
	}
	return h
}

func Ptr[T any](v T) *T {
	return &v
}

func FuncPtr(f any) uintptr {
	return *(*uintptr)(unsafe.Pointer(&f))
}
