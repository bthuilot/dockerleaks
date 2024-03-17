package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReverse(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "empty slice",
			in:   []string{},
			want: []string{},
		},
		{
			name: "single element",
			in:   []string{"a"},
			want: []string{"a"},
		},
		{
			name: "multiple elements",
			in:   []string{"a", "b", "c"},
			want: []string{"c", "b", "a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Reverse(tt.in)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestAny(t *testing.T) {
	tests := []struct {
		name string
		f    func(string) bool
		in   []string
		want bool
	}{
		{
			name: "empty slice",
			in:   []string{},
			f: func(s string) bool {
				return true
			},
			want: false,
		},
		{
			name: "single element passes",
			in:   []string{"a"},
			f: func(s string) bool {
				return s == "a"
			},
			want: true,
		},
		{
			name: "single element fails",
			in:   []string{"a"},
			f: func(s string) bool {
				return s == "b"
			},
			want: false,
		},
		{
			name: "multiple elements, two pass",
			in:   []string{"a", "b", "c"},
			f: func(s string) bool {
				return s != "a"
			},
			want: true,
		},
		{
			name: "multiple elements, none pass",
			in:   []string{"a", "b", "c"},
			f: func(s string) bool {
				return s == "d"
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Any(tt.in, tt.f)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestApply(t *testing.T) {
	tests := []struct {
		name string
		f    func(string) string
		in   []string
		want []string
	}{
		{
			name: "empty slice",
			in:   []string{},
			f: func(s string) string {
				return s
			},
			want: []string{},
		},
		{
			name: "single element",
			in:   []string{"a"},
			f: func(s string) string {
				return s + "b"
			},
			want: []string{"ab"},
		},
		{
			name: "multiple elements",
			in:   []string{"a", "b", "c"},
			f: func(s string) string {
				return s + "b"
			},
			want: []string{"ab", "bb", "cb"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Apply(tt.in, tt.f)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestZipApply(t *testing.T) {
	tests := []struct {
		name string
		f    func(string, string) string
		as   []string
		bs   []string
		want []string
	}{
		{
			name: "empty slices",
			as:   []string{},
			bs:   []string{},
			f: func(a, b string) string {
				return a + b
			},
			want: []string{},
		},
		{
			name: "single element",
			as:   []string{"a"},
			bs:   []string{"b"},
			f: func(a, b string) string {
				return a + b
			},
			want: []string{"ab"},
		},
		{
			name: "multiple elements",
			as:   []string{"a", "b", "c"},
			bs:   []string{"d", "e", "f"},
			f: func(a, b string) string {
				return a + b
			},
			want: []string{"ad", "be", "cf"},
		},
		{
			name: "different length slices",
			as:   []string{"a", "b", "c"},
			bs:   []string{"d", "e"},
			f: func(a, b string) string {
				return a + b
			},
			want: []string{"ad", "be"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ZipApply(tt.f, tt.as, tt.bs)
			assert.Equal(t, got, tt.want)
		})
	}
}
