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

type enumValue[T any] struct {
	enum       *Enum[T]
	name       string
	ordinal    int
	properties []any
}

type EnumValue[T any] struct {
	value *enumValue[T]
}

func NewEnumValue[T any](name string, properties ...any) EnumValue[T] {
	return EnumValue[T]{
		value: &enumValue[T]{
			name:       name,
			ordinal:    -1,
			properties: properties,
		},
	}
}

func (this EnumValue[T]) Name() string {
	return this.value.name
}

func (this EnumValue[T]) Ordinal() int {
	return this.value.ordinal
}

func (this EnumValue[T]) Enum() *Enum[T] {
	return this.value.enum
}

func (this EnumValue[T]) Property(index int) any {
	return this.value.properties[index]
}

func (this EnumValue[T]) String() string {
	return this.value.name
}

func (this EnumValue[T]) setOrdinal(ordinal int) {
	this.value.ordinal = ordinal
}

func (this EnumValue[T]) setDescriptor(descriptor *Enum[T]) {
	this.value.enum = descriptor
}
