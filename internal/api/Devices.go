package api

import (
	"fmt"
	"github.com/bart-lute/addigy-tools/internal/models/device_entities"
	"log/slog"
	"net/http"
)

func Devices(desiredFactIdentifiers *[]string, queryFilter *device_entities.QueryFilter, sortDirection string, sortField string, items *[]device_entities.DeviceAudit) {
	slog.Debug(fmt.Sprintf("Devices API"))

	deviceFilter := device_entities.DeviceFilter{
		DesiredFactIdentifiers: *desiredFactIdentifiers,
		Page:                   1,
		PerPage:                50,
		Query:                  *queryFilter,
		SortDirection:          sortDirection,
		SortField:              sortField,
	}

	for {
		var deviceAuditResponse device_entities.DeviceAuditResponse
		doRequest(http.MethodPost, "devices/", &deviceFilter, &deviceAuditResponse)
		*items = append(*items, deviceAuditResponse.Items...)
		if deviceAuditResponse.Metadata.Page == deviceAuditResponse.Metadata.PageCount {
			break
		}
		deviceFilter.Page++
	}
	slog.Debug(fmt.Sprintf("Found %d Devices API items", len(*items)))
}
