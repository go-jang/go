package enum_test

import (
	"encoding/json"
	"testing"

	"github.com/go-jang/go/lang/enum"
	"github.com/stretchr/testify/require"
)

type Color struct {
	enum.EnumValue[Color]
}

func NewColor(name, description string) Color {
	return Color{EnumValue: enum.NewEnumValue[Color](name, description)}
}

func (this Color) Description() string {
	return this.Property(0).(string)
}

func (this Color) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.Name())
}

func (this *Color) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}
	*this = RGB.ValueOf(name)
	return nil
}

var RGB = enum.NewEnum(
	NewColor("Red", "Description of red"),
	NewColor("Green", "Description of green"),
	NewColor("Blue", "Description of blue"),
)

var RGB2 = enum.NewEnum(
	NewColor("Red", "Description of red"),
	NewColor("Green", "Description of green"),
	NewColor("Blue", "Description of blue"),
)

var NotRGB = enum.NewEnum(
	NewColor("Red", "Description of red"),
	NewColor("Green", "Description of green"),
	NewColor("Blue", "Description of blue"),
	NewColor("Black", "Description of black"),
)

func TestEnumDescriptor_ValueOfOrdinal(t *testing.T) {
	t.Run("valid ordinals", func(t *testing.T) {
		require.Equal(t, "Red", RGB.ValueOfOrdinal(0).Name())
		require.Equal(t, "Green", RGB.ValueOfOrdinal(1).Name())
		require.Equal(t, "Blue", RGB.ValueOfOrdinal(2).Name())
	})

	t.Run("same instance returned", func(t *testing.T) {
		require.True(t, RGB.ValueOfOrdinal(0) == RGB.ValueOf("Red"))
		require.True(t, RGB.ValueOfOrdinal(1) == RGB.ValueOf("Green"))
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
		require.True(t, RGB.ValueOf("Red") == red)
	})

	t.Run("different types with the same value", func(t *testing.T) {
		require.Equal(t, RGB.ValueOf("Red"), RGB2.ValueOf("Red"))
		require.Equal(t, RGB.ValueOf("Red").Name(), RGB2.ValueOf("Red").Name())
		require.False(t, RGB.ValueOf("Red") == RGB2.ValueOf("Red"))

		require.NotEqual(t, RGB.ValueOf("Red"), NotRGB.ValueOf("Red"))
		require.Equal(t, RGB.ValueOf("Red").Name(), NotRGB.ValueOf("Red").Name())
		require.False(t, RGB.ValueOf("Red") == NotRGB.ValueOf("Red"))
	})
}

func TestEnumJSON(t *testing.T) {
	type Model struct {
		Text  string `json:"text"`
		Color Color  `json:"color"`
	}

	red := RGB.ValueOf("Red")
	green := RGB.ValueOf("Green")
	blue := RGB.ValueOf("Blue")

	source := Model{Text: "test", Color: red}
	data, err := json.Marshal(source)
	require.NoError(t, err)
	require.Equal(t, `{"text":"test","color":"Red"}`, string(data))
	require.Equal(t, `Red`, red.Name())
	require.Equal(t, 0, red.Ordinal())
	require.Equal(t, `Description of red`, red.Description())

	target := Model{Color: green}
	err = json.Unmarshal(data, &target)
	require.NoError(t, err)

	require.Equal(t, "test", target.Text)
	require.True(t, target.Color == red)
	require.Equal(t, "Red", target.Color.Name())
	require.Equal(t, 0, target.Color.Ordinal())
	require.Equal(t, "Description of red", target.Color.Description())

	require.True(t, RGB.ValueOf("Red") == red)
	require.True(t, RGB.ValueOf("Green") == green)
	require.True(t, RGB.ValueOf("Blue") == blue)

	require.Equal(t, "Red", RGB.ValueOf("Red").Name())
	require.Equal(t, "Green", RGB.ValueOf("Green").Name())
	require.Equal(t, "Blue", RGB.ValueOf("Blue").Name())

	require.Equal(t, 0, RGB.ValueOf("Red").Ordinal())
	require.Equal(t, 1, RGB.ValueOf("Green").Ordinal())
	require.Equal(t, 2, RGB.ValueOf("Blue").Ordinal())

	require.Equal(t, "Description of red", RGB.ValueOf("Red").Description())
	require.Equal(t, "Description of green", RGB.ValueOf("Green").Description())
	require.Equal(t, "Description of blue", RGB.ValueOf("Blue").Description())
}
