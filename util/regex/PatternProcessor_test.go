package regex_test

import (
	"testing"

	"github.com/go-jang/go/util/regex"
	"github.com/stretchr/testify/require"
)

func Test_PatternProcessor_Process(t *testing.T) {
	processor := regex.PatternProcessorOf("\\w+")
	processor.OverrideResolve(func(match *regex.Match,
		super func(*regex.Match) any) any {

		switch match.Expr() {
		case "World":
			return "Mike"
		default:
			return match.Expr()
		}
	})

	require.Equal(t, " Hello Mike! ", processor.Process(" Hello World! "))
}

func TestPatternProcessorRecursive(test *testing.T) {
	values := map[string]string{
		"a": "${b}",
		"b": "${c}",
		"c": "resolved",
	}

	processor := regex.PatternProcessorOf(`\$\{(?P<key>[^}]+)\}`)
	processor.OverrideResolve(func(match *regex.Match, super func(*regex.Match) any) any {
		return values[match.NamedGroup("key").Value()]
	})

	require.Equal(test, "resolved", processor.Process("${a}"))
}

func TestPatternProcessorRecursiveCycle(test *testing.T) {
	values := map[string]string{
		"a": "${b}",
		"b": "${a}",
	}

	processor := regex.PatternProcessorOf(`\$\{(?P<key>[^}]+)\}`)
	processor.OverrideResolve(func(match *regex.Match, super func(*regex.Match) any) any {
		return values[match.NamedGroup("key").Value()]
	})

	require.Panics(test, func() {
		processor.Process("${a}")
	})
}
