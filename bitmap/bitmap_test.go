package bitmap

import (
	"strconv"
	"testing"
)

func TestNewBitMap(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
		want *BitMap
	}{
		{name: "test1", args: args{size: 20}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm := NewBitMap(tt.args.size)
			if bm == nil {
				t.Errorf("NewBitMap() = %v, want %v", bm, tt.want)
			}
			bm.Add([]byte("1"))
			bm.Add([]byte("2"))
			bm.Add([]byte("3"))
			bm.Add([]byte("4"))
			bm.Add([]byte("5"))
			bm.Add([]byte("6"))
			bm.Add([]byte("7"))
			bm.Add([]byte("8"))
			bm.Add([]byte("9"))
			bm.Add([]byte("10"))
			if !bm.Contains([]byte("1")) {
				t.Errorf("NewBitMap() = %v, want %v", bm, tt.want)
			}
			if !bm.Contains([]byte("2")) {
				t.Errorf("NewBitMap() = %v, want %v", bm, tt.want)
			}
			if !bm.Contains([]byte("3")) {
				t.Errorf("NewBitMap() = %v, want %v", bm, tt.want)
			}
			if !bm.Contains([]byte("4")) {
				t.Errorf("NewBitMap() = %v, want %v", bm, tt.want)
			}
			if !bm.Contains([]byte("5")) {
				t.Errorf("NewBitMap() = %v, want %v", bm, tt.want)
			}
			if !bm.Contains([]byte("6")) {
				t.Errorf("NewBitMap() = %v, want %v", bm, tt.want)
			}
			if !bm.Contains([]byte("7")) {
				t.Errorf("NewBitMap() = %v, want %v", bm, tt.want)
			}
			if !bm.Contains([]byte("8")) {
				t.Errorf("NewBitMap() = %v, want %v", bm, tt.want)
			}
			if !bm.Contains([]byte("9")) {
				t.Errorf("NewBitMap() = %v, want %v", bm, tt.want)
			}
			if !bm.Contains([]byte("10")) {
				t.Errorf("NewBitMap() = %v, want %v", bm, tt.want)
			}
		})
	}
}

func BenchmarkNewBitMap(b *testing.B) {
	size := 10000000
	bm := NewBitMap(size)
	b.Run("Add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bm.Add([]byte(strconv.Itoa(i)))
		}
	})
}

func TestBitMap_AddCheck(t *testing.T) {
	type args struct {
		size int
		num  int
	}
	tests := []struct {
		name              string
		args              args
		wantCollisionRate float64
	}{
		{name: "t1", args: args{size: 1, num: 1}, wantCollisionRate: 0},
		{name: "t2", args: args{size: 1, num: 9}, wantCollisionRate: 8},
		{name: "t3", args: args{size: 3, num: 1}, wantCollisionRate: 0.01},
		{name: "t4", args: args{size: 30, num: 10}, wantCollisionRate: 0.01},
		{name: "t5", args: args{size: 300, num: 100}, wantCollisionRate: 0.06},
		{name: "t6", args: args{size: 3000, num: 1000}, wantCollisionRate: 0.09},
		{name: "t7", args: args{size: 30000, num: 10000}, wantCollisionRate: 0.099},
		{name: "t8", args: args{size: 300000, num: 100000}, wantCollisionRate: 0.09},
		{name: "t9", args: args{size: 3000000, num: 1000000}, wantCollisionRate: 0.09},
		{name: "t10", args: args{size: 30000000, num: 10000000}, wantCollisionRate: 0.09},
		{name: "t11", args: args{size: 300000000, num: 100000000}, wantCollisionRate: 0.09},
		{name: "t12", args: args{size: 3000000000, num: 1000000000}, wantCollisionRate: 0.09},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm := NewBitMap(tt.args.size)
			collisions := 0
			for i := 0; i < tt.args.num; i++ {
				if ok := bm.Add([]byte(strconv.Itoa(i))); ok {
					collisions++
				}
			}
			collisionRate := float64(collisions) / float64(tt.args.num)
			if collisionRate > tt.wantCollisionRate {
				t.Errorf("冲突率：got %.2f ,want %.2f", collisionRate, tt.wantCollisionRate)
				return
			}
			t.Logf("冲突率： %.2f", collisionRate)
		})
	}
}

func TestCollisionRate(t *testing.T) {
	type args struct {
		m int64
		k int64
		n int64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{name: "t1", args: args{m: 20, k: 10, n: 100}, want: 1},
		{name: "t2", args: args{m: 20, k: 3, n: 6}, want: 0.61},
		{name: "t3", args: args{m: 900, k: 3, n: 3}, want: 0.01},
		{name: "t4", args: args{m: 10000000, k: 3, n: 3}, want: 0.01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CollisionRate(tt.args.m, tt.args.k, tt.args.n); got > tt.want {
				t.Errorf("CollisionRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitMap_Contains(t *testing.T) {
	type fields struct {
		size int
	}
	type args struct {
		value   [][]byte
		contain []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{name: "t1", fields: fields{size: 12}, args: args{value: [][]byte{{'1'}, []byte{'2'}, []byte{'3'}, []byte{'4'}}, contain: []byte{'1'}}, want: true},
		{name: "t2", fields: fields{size: 12}, args: args{value: [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d")}, contain: []byte("a")}, want: true},
		{name: "t3", fields: fields{size: 12}, args: args{value: [][]byte{{'a'}, []byte{'b'}, []byte{'c'}, []byte{'d'}}, contain: []byte{'a'}}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm := NewBitMap(tt.fields.size)
			for _, v := range tt.args.value {
				bm.Add(v)
			}
			if got := bm.Contains(tt.args.contain); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
