package diehard

import "testing"

func TestTransfer3To5(t *testing.T) {
	testCases := []struct {
		desc         string
		j3, j5       int
		want3, want5 int
	}{
		{
			desc:  "overflow",
			j3:    3,
			j5:    3,
			want3: 1,
			want5: 5,
		}, {
			desc:  "no overflow",
			j3:    1,
			j5:    1,
			want3: 0,
			want5: 2,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctx := &context{
				jug_3: tC.j3,
				jug_5: tC.j5,
			}

			transfer_3_to_5(ctx)

			if ctx.jug_3 != tC.want3 {
				t.Fatalf("j3, want %d, got %d", tC.want3, ctx.jug_3)
			}
			if ctx.jug_5 != tC.want5 {
				t.Fatalf("j5, want %d, got %d", tC.want5, ctx.jug_5)
			}
		})
	}
}

func TestTransfer5To3(t *testing.T) {
	testCases := []struct {
		desc         string
		j3, j5       int
		want3, want5 int
	}{
		{
			desc:  "exact",
			j3:    0,
			j5:    3,
			want3: 3,
			want5: 0,
		}, {
			desc:  "overflow",
			j3:    2,
			j5:    2,
			want3: 3,
			want5: 1,
		}, {
			desc:  "no overflow",
			j3:    1,
			j5:    2,
			want3: 3,
			want5: 0,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctx := &context{
				jug_3: tC.j3,
				jug_5: tC.j5,
			}

			transfer_5_to_3(ctx)

			if ctx.jug_3 != tC.want3 {
				t.Fatalf("j3, want %d, got %d", tC.want3, ctx.jug_3)
			}
			if ctx.jug_5 != tC.want5 {
				t.Fatalf("j5, want %d, got %d", tC.want5, ctx.jug_5)
			}
		})
	}
}
