package enum_test

import (
	"testing"

	"github.com/go-jang/go/lang/enum"
	"github.com/stretchr/testify/require"
)

type Color struct {
	enum.EnumValue[*Color]
}

func NewColor(name string) *Color {
	return &Color{EnumValue: *enum.NewEnumValue[*Color](name)}
}

var RGB = enum.NewEnum(
	NewColor("Red"),
	NewColor("Green"),
	NewColor("Blue"),
)

var RGB2 = enum.NewEnum(
	NewColor("Red"),
	NewColor("Green"),
	NewColor("Blue"),
)

var NotRGB = enum.NewEnum(
	NewColor("Red"),
	NewColor("Green"),
	NewColor("Blue"),
	NewColor("Black"),
)

func TestEnumDescriptor_ValueOfOrdinal(t *testing.T) {
	t.Run("valid ordinals", func(t *testing.T) {
		require.Equal(t, "Red", RGB.ValueOfOrdinal(0).Name())
		require.Equal(t, "Green", RGB.ValueOfOrdinal(1).Name())
		require.Equal(t, "Blue", RGB.ValueOfOrdinal(2).Name())
	})

	t.Run("same instance returned", func(t *testing.T) {
		require.Same(t, RGB.ValueOfOrdinal(0), RGB.ValueOf("Red"))
		require.Same(t, RGB.ValueOfOrdinal(1), RGB.ValueOf("Green"))
	})

	t.Run("float discrete values", func(t *testing.T) {
		require.Equal(t, "Red", RGB.ValueOfOrdinal(0).Name())
		require.Equal(t, "Green", RGB.ValueOfOrdinal(1).Name())
	})

	t.Run("invalid ordinals", func(t *testing.T) {
		require.Panics(t, func() { RGB.ValueOfOrdinal(-1) })
		require.Panics(t, func() { RGB.ValueOfOrdinal(3) })
		require.Panics(t, func() { RGB.ValueOfOrdinal(100) })
	})

	t.Run("ordinal consistency", func(t *testing.T) {
		blue := RGB.ValueOfOrdinal(2)
		require.Equal(t, 2, blue.Ordinal())
		require.Equal(t, "Blue", blue.Name())

		red := blue.Enum().ValueOf("Red")
		require.Equal(t, "Red", red.Name())
		require.Same(t, RGB.ValueOf("Red"), red)
	})

	t.Run("different types with the save value", func(t *testing.T) {
		// same universe
		require.Equal(t, RGB.ValueOf("Red"), RGB2.ValueOf("Red"))
		require.Equal(t, RGB.ValueOf("Red").Name(), RGB2.ValueOf("Red").Name())
		require.Equal(t, RGB.ValueOf("Red").Ordinal(), RGB2.ValueOf("Red").Ordinal())

		// different universe
		require.NotEqual(t, RGB.ValueOf("Red"), NotRGB.ValueOf("Red"))
		require.Equal(t, RGB.ValueOf("Red").Name(), NotRGB.ValueOf("Red").Name())
		require.Equal(t, RGB.ValueOf("Red").Ordinal(), NotRGB.ValueOf("Red").Ordinal())
	})
}
