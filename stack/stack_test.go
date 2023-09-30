package stack

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStack_IsEmpty(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		s    Stack
		want bool
	}{
		{
			name: "Empty Stack",
			s:    Stack{},
			want: true,
		},
		{
			name: "Non-empty Stack",
			s:    Stack{"item1", "item2"},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, tt.s.IsEmpty())
		})
	}
}

func TestStack_Peek(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		s    Stack
		want string
	}{
		{
			name: "Peek on Empty Stack",
			s:    Stack{},
			want: "",
		},
		{
			name: "Peek on Non-empty Stack",
			s:    Stack{"item1", "item2"},
			want: "item2",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, tt.s.Peek())
		})
	}
}

func TestStack_Pop(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		s    Stack
		want string
	}{
		{
			name: "Pop from Empty Stack",
			s:    Stack{},
			want: "",
		},
		{
			name: "Pop from Non-empty Stack",
			s:    Stack{"item1", "item2"},
			want: "item2",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, tt.s.Pop())
		})
	}
}

// ?
func TestStack_Push(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		s    Stack
		args string
		want []string
	}{
		{
			name: "Push to Empty Stack",
			s:    Stack{},
			args: "item1",
			want: []string{"item1"},
		},
		{
			name: "Push to Non-empty Stack",
			s:    Stack{"item1", "item2"},
			args: "item3",
			want: []string{"item1", "item2", "item3"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.s.Push(tt.args)
		})
	}
}
