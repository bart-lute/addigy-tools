package auditor_service

import (
	"github.com/bart-lute/addigy-tools/internal/models/fact_entities"
	"github.com/bart-lute/addigy-tools/internal/models/response_entities"
)

type FactsResponse struct {
	Items    []fact_entities.Fact       `json:"items"`
	Metadata response_entities.Metadata `json:"metadata"`
}
