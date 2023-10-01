package util

import (
	"encoding/json"
	"github.com/rogpeppe/go-internal/diff"
	"reflect"
)

const (
	DiffAdd string = "add"
	DiffDel string = "delete"
	DiffMod string = "modify"
)

type DiffMarshaler interface {
	MarshalForDiff() ([]byte, *Result)
}

type Diff struct {
	From []byte
	To   []byte
	Type string
	Path string
	Diff []byte
}

func NewDiff(path string, from, to DiffMarshaler) (*Diff, *Result) {
	if from != nil && (reflect.ValueOf(from).Kind() == reflect.Ptr && reflect.ValueOf(from).IsNil()) {
		from = nil
	}
	if to != nil && (reflect.ValueOf(to).Kind() == reflect.Ptr && reflect.ValueOf(to).IsNil()) {
		to = nil
	}

	if from == to {
		return nil, nil
	}

	var fromBuf, toBuf []byte
	var res *Result
	if from != nil {
		fromBuf, res = from.MarshalForDiff()
		if res != nil {
			return nil, res.With("from.MarshalForDiff")
		}
		fromBuf = append(fromBuf, '\n')
	}
	if to != nil {
		toBuf, res = to.MarshalForDiff()
		if res != nil {
			return nil, res.With("to.MarshalForDiff")
		}
		toBuf = append(toBuf, '\n')
	}

	diffBuf := diff.Diff("from", fromBuf, "to", toBuf)
	if len(diffBuf) == 0 {
		return nil, nil
	}

	d := Diff{
		From: fromBuf,
		To:   toBuf,
		Path: path,
		Diff: diffBuf,
	}
	if from == nil && to != nil {
		d.Type = DiffAdd
	} else if to == nil && from != nil {
		d.Type = DiffDel
	} else {
		d.Type = DiffMod
	}

	return &d, nil
}

func NewJsonDiff[T any](path string, from, to *T) (*Diff, *Result) {
	if from == to {
		return nil, nil
	}

	fromBuf, err := json.MarshalIndent(from, "", "  ")
	if err != nil {
		return nil, Error("MarshalFrom", err)
	}
	fromBuf = append(fromBuf, '\n')
	toBuf, err := json.MarshalIndent(to, "", "  ")
	if err != nil {
		return nil, Error("MarshalTo", err)
	}
	toBuf = append(toBuf, '\n')
	diffBuf := diff.Diff("from", fromBuf, "to", toBuf)

	d := Diff{
		From: fromBuf,
		To:   toBuf,
		Path: path,
		Diff: diffBuf,
	}
	if from == nil {
		d.Type = DiffAdd
	} else if to == nil {
		d.Type = DiffDel
	} else {
		d.Type = DiffMod
	}

	return &d, nil
}