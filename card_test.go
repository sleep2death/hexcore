package hexcore

import "testing"

func TestCard_ToString(t *testing.T) {
	type fields struct {
		ID             string
		Name           string
		ImgURL         string
		Cost           int
		RawDescription string
		Type           CardType
		Color          CardColor
		Rarity         CardRarity
		Target         CardTarget
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card := &Card{
				ID:             tt.fields.ID,
				Name:           tt.fields.Name,
				ImgURL:         tt.fields.ImgURL,
				Cost:           tt.fields.Cost,
				RawDescription: tt.fields.RawDescription,
				Type:           tt.fields.Type,
				Color:          tt.fields.Color,
				Rarity:         tt.fields.Rarity,
				Target:         tt.fields.Target,
			}
			if got := card.ToString(); got != tt.want {
				t.Errorf("Card.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
