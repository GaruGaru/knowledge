package data

import (
	"testing"
)

func TestPaginationRequest_Offset(t *testing.T) {
	type fields struct {
		Page     int
		PageSize int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "offset 0 with page 1",
			fields: fields{
				Page:     1,
				PageSize: 100,
			},
			want: 0,
		},
		{
			name: "offset 10 with page 2 page size 10",
			fields: fields{
				Page:     2,
				PageSize: 10,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PaginationRequest{
				Page:     tt.fields.Page,
				PageSize: tt.fields.PageSize,
			}
			if got := p.Offset(); got != tt.want {
				t.Errorf("Offset() = %v, want %v", got, tt.want)
			}
		})
	}
}
