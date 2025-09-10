package api

import (
	"fmt"
	"github.com/bart-lute/addigy-tools/internal/models/auditor_service"
	"github.com/bart-lute/addigy-tools/internal/models/fact_entities"
	"github.com/bart-lute/addigy-tools/internal/models/facts"
	"log/slog"
	"net/http"
)

func FactsCustomQuery(
	filter *facts.Filter,
	sortDirection string,
	sortField string,
	items *[]fact_entities.Fact,
) {

	paginatedCustomFactsRequestQuery := facts.PaginatedCustomFactsRequestQuery{
		Page:          1,
		PerPage:       50,
		Query:         filter,
		SortDirection: sortDirection,
		SortField:     sortField,
	}

	for {
		var factsResponse *auditor_service.FactsResponse
		doRequest(http.MethodPost, "facts/custom/query", &paginatedCustomFactsRequestQuery, &factsResponse)
		*items = append(*items, factsResponse.Items...)
		if factsResponse.Metadata.Page == factsResponse.Metadata.PageCount {
			break
		}
		paginatedCustomFactsRequestQuery.Page++
	}
	slog.Debug(fmt.Sprintf("Found %d Custom Facts", len(*items)))
}
