package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

type idDispatch struct {
	description string
	keys        []string
	handler     func(ids []int64)
}

func (d *idDispatch) describe() string {
	if len(d.keys) > 1 {
		d.keys[len(d.keys)-2] = d.keys[len(d.keys)-2] + " and " + d.keys[len(d.keys)-1]
		d.keys = d.keys[:len(d.keys)-1]
	}
	return fmt.Sprintf("with %s\n\t this endpoint returns %s", strings.Join(d.keys, ", "), d.description)
}

// urlParamFilter calls the first handler that has all keys in the url params with the parsed ids.
func urlParamDispatch(params url.Values, targets []idDispatch) error {
	for _, target := range targets {
		var ids []int64
		for _, key := range target.keys {
			id, err := strconv.ParseInt(params.Get(key), 10, 64)
			if err != nil {
				break
			}
			ids = append(ids, id)
		}
		if len(ids) == len(target.keys) {
			target.handler(ids)
			return nil
		}
	}

	var descriptions []string
	for _, target := range targets {
		descriptions = append(descriptions, target.describe())
	}
	return fmt.Errorf("Error - bad combination of URL paramters. Valid usages: \n\n%s", strings.Join(descriptions, "\n"))
}
