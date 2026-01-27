package d2_test

import (
	"testing"

	"github.com/kgatilin/archai/internal/adapter/d2"
	"github.com/kgatilin/archai/internal/domain"
)

func TestStereotypeColor(t *testing.T) {
	tests := []struct {
		name       string
		stereotype domain.Stereotype
		wantColor  string
	}{
		{
			name:       "service stereotype returns purple",
			stereotype: domain.StereotypeService,
			wantColor:  d2.ColorPurple,
		},
		{
			name:       "repository stereotype returns purple",
			stereotype: domain.StereotypeRepository,
			wantColor:  d2.ColorPurple,
		},
		{
			name:       "port stereotype returns purple",
			stereotype: domain.StereotypePort,
			wantColor:  d2.ColorPurple,
		},
		{
			name:       "interface stereotype returns purple",
			stereotype: domain.StereotypeInterface,
			wantColor:  d2.ColorPurple,
		},
		{
			name:       "factory stereotype returns green",
			stereotype: domain.StereotypeFactory,
			wantColor:  d2.ColorGreen,
		},
		{
			name:       "aggregate stereotype returns blue",
			stereotype: domain.StereotypeAggregate,
			wantColor:  d2.ColorBlue,
		},
		{
			name:       "entity stereotype returns blue",
			stereotype: domain.StereotypeEntity,
			wantColor:  d2.ColorBlue,
		},
		{
			name:       "value stereotype returns gray",
			stereotype: domain.StereotypeValue,
			wantColor:  d2.ColorGray,
		},
		{
			name:       "enum stereotype returns gray",
			stereotype: domain.StereotypeEnum,
			wantColor:  d2.ColorGray,
		},
		{
			name:       "no stereotype returns gray",
			stereotype: domain.StereotypeNone,
			wantColor:  d2.ColorGray,
		},
		{
			name:       "unknown stereotype returns gray",
			stereotype: domain.Stereotype("unknown"),
			wantColor:  d2.ColorGray,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d2.StereotypeColor(tt.stereotype)
			if got != tt.wantColor {
				t.Errorf("StereotypeColor(%q) = %q, want %q", tt.stereotype, got, tt.wantColor)
			}
		})
	}
}

func TestStereotypeLabel(t *testing.T) {
	tests := []struct {
		name       string
		stereotype domain.Stereotype
		wantLabel  string
	}{
		{
			name:       "service stereotype returns formatted label",
			stereotype: domain.StereotypeService,
			wantLabel:  "<<service>>",
		},
		{
			name:       "repository stereotype returns formatted label",
			stereotype: domain.StereotypeRepository,
			wantLabel:  "<<repository>>",
		},
		{
			name:       "factory stereotype returns formatted label",
			stereotype: domain.StereotypeFactory,
			wantLabel:  "<<factory>>",
		},
		{
			name:       "aggregate stereotype returns formatted label",
			stereotype: domain.StereotypeAggregate,
			wantLabel:  "<<aggregate>>",
		},
		{
			name:       "no stereotype returns empty string",
			stereotype: domain.StereotypeNone,
			wantLabel:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d2.StereotypeLabel(tt.stereotype)
			if got != tt.wantLabel {
				t.Errorf("StereotypeLabel(%q) = %q, want %q", tt.stereotype, got, tt.wantLabel)
			}
		})
	}
}
