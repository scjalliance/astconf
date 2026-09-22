package astorg_test

import (
	"testing"

	"github.com/scjalliance/astconf/astorg"
)

func TestDataSetEqual(t *testing.T) {
	tests := []struct {
		name string
		a, b astorg.DataSet
		want bool
	}{
		{
			name: "empty",
			want: true,
		},
		{
			name: "same location",
			a:    astorg.DataSet{Locations: astorg.LocationList{{Name: "Seattle"}}},
			b:    astorg.DataSet{Locations: astorg.LocationList{{Name: "Seattle"}}},
			want: true,
		},
		{
			name: "different location",
			a:    astorg.DataSet{Locations: astorg.LocationList{{Name: "Seattle"}}},
			b:    astorg.DataSet{Locations: astorg.LocationList{{Name: "Tacoma"}}},
			want: false,
		},
		{
			name: "different location count",
			a:    astorg.DataSet{Locations: astorg.LocationList{{Name: "Seattle"}}},
			b:    astorg.DataSet{},
			want: false,
		},
		{
			name: "same person",
			a:    astorg.DataSet{People: astorg.PersonList{{Username: "fred"}}},
			b:    astorg.DataSet{People: astorg.PersonList{{Username: "fred"}}},
			want: true,
		},
		{
			name: "different person",
			a:    astorg.DataSet{People: astorg.PersonList{{Username: "fred"}}},
			b:    astorg.DataSet{People: astorg.PersonList{{Username: "barney"}}},
			want: false,
		},
		{
			name: "same phone",
			a:    astorg.DataSet{Phones: astorg.PhoneList{{MAC: "000000000001"}}},
			b:    astorg.DataSet{Phones: astorg.PhoneList{{MAC: "000000000001"}}},
			want: true,
		},
		{
			name: "different phone",
			a:    astorg.DataSet{Phones: astorg.PhoneList{{MAC: "000000000001"}}},
			b:    astorg.DataSet{Phones: astorg.PhoneList{{MAC: "000000000002"}}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Equal(&tt.b); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
