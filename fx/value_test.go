package fx

import (
	"testing"

	"github.com/soderasen-au/go-common/util"
)

func TestValue_Less(t *testing.T) {
	type fields struct {
		Text      *string
		IsNumeric bool
		Number    *float64
		Error     *util.Result
	}
	type args struct {
		rh Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  *Value
	}{
		{
			name: "less",
			fields: fields{
				Text:      util.Ptr("abc"),
				IsNumeric: false,
				Number:    nil,
				Error:     nil,
			},
			args: args{rh: Value{
				Text:      util.Ptr("bcd"),
				IsNumeric: false,
				Number:    nil,
				Error:     nil,
			}},
			want:  true,
			want1: True(),
		},
		{
			name: "great",
			fields: fields{
				Text:      util.Ptr("bcd"),
				IsNumeric: false,
				Number:    nil,
				Error:     nil,
			},
			args: args{rh: Value{
				Text:      util.Ptr("bcd"),
				IsNumeric: false,
				Number:    nil,
				Error:     nil,
			}},
			want:  false,
			want1: False(),
		},
		{
			name: "equal",
			fields: fields{
				Text:      util.Ptr("abc"),
				IsNumeric: false,
				Number:    nil,
				Error:     nil,
			},
			args: args{rh: Value{
				Text:      util.Ptr("abc"),
				IsNumeric: false,
				Number:    nil,
				Error:     nil,
			}},
			want:  false,
			want1: False(),
		},
		{
			name: "lessN",
			fields: fields{
				Text:      nil,
				IsNumeric: true,
				Number:    util.Ptr(1.0),
				Error:     nil,
			},
			args: args{rh: Value{
				Text:      nil,
				IsNumeric: true,
				Number:    util.Ptr(1.1),
				Error:     nil,
			}},
			want:  true,
			want1: True(),
		},
		{
			name: "greatN",
			fields: fields{
				Text:      nil,
				IsNumeric: true,
				Number:    util.Ptr(1.1),
				Error:     nil,
			},
			args: args{rh: Value{
				Text:      nil,
				IsNumeric: true,
				Number:    util.Ptr(1.0),
				Error:     nil,
			}},
			want:  false,
			want1: False(),
		},
		{
			name: "equalN",
			fields: fields{
				Text:      nil,
				IsNumeric: true,
				Number:    util.Ptr(1.0),
				Error:     nil,
			},
			args: args{rh: Value{
				Text:      nil,
				IsNumeric: true,
				Number:    util.Ptr(1.0),
				Error:     nil,
			}},
			want:  false,
			want1: False(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lh := &Value{
				Text:      tt.fields.Text,
				IsNumeric: tt.fields.IsNumeric,
				Number:    tt.fields.Number,
				Error:     tt.fields.Error,
			}
			got := lh.Less(&tt.args.rh)
			if got.Equal(tt.want1).Not().True() {
				t.Errorf("Less() got = %v, want %v", util.JsonStr(got), util.JsonStr(tt.want1))
			}
		})
	}
}
