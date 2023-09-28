package parameters

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCheckFlags(t *testing.T) {
	tests := []struct {
		name        string
		flags       CmdFlags
		shouldError bool
	}{
		{
			name: "Only count flag",
			flags: CmdFlags{
				Count: true,
			},
			shouldError: false,
		},
		{
			name: "Only duplicates flag",
			flags: CmdFlags{
				Duplicates: true,
			},
			shouldError: false,
		},
		{
			name: "Only unique flag",
			flags: CmdFlags{
				Unique: true,
			},
			shouldError: false,
		},
		{
			name: "Count && Duplicates flag",
			flags: CmdFlags{
				Count:      true,
				Duplicates: true,
			},
			shouldError: true,
		},
		{
			name: "Count && Unique flag",
			flags: CmdFlags{
				Count:  true,
				Unique: true,
			},
			shouldError: true,
		},
		{
			name: "Duplicates && Unique flag",
			flags: CmdFlags{
				Duplicates: true,
				Unique:     true,
			},
			shouldError: true,
		},
		{
			name: "Count && Duplicates && Unique flag",
			flags: CmdFlags{
				Count:      true,
				Duplicates: true,
				Unique:     true,
			},
			shouldError: true,
		},
		{
			name: "Any other flags",
			flags: CmdFlags{
				IgnoreCase: true,
				Fields:     2,
			},
			shouldError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckFlags(tt.flags)
			if tt.shouldError {
				require.Error(t, err, "Expected an error, but got nil")
			} else {
				require.NoError(t, err, "Expected no error, but got an error")
			}
		})
	}
}
