package parser

import (
	"github.com/sermojohn/postgres-client/pkg/config"
	"github.com/sermojohn/postgres-client/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"reflect"
	"testing"
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
		data, err := p.FetchIXPData(nonce, opt)
		assert.NoError(t, err)
		assert.True(t, len(nonce) > 0, "expected non empty nonce")
		t.Logf("params:%+v, summary: %+v", opt, data)
	}
}

func Test_filterServers(t *testing.T) {
	t.Parallel()

	type args struct {
		servers      IXPServerOptions
		clientConfig config.ClientConfig
	}
	tests := []struct {
		name string
		args args
		want []IXPServerOption
	}{
		{
			name: "happy case",
			args: args{
				servers: IXPServerOptions{
					{IXPServer: domain.IXPServer{IXP: "ixp-1", Country: "Greece"}},
					{IXPServer: domain.IXPServer{City: "Heraklion", Country: "Greece"}},
					{IXPServer: domain.IXPServer{Country: "Greece"}},
				},
				clientConfig: config.ClientConfig{
					Country:     "greece",
					ServerLimit: 3,
				},
			},
			want: []IXPServerOption{
				{IXPServer: domain.IXPServer{IXP: "ixp-1", Country: "Greece"}},
				{IXPServer: domain.IXPServer{City: "Heraklion", Country: "Greece"}},
				{IXPServer: domain.IXPServer{Country: "Greece"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterServers(tt.args.servers, tt.args.clientConfig); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterServers() = %v, want %v", got, tt.want)
			}
		})
	}
}
