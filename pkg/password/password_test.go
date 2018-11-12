package password

import "testing"

func TestGenerate(t *testing.T) {
	tests := []struct {
		name string
		len  int
		want string
	}{
		{
			name: "simple",
			len:  10,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Generate(tt.len)
			if len(got) != tt.len {
				t.Errorf("Generate() = %d, want %d", len(got), tt.len)
			}
			got2 := Generate(tt.len)
			if got == got2 {
				t.Errorf("expected different password, but got the same")
			}
		})
	}
}
