package cmd

import "testing"

func TestFormatTime(t *testing.T) {
	tests := []struct {
		name      string
		givenTz        string
		wantErr bool
	}{
		{"Valid timezone", "Asia/Kathmandu", false},
		{"Invalid timezone", "Invalid/Timezone", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, gotErr := formatTime(test.givenTz)
			if (gotErr != nil) != test.wantErr {
				t.Errorf("Expected error: %v, but got: %v", test.wantErr, gotErr)
			}
		})
	}
}
