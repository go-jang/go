package objects

type EqualsProvider interface {
	Equals(o EqualsProvider) bool
}
