package qlog

import (
	"reflect"
	"testing"
)

func TestNewMapS(t *testing.T) {
	tests := []struct {
		name string
		want *FlatMapS
	}{
		{
			"New FlatMapS",
			&FlatMapS{
				[]string{},
				[]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMapS(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMapS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMapI(t *testing.T) {
	tests := []struct {
		name string
		want *FlatMapI
	}{
		{
			"New FlatMapI",
			&FlatMapI{
				[]string{},
				[]interface{}{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMapI(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMapI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlatMapS_Add(t *testing.T) {
	type fields struct {
		K []string
		V []string
	}
	type args struct {
		key string
		val string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *FlatMapS
		wantErr bool
	}{
		{
			"Add FlatMapS",
			fields{
				[]string{"new"},
				[]string{"one"},
			},
			args{
				"foo",
				"bar",
			},
			&FlatMapS{
				[]string{"new", "foo"},
				[]string{"one", "bar"},
			},
			false,
		},
		{
			"Add FlatMapS with replace val",
			fields{
				[]string{"new"},
				[]string{"one"},
			},
			args{
				"new",
				"two",
			},
			&FlatMapS{
				[]string{"new"},
				[]string{"two"},
			},
			false,
		},
		{
			"Add FlatMapS with replace val (error)",
			fields{
				[]string{"new"},
				[]string{"one"},
			},
			args{
				"new",
				"two",
			},
			&FlatMapS{
				[]string{"new"},
				[]string{"one"},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &FlatMapS{
				K: tt.fields.K,
				V: tt.fields.V,
			}
			m.Add(tt.args.key, tt.args.val)
			if !reflect.DeepEqual(m, tt.want) && !tt.wantErr {
				t.Errorf("FlatMapS_Add() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestFlatMapI_Add(t *testing.T) {
	type fields struct {
		K []string
		V []interface{}
	}
	type args struct {
		key string
		val interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *FlatMapI
		wantErr bool
	}{
		{
			"Add FlatMapI",
			fields{
				[]string{"new"},
				[]interface{}{"one"},
			},
			args{
				"foo",
				"bar",
			},
			&FlatMapI{
				[]string{"new", "foo"},
				[]interface{}{"one", "bar"},
			},
			false,
		},
		{
			"Add FlatMapI with replace val",
			fields{
				[]string{"new"},
				[]interface{}{"one"},
			},
			args{
				"new",
				"two",
			},
			&FlatMapI{
				[]string{"new"},
				[]interface{}{"two"},
			},
			false,
		},
		{
			"Add FlatMapI with replace val (error)",
			fields{
				[]string{"new"},
				[]interface{}{"one"},
			},
			args{
				"new",
				"two",
			},
			&FlatMapI{
				[]string{"new"},
				[]interface{}{"one"},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &FlatMapI{
				K: tt.fields.K,
				V: tt.fields.V,
			}
			m.Add(tt.args.key, tt.args.val)
			if !reflect.DeepEqual(m, tt.want) && !tt.wantErr {
				t.Errorf("FlatMapS_Add() = %v, want %v", m, tt.want)
			}
		})
	}
}
