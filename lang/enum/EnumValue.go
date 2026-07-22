package enum

func ValueOf[T any](d *Enum[T], name string) T {
	return d.ValueOf(name)
}

func ValueOfOrdinal[T any](d *Enum[T], ordinal int) T {
	return d.ValueOfOrdinal(ordinal)
}

func ValuesOf[T any](d *Enum[T]) []T {
	return d.Values()
}

type EnumValue[T any] struct {
	enum    *Enum[T]
	name    string
	ordinal int
}

func NewEnumValue[T any](name string) *EnumValue[T] {
	return &EnumValue[T]{name: name, ordinal: -1}
}

func (this *EnumValue[T]) Name() string {
	return this.name
}

func (this *EnumValue[T]) Ordinal() int {
	return this.ordinal
}

func (this *EnumValue[T]) Enum() *Enum[T] {
	return this.enum
}

func (this *EnumValue[T]) String() string {
	return this.name
}

func (this *EnumValue[T]) setOrdinal(ordinal int) {
	this.ordinal = ordinal
}

func (this *EnumValue[T]) setDescriptor(descriptor *Enum[T]) {
	this.enum = descriptor
}
