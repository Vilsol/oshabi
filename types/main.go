package types

import "github.com/vilsol/oshabi/data"

type ParsedListing struct {
	Type  data.HarvestType `json:"type"`
	Count int              `json:"count"`
	Level int              `json:"level"`
}
