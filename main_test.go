package main

import (
	"reflect"
	"testing"
)

func TestEvaluateRange_SingleNumber(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		step     float64
		expected []float64
	}{
		{
			name:     "single digit",
			expr:     "5",
			step:     1.0,
			expected: []float64{0, 1, 2, 3, 4},
		},
		{
			name:     "double digit",
			expr:     "10",
			step:     1.0,
			expected: []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name:     "single number with step",
			expr:     "3",
			step:     1.0,
			expected: []float64{0, 1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluateRange(tt.expr, tt.step)
			if err != nil {
				t.Fatalf("evaluateRange() error = %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("evaluateRange() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEvaluateRange_DotDotNotation(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		step     float64
		expected []float64
	}{
		{
			name:     "integer range",
			expr:     "2..4",
			step:     1.0,
			expected: []float64{2, 3, 4},
		},
		{
			name:     "float range",
			expr:     "1.5..3.5",
			step:     1.0,
			expected: []float64{1.5, 2.5, 3.5},
		},
		{
			name:     "range with custom step",
			expr:     "0..5",
			step:     0.5,
			expected: []float64{0, 0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5},
		},
		{
			name:     "descending range",
			expr:     "5..2",
			step:     1.0,
			expected: []float64{5, 4, 3, 2},
		},
		{
			name:     "single value range",
			expr:     "3..3",
			step:     1.0,
			expected: []float64{3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluateRange(tt.expr, tt.step)
			if err != nil {
				t.Fatalf("evaluateRange() error = %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("evaluateRange() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEvaluateRange_BracketNotation(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		step     float64
		expected []float64
	}{
		{
			name:     "both inclusive",
			expr:     "[2,5]",
			step:     1.0,
			expected: []float64{2, 3, 4, 5},
		},
		{
			name:     "both exclusive",
			expr:     "(2,5)",
			step:     1.0,
			expected: []float64{3, 4},
		},
		{
			name:     "start inclusive, end exclusive",
			expr:     "[2,5)",
			step:     1.0,
			expected: []float64{2, 3, 4},
		},
		{
			name:     "start exclusive, end inclusive",
			expr:     "(2,5]",
			step:     1.0,
			expected: []float64{3, 4, 5},
		},
		{
			name:     "float range both inclusive",
			expr:     "[2.5,4.5]",
			step:     1.0,
			expected: []float64{2.5, 3.5, 4.5},
		},
		{
			name:     "float range both exclusive",
			expr:     "(2.5,4.5)",
			step:     1.0,
			expected: []float64{3.5},
		},
		{
			name:     "float range with custom step",
			expr:     "[2.5,4.5]",
			step:     0.5,
			expected: []float64{2.5, 3, 3.5, 4, 4.5},
		},
		{
			name:     "exclusive range with custom step",
			expr:     "(2.5,4.5)",
			step:     0.5,
			expected: []float64{3, 3.5, 4},
		},
		{
			name:     "quoted bracket notation",
			expr:     "'[2.5,4.5]'",
			step:     1.0,
			expected: []float64{2.5, 3.5, 4.5},
		},
		{
			name:     "double quoted bracket notation",
			expr:     "\"[2.5,4.5]\"",
			step:     1.0,
			expected: []float64{2.5, 3.5, 4.5},
		},
		{
			name:     "descending bracket range",
			expr:     "[5,2]",
			step:     1.0,
			expected: []float64{5, 4, 3, 2},
		},
		{
			name:     "descending exclusive range",
			expr:     "(5,2)",
			step:     1.0,
			expected: []float64{4, 3},
		},
		{
			name:     "single value inclusive",
			expr:     "[3,3]",
			step:     1.0,
			expected: []float64{3},
		},
		{
			name:     "single value exclusive",
			expr:     "(3,3)",
			step:     1.0,
			expected: []float64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluateRange(tt.expr, tt.step)
			if err != nil {
				t.Fatalf("evaluateRange() error = %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("evaluateRange() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEvaluateRange_Errors(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		step    float64
		wantErr bool
	}{
		{
			name:    "invalid dot dot format",
			expr:    "2..3..4",
			step:    1.0,
			wantErr: true,
		},
		{
			name:    "invalid bracket format",
			expr:    "[2,3,4]",
			step:    1.0,
			wantErr: true,
		},
		{
			name:    "invalid number in dot dot",
			expr:    "abc..def",
			step:    1.0,
			wantErr: true,
		},
		{
			name:    "unrecognized format",
			expr:    "invalid",
			step:    1.0,
			wantErr: true,
		},
		{
			name:    "mismatched brackets",
			expr:    "[2,3)",
			step:    1.0,
			wantErr: false, // This should work, brackets can be mixed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := evaluateRange(tt.expr, tt.step)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateRange(t *testing.T) {
	tests := []struct {
		name           string
		start          float64
		end            float64
		step           float64
		startInclusive bool
		endInclusive   bool
		expected       []float64
	}{
		{
			name:           "ascending inclusive both",
			start:          1,
			end:            3,
			step:           1.0,
			startInclusive: true,
			endInclusive:   true,
			expected:       []float64{1, 2, 3},
		},
		{
			name:           "ascending exclusive both",
			start:          1,
			end:            3,
			step:           1.0,
			startInclusive: false,
			endInclusive:   false,
			expected:       []float64{2},
		},
		{
			name:           "ascending start inclusive end exclusive",
			start:          1,
			end:            3,
			step:           1.0,
			startInclusive: true,
			endInclusive:   false,
			expected:       []float64{1, 2},
		},
		{
			name:           "ascending start exclusive end inclusive",
			start:          1,
			end:            3,
			step:           1.0,
			startInclusive: false,
			endInclusive:   true,
			expected:       []float64{2, 3},
		},
		{
			name:           "descending inclusive both",
			start:          5,
			end:            2,
			step:           1.0,
			startInclusive: true,
			endInclusive:   true,
			expected:       []float64{5, 4, 3, 2},
		},
		{
			name:           "descending exclusive both",
			start:          5,
			end:            2,
			step:           1.0,
			startInclusive: false,
			endInclusive:   false,
			expected:       []float64{4, 3},
		},
		{
			name:           "zero step defaults to 1.0",
			start:          1,
			end:            3,
			step:           0,
			startInclusive: true,
			endInclusive:   true,
			expected:       []float64{1, 2, 3},
		},
		{
			name:           "small step",
			start:          0,
			end:            1,
			step:           0.25,
			startInclusive: true,
			endInclusive:   true,
			expected:       []float64{0, 0.25, 0.5, 0.75, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateRange(tt.start, tt.end, tt.step, tt.startInclusive, tt.endInclusive)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("generateRange() = %v, want %v", result, tt.expected)
			}
		})
	}
}
