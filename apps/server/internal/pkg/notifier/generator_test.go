package notifier

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"strconv"
	"testing"
)

func TestBasicGenerator_VerifyLink(t *testing.T) {
	regex, err := regexp.Compile(".*/verify\\?id=[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
	require.NoError(t, err)
	for i := range 10 {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b := &BasicGenerator{
				HostDomain: fmt.Sprintf("%d.com", i),
			}
			link, _, err := b.VerifyLink()
			assert.NoError(t, err)

			t.Log(link)
			match := regex.Match([]byte(link))
			assert.True(t, match)
		})
	}
}

func TestBasicGenerator_UnsubLink(t *testing.T) {
	regex, err := regexp.Compile(".*/unsub\\?id=[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
	require.NoError(t, err)
	for i := range 10 {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b := &BasicGenerator{
				HostDomain: fmt.Sprintf("%d.com", i),
			}
			link, _, err := b.UnsubLink()
			assert.NoError(t, err)

			t.Log(link)
			match := regex.Match([]byte(link))
			assert.True(t, match)
		})
	}
}
