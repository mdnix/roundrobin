package roundrobin_test

import (
	"reflect"
	"testing"

	"github.com/mdnix/roundrobin"
)

func TestNewService(t *testing.T) {

	var tests = []struct {
		addresses []string
		expected  []string
		err       bool
	}{
		{
			addresses: []string{
				"192.168.1.1:4444",
				"192.168.1.2:4444",
				"192.168.1.3:4444",
			},
			expected: []string{
				"192.168.1.1:4444",
				"192.168.1.2:4444",
				"192.168.1.3:4444",
				"192.168.1.1:4444",
			},
			err: false,
		},
		{
			addresses: []string{
				"192.168.1.1:4444",
				"192.168.1.1:4444",
				"192.168.1.2:4444",
			},
			expected: []string{},
			err:      true,
		},
		{
			addresses: []string{
				"192.168.1.1:4444",
				"192.168.1.2:4444",
				"192.168.1.293:4444",
			},
			expected: []string{},
			err:      true,
		},
		{
			addresses: []string{
				"192.168.1.1:67000",
				"192.168.1.2:67000",
				"192.168.1.3:67000",
			},
			expected: []string{},
			err:      true,
		},
		{
			addresses: []string{},
			expected:  []string{},
			err:       true,
		},
	}

	for i, test := range tests {
		service, err := roundrobin.NewService(test.addresses)

		if gotErr, wantErr := !(err == nil), test.err; gotErr != wantErr {
			t.Errorf("test[%d]: have: %v, expected: %v", i, gotErr, wantErr)
		}

		var have []string
		for i := 0; i < len(test.expected); i++ {
			have = append(have, service.Next().Address)
		}

		if !reflect.DeepEqual(have, test.expected) && !test.err {
			t.Errorf("test[%d]: have: %v, expected: %v", i, have, test.expected)
		}

	}
}
