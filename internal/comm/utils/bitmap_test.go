package utils

import (
	"reflect"
	"testing"
)

func TestBitMap_FirstZero(t *testing.T) {
	type fields struct {
		bits []uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"1", fields{
			bits: []uint64{0, 0},
		}, 1},
		{"2", fields{
			bits: []uint64{1, 0},
		}, 2},
		{"7", fields{
			bits: []uint64{63, 0},
		}, 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BitMap{
				Bits: tt.fields.bits,
			}
			if got := b.FirstZero(); got != tt.want {
				t.Errorf("BitMap.FirstZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitMap_Set(t *testing.T) {
	type args struct {
		loc int
	}
	tests := []struct {
		name string
		args args
	}{
		{"set 1", args{1}},
		{"set 2", args{2}},
	}
	for _, tt := range tests {
		b := NewBitMap()
		t.Run(tt.name, func(t *testing.T) {
			b.Set(tt.args.loc)
			t.Logf("[%b]", b.Bits[0])
		})
	}
}

func TestNewBitMap(t *testing.T) {
	tests := []struct {
		name string
		want *BitMap
	}{
		{"ok", &BitMap{Bits: make([]uint64, 0)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBitMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBitMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitMap_Unset(t *testing.T) {
	type fields struct {
		Bits []uint64
	}
	type args struct {
		loc int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"first", fields{[]uint64{1, 0}}, args{1}},
		{"second", fields{[]uint64{2, 0}}, args{2}},
		{"third", fields{[]uint64{3, 0}}, args{1}},
		{"forth", fields{[]uint64{4, 0}}, args{3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BitMap{
				Bits: tt.fields.Bits,
			}
			t.Logf("[%b]", b.Bits[0])
			b.Unset(tt.args.loc)
			t.Logf("[%b]", b.Bits[0])
		})
	}
}
