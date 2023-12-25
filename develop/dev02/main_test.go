package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpacked(t *testing.T) {
	type test struct {
		name    string
		inStr   string
		outStr  string
		wantErr error
	}

	cases := []test{
		{
			name:   "success",
			inStr:  "a4bc2d5e",
			outStr: "aaaabccddddde",
		},
		{
			name:    "invalid line",
			inStr:   "1a4bc2d5e",
			outStr:  "",
			wantErr: ErrInvalidLine,
		},
		{
			name:   "successEscape",
			inStr:  `qwe\\5`,
			outStr: `qwe\\\\\`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			str, err := Unpacked(c.inStr)
			if c.wantErr != nil {
				require.ErrorIs(t, err, c.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, str, c.outStr)
		})
	}
}
