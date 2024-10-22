package calculator

import "testing"

func TestEval(t *testing.T) {

	tests := []struct {
		expr    string
		want    float64
		wantErr bool
	}{
		{"1.5  + 2", 3.5, false},
		{"1+", 0, true},
		{"+", 0, true},
		{"2*3.5", 7, false},
		{"2*", 0, true},
		{"*3", 0, true},
		{"*", 0, true},
		{"2*3", 6, false},
		{"1-2", -1, false},
		{"1-", 0, true},
		{"-", 0, true},
		{"1/2", 0.5, false},
		{"1/", 0, true},
		{"/3", 0, true},
		{"/", 0, true},
		{"2.0", 2.0, false},
		{"-2.5", -2.5, false},
		{"   +2", 2, false},
		{"2..0", 0, true},
		{"2.", 2, false},
		{".2", .2, false},
		{" ", 0, true},
		{"1+2+3", 6, false},
		{"2*3*4", 24, false},
		{"2*3+4", 10, false},
		{"4-2-1", 1, false},
		{"100/4/5", 5, false},
		{"100/4*8", 200, false},
		{"1+2 -7", -4, false},
		{"1+2*3", 7, false},
		{"(1+2)*3", 9, false},
		{"1+-2", -1, false},
		{"2^2", 4, false},
		{"2^(2^3)", 256, false},
		{"2^2^3", 64, false},
		{"2x2-2^5/8", 0, false},
		{"--2", 0, true},
		{"1x-2", -2, false},
		{"()", 0, true},
		{"(123)", 123, false},
		{"(123", 0, true},
		{"~1", 0, true},
		{"1)", 0, true},
		{"1+", 0, true},
		{"1^(*)", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			got, err := Calculate(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eval(%s) error = %v, wantErr %v", tt.expr, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Eval(%s) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}
