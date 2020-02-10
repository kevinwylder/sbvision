package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func getIDs(r *http.Request, fields []string) ([]int64, error) {
	var ids []int64
	for _, name := range fields {
		id, err := strconv.ParseInt(r.Form.Get(name), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Missing %s in query params", name)
		}
		ids = append(ids, id)
	}
	return ids, nil
}
