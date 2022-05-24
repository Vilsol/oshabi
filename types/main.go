package types

type HarvestType string

type ParsedListing struct {
	Type  HarvestType `json:"type"`
	Count int         `json:"count"`
	Level int         `json:"level"`
}
