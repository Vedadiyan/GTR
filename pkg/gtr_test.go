package gtr

import (
	"net/url"
	"testing"
	"time"
)

func PrepareURLTemplate(t *testing.T) *url.URL {
	const (
		TEMPLATE = "http://www.abcdefg.com/api/v1/users/:username/details?type=cache"
	)
	url, err := url.Parse(TEMPLATE)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	return url
}

func PrepareURL(t *testing.T) *url.URL {
	const (
		TEMPLATE = "http://www.abcdefg.com/api/v1/users/ken/details?type=cache"
	)
	url, err := url.Parse(TEMPLATE)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	return url
}

func PrepareParseRoute(t *testing.T) *Route {
	url := PrepareURLTemplate(t)
	route := ParseRoute(url)
	if len(route.routeParams) != 5 {
		t.Log("route params parsed incorrectly")
		t.FailNow()
	}
	if len(route.queryParams) != 1 {
		t.Log("query params parsed incorrectly")
		t.FailNow()
	}
	return route
}

func PrepareFind(t *testing.T) string {
	config := make(map[string]any)
	config["ttl"] = time.Second
	DefaultRouteTable().Register(PrepareURLTemplate(t), config)
	hash, err := DefaultRouteTable().Find(PrepareURL(t))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	return hash
}

func TestRouteCompare(t *testing.T) {
	template := PrepareURLTemplate(t)
	preferredRoute := ParseRoute(template)
	url := PrepareURL(t)
	route := ParseRoute(url)
	rank := RouteCompare(preferredRoute, route)
	if rank == 0 {
		t.Log("route matching failed")
		t.FailNow()
	}
}

func TestParseRoute(t *testing.T) {
	_ = PrepareParseRoute(t)
}

func TestFind(t *testing.T) {
	_ = PrepareFind(t)
}

func TestGetConfig(t *testing.T) {
	hash := PrepareFind(t)
	config := DefaultRouteTable().GetConfig(hash)
	value, ok := config["ttl"]
	if !ok {
		t.Log("could not get correct config")
		t.FailNow()
	}
	if value != time.Second {
		t.Log("config is invalid")
		t.FailNow()
	}
}
