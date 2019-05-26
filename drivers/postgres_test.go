package drivers

import "testing"

func TestPostgresEncode(t *testing.T) {
	tests := []struct{
		in string
		out string
	}{
		{
			"it'smypass",
			`'it\'smypass'`,
		},
		{
			"itsmy pass",
			`'itsmy\ pass'`,
		},
		{
			`its\mine`,
			`'its\\mine'`,
		},
		{
			`it's mine t\Own`,
			`'it\'s\ mine\ t\\Own'`,
		},
		{
			`C:\Users\ghost\certs and keys\key.file`,
			`'C:\\Users\\ghost\\certs\ and\ keys\\key.file'`,
		},
		{
			"pa%ss",
			`pa%ss`,
		},
		{
			"notExactlyAnAmazingPass",
			"notExactlyAnAmazingPass",
		},
		{
			" ",
			`'\ '`,
		},
		{
			"",
			"",
		},
	}

	var p postgres

	for _, tc := range tests {
		out := p.Encode(tc.in)

		if out != tc.out {
			t.Errorf("expected %s but got %s", tc.out, out)
		}
	}
}
