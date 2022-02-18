package main

import "testing"

func Test_unification(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
	}{
		{
			name:       "1",
			configPath: "config.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unification(tt.configPath)
		})
	}
}
