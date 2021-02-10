package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestIxpParser_InitIXPOptions(t *testing.T) {
	t.Parallel()

	p := &ixpParser{hc: &http.Client{}}

	options, err := p.InitIXPServers()
	require.NoError(t, err)

	assert.True(t, len(options.Nonce) > 0, "expected non empty nonce")
	assert.True(t, len(options.Servers) > 0, "expected non empty options")
}

func TestIxpParser_FetchIXPData(t *testing.T) {
	t.Parallel()

	p := &ixpParser{hc: &http.Client{}}

	options, err := p.InitIXPServers()
	require.NoError(t, err)

	nonce := options.Nonce
	for _, opt := range options.Servers[0:5] {
		data, err := p.FetchIXPData(nonce, opt, false)
		assert.NoError(t, err)
		assert.True(t, len(nonce) > 0, "expected non empty nonce")
		t.Logf("params:%+v, summary: %+v", opt, data)
	}
}

func addDelay() {
	time.Sleep(1 * time.Second)
}
