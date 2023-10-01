package util

import (
	"fmt"
	"testing"
)

func TestNewItemDiff(t *testing.T) {
	type args[T any] struct {
		path string
		from *T
		to   *T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want string
	}
	tests := []testCase[string]{
		{
			name: "string-del",
			args: args[string]{
				path: "string-del",
				from: Ptr("from"),
				to:   nil,
			},
			want: DiffDel,
		},
		{
			name: "string-add",
			args: args[string]{
				path: "string-add",
				from: nil,
				to:   Ptr("to"),
			},
			want: DiffAdd,
		},
		{
			name: "string-mod",
			args: args[string]{
				path: "string-del",
				from: Ptr("from"),
				to:   Ptr("to"),
			},
			want: DiffMod,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff, res := NewJsonDiff(tt.args.path, tt.args.from, tt.args.to)
			if res != nil {
				t.Errorf("error: %s", res.Error())
			}
			if diff != nil && diff.Type != tt.want {
				t.Errorf("want type `%s` but got `%s`", tt.want, diff.Type)
			} else {
				fmt.Printf("diff: %s\n", string(diff.Diff))
			}
		})
	}
}

func (r *Result) MarshalForDiff() ([]byte, *Result) {
	return Jsonify(r), nil
}

func TestNewDiff(t *testing.T) {
	r1 := MsgError("test", "res1")
	r2 := MsgError("test", "res2")
	r3 := MsgError("test", "res1")
	var r4 *Result

	type args struct {
		path string
		from DiffMarshaler
		to   DiffMarshaler
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Result-Add",
			args: args{
				path: "Result-Add",
				from: nil,
				to:   r2,
			},
			want: DiffAdd,
		},
		{
			name: "Result-Delete",
			args: args{
				path: "Result-Delete",
				from: r1,
				to:   nil,
			},
			want: DiffDel,
		},
		{
			name: "Result-Mod",
			args: args{
				path: "Result-Mod",
				from: r1,
				to:   r2,
			},
			want: DiffMod,
		},
		{
			name: "Result-Equal1",
			args: args{
				path: "Result-Equal1",
				from: nil,
				to:   nil,
			},
			want: "nil",
		},
		{
			name: "Result-Equal2",
			args: args{
				path: "Result-Equal2",
				from: r1,
				to:   r1,
			},
			want: "nil",
		},
		{
			name: "Result-Equal3",
			args: args{
				path: "Result-Equal3",
				from: r1,
				to:   r3,
			},
			want: "nil",
		},
		{
			name: "Result-Equal4-nil-1",
			args: args{
				path: "Result-Equal4-nil-1",
				from: nil,
				to:   r4,
			},
			want: "nil",
		},
		{
			name: "Result-Equal4-nil-2",
			args: args{
				path: "Result-Equal4-nil-2",
				from: r4,
				to:   r4,
			},
			want: "nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff, res := NewDiff(tt.args.path, tt.args.from, tt.args.to)
			if res != nil {
				t.Errorf("error: %s", res.Error())
			}

			if diff != nil {
				if tt.want == "nil" {
					t.Errorf("want nil diff but got: `%s`", string(diff.Diff))
					return
				}
				if diff.Type != tt.want {
					t.Errorf("want type `%s` but got `%s`", tt.want, diff.Type)
				} else {
					fmt.Printf("diff: %s\n", string(diff.Diff))
				}
			} else {
				if tt.want != "nil" {
					t.Errorf("want diff but got nil")
				}
			}
		})
	}
}
