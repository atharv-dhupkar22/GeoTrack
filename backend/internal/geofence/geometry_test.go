package geofence

import "testing"

func validPolygon() [][]float64 {
	return [][]float64{
		{19.0760, 72.8777},
		{19.0860, 72.8777},
		{19.0860, 72.8877},
		{19.0760, 72.8877},
		{19.0760, 72.8777},
	}
}

func TestValidatePolygon(t *testing.T) {
	tests := []struct {
		name        string
		coordinates [][]float64
		wantError   bool
	}{
		{
			name:        "valid closed polygon",
			coordinates: validPolygon(),
			wantError:   false,
		},
		{
			name: "polygon is not closed",
			coordinates: [][]float64{
				{19.0760, 72.8777},
				{19.0860, 72.8777},
				{19.0860, 72.8877},
				{19.0760, 72.8877},
			},
			wantError: true,
		},
		{
			name: "latitude out of range",
			coordinates: [][]float64{
				{91.0, 72.8777},
				{19.0860, 72.8777},
				{19.0860, 72.8877},
				{91.0, 72.8777},
			},
			wantError: true,
		},
		{
			name: "longitude out of range",
			coordinates: [][]float64{
				{19.0760, 181.0},
				{19.0860, 72.8777},
				{19.0860, 72.8877},
				{19.0760, 181.0},
			},
			wantError: true,
		},
		{
			name: "polygon has insufficient coordinates",
			coordinates: [][]float64{
				{19.0760, 72.8777},
				{19.0860, 72.8777},
				{19.0760, 72.8777},
			},
			wantError: true,
		},
		{
			name: "self intersecting polygon",
			coordinates: [][]float64{
				{0, 0},
				{1, 1},
				{0, 1},
				{1, 0},
				{0, 0},
			},
			wantError: true,
		},
		{
			name: "polygon with repeated vertices",
			coordinates: [][]float64{
				{19.0760, 72.8777},
				{19.0760, 72.8777},
				{19.0760, 72.8777},
				{19.0760, 72.8777},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePolygon(tt.coordinates)

			if tt.wantError && err == nil {
				t.Fatal("expected an error, but got nil")
			}

			if !tt.wantError && err != nil {
				t.Fatalf("expected no error, but got: %v", err)
			}
		})
	}
}

func TestIsPointInPolygon(t *testing.T) {
	polygon := validPolygon()

	tests := []struct {
		name      string
		latitude  float64
		longitude float64
		want      bool
	}{
		{
			name:      "point inside polygon",
			latitude:  19.0810,
			longitude: 72.8820,
			want:      true,
		},
		{
			name:      "point outside polygon",
			latitude:  19.1000,
			longitude: 72.9000,
			want:      false,
		},
		{
			name:      "point south of polygon",
			latitude:  19.0700,
			longitude: 72.8820,
			want:      false,
		},
		{
			name:      "point west of polygon",
			latitude:  19.0810,
			longitude: 72.8700,
			want:      false,
		},
		{
			name:      "point on polygon vertex",
			latitude:  19.0760,
			longitude: 72.8777,
			want:      true,
		},
		{
			name:      "point on polygon edge",
			latitude:  19.0760,
			longitude: 72.8820,
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPointInPolygon(
				tt.latitude,
				tt.longitude,
				polygon,
			)

			if got != tt.want {
				t.Errorf(
					"IsPointInPolygon(%f, %f) = %v, want %v",
					tt.latitude,
					tt.longitude,
					got,
					tt.want,
				)
			}
		})
	}
}
