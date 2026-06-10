package notifire

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseRoutes(t *testing.T) {
	cases := []struct {
		name    string
		raw     string
		want    []RouteSt
		wantErr bool
	}{
		{
			name: "empty",
			raw:  "",
			want: nil,
		},
		{
			name: "whitespace only",
			raw:  "   \n  \t ",
			want: nil,
		},
		{
			name: "single route",
			raw:  "- channel: sms\n  ttl: 60\n",
			want: []RouteSt{{Channel: "sms", Ttl: 60}},
		},
		{
			name: "multiple routes",
			raw:  "- channel: sms\n  ttl: 60\n- channel: whatsapp\n  ttl: 120\n",
			want: []RouteSt{
				{Channel: "sms", Ttl: 60},
				{Channel: "whatsapp", Ttl: 120},
			},
		},
		{
			name: "route without ttl",
			raw:  "- channel: sms\n",
			want: []RouteSt{{Channel: "sms"}},
		},
		{
			name: "flow style",
			raw:  "[{channel: sms, ttl: 30}, {channel: whatsapp}]",
			want: []RouteSt{
				{Channel: "sms", Ttl: 30},
				{Channel: "whatsapp"},
			},
		},
		{
			name:    "invalid yaml",
			raw:     "- channel: sms\n ttl: : oops",
			wantErr: true,
		},
		{
			name:    "wrong type for ttl",
			raw:     "- channel: sms\n  ttl: abc",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseRoutes(tc.raw)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}
