package handlers

import (
	"net/url"
	"strconv"
)

const defaultPageSize = 10

func getPageSize(v url.Values) uint64 {
	pageSizeString := v.Get("page_size")

	if pageSizeString == "" {
		return defaultPageSize
	}

	if pageSize, err := strconv.ParseUint(pageSizeString, 10, 64); err == nil {
		return pageSize
	}

	return defaultPageSize
}