package util_test

import (
	"testing"

	"github.com/go-jang/go/util"
	"github.com/stretchr/testify/require"
)

var updates = make([]string, 0)

func TestMessageReadWrite(t *testing.T) {
	t.Run("subsribtion test'", func(t *testing.T) {
		observable := util.NewObservable[string]()
		observable.AddObserver(onMessage)
		observable.NotifyObservers("str1")
		observable.NotifyObservers("str2")
		observable.NotifyObservers("str3")

		require.Equal(t, []string{"str1", "str2", "str3"}, updates)

		observable.RemoveObserver(onMessage)
		observable.NotifyObservers("str4")

		require.Equal(t, []string{"str1", "str2", "str3"}, updates)
	})
}

func onMessage(observable *util.Observable[string], arg string) {
	updates = append(updates, arg)
}
