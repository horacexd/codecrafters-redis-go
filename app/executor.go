package main

import (
	"fmt"
	"time"
)


type d_val struct {
	val string
	ttl int64
	
	// internal
	createMilliTimestamp int64
} 

var dict = map[string]d_val{}

func executeCommandSetWithExpiry(k string, v string, ttl int64) error {
	val := d_val{
		val: v,
		ttl: ttl,
		createMilliTimestamp: time.Now().UnixMilli(),
	}
	dict[k] = val
	return nil
}

func executeCommandGet(k string) (string, error) {
	v, ok := dict[k]
	if !ok {
		return "-1", nil
	}

	// TODO: add periodically key clean up for dict with ttl
	if v.ttl != -1 && v.createMilliTimestamp + v.ttl < time.Now().UnixMilli() {
		delete(dict, k)
		return "-1", nil
	}

	return v.val, nil
}

func executeCommandConfigGet(k string) ([]string, error) {
	switch k {
	case "dir":
		return []string{k, DIR}, nil
	case "dirfilename":
		return []string{k, DB_FILE_NAME}, nil
	default:
		return nil, fmt.Errorf("unknown config name. %s", k)
	}
}