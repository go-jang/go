package enum

import "github.com/go-jang/go/lang"

type value[T any] interface {
	Name() string
	setOrdinal(int)
	setDescriptor(*Enum[T])
}

type Enum[T any] struct {
	byName    map[string]T
	byOrdinal []T
}

func NewEnum[T value[T]](values ...T) *Enum[T] {
	d := &Enum[T]{
		byName:    make(map[string]T, len(values)),
		byOrdinal: make([]T, len(values)),
	}
	for i, value := range values {
		value.setOrdinal(i)
		value.setDescriptor(d)
		if _, exists := d.byName[value.Name()]; exists {
			panic("Duplicate enum constant: " + value.Name())
		}
		d.byName[value.Name()] = value
		d.byOrdinal[i] = value
	}
	return d
}

func (this *Enum[T]) ValueOf(name string) T {
	v, ok := this.byName[name]
	lang.Assert(ok, "No enum constant %T.%s", this, name)
	return v
}

func (this *Enum[T]) ValueOfOrdinal(ordinal int) T {
	lang.Assert(ordinal >= 0 && ordinal < len(this.byOrdinal), "Invalid ordinal %T: %d", this, ordinal)
	return this.byOrdinal[ordinal]
}

func (this *Enum[T]) Values() []T {
	return append([]T(nil), this.byOrdinal...)
}
