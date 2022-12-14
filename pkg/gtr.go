/*
GTR (being short for Go Tiny Router) is a minimalistic router
that was initially developed to identify RESTful call signatures
for the purpose of caching them (based on parameters).

[Usage Examples]
GTR parses URLs based on Express style templates, for example:

	`http://www.abcdefg.com/api/v1/users/:username/details`

The provided URL can be successfully matched against the following
URLs:

	`http://www.abcdefg.com/api/v1/users/ken/details`
	`http://www.abcdefg.com/api/v1/users/dennis/details`

Alongside with route parameter, GTR also supports query string
specification, in such a way that if a query parameter is specified
a match will be successfull only and only if the target URL also
specifies that query parameter and that it is of the same value as
originally specified. For example:

	`http://www.abcdefg.com/api/v1/users/:username/details?type=cached`

The provided URL can be successfully matched against the following
URLs:

	`http://www.abcdefg.com/api/v1/users/ken/details?type=cached&format=JSON`
	`http://www.abcdefg.com/api/v1/users/dennis/details?type=cached`

However, the following URLs will NOT be successfullt matched:

	`http://www.abcdefg.com/api/v1/users/ken/details?format=JSON`
	`http://www.abcdefg.com/api/v1/users/dennis/details?`

This behavior has been designed intentional to serve the original purpose
of the library.
*/
package gtr

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
	"sync"
)

// Sentinel errors section
type RouterError string

func (routerError RouterError) Error() string {
	return string(routerError)
}

const (
	HOST_NOT_REGISTERED RouterError = "host not registered"
	NO_MATCH_FOUND      RouterError = "no match found"
	NO_URL_REGISTERED   RouterError = "no url registered"
)

var (
	_routeTable RouteTable
	_once       sync.Once
)

// The RouterTable is used to store information relating to routes
type RouteTable struct {
	routes  map[int][]*Route
	configs map[string]map[string]any
}

// The Route struct is used for breaking down a URL to segments
// based on which a route matching can take place
type Route struct {
	host        string
	routeParams map[int]string
	queryParams map[string]string
	hash        string
}

// Parses a URL to Route struct
// Examples:
//
//	   url, err := url.Parse("http://www.abcdefg.com/api/v1/users/:username/details")
//
//		  if err != nil {
//		      ...
//		  }
//
//	   route := ParseRoute(url)
func ParseRoute(url *url.URL) *Route {
	routeParams := make(map[int]string)
	queryParams := make(map[string]string)
	for index, segment := range strings.Split(url.Path, "/") {
		if len(segment) == 0 {
			continue
		}
		if strings.HasPrefix(segment, ":") {
			routeParams[index] = "?"
			continue
		}
		routeParams[index] = segment
	}

	for key, value := range url.Query() {
		sort.Slice(value, func(i, j int) bool {
			return value[i] > value[j]
		})
		queryParams[key] = strings.Join(value, ",")
	}
	hash := CreateHash(url)
	route := Route{
		host:        url.Host,
		routeParams: routeParams,
		queryParams: queryParams,
		hash:        hash,
	}
	return &route
}

// Compares two routes against each other
// Params:
//   - preferredRoute: The route template
//   - route: The route to match against the route template
func RouteCompare(preferredRoute *Route, route *Route) int {
	if len(preferredRoute.routeParams) != len(route.routeParams) {
		return 0
	}
	rank := 0
	for key, value := range preferredRoute.routeParams {
		if value == "?" {
			rank += 1
			continue
		}
		if value != route.routeParams[key] {
			rank = 0
			break
		}
		rank += 2
	}
	for key, value := range preferredRoute.queryParams {
		val, ok := route.queryParams[key]
		if !ok {
			return 0
		}
		if val != value {
			return 0
		}
	}
	return rank
}

// Creates a unique hash for a URL
func CreateHash(url *url.URL) string {
	buffer := bytes.NewBufferString(url.Path)
	if len(url.RawQuery) > 0 {
		buffer.WriteString("?")
		buffer.WriteString(url.RawQuery)
	}
	sha256 := sha256.New()
	sha256.Write(buffer.Bytes())
	hash := hex.EncodeToString(sha256.Sum(nil))
	return hash
}

// Gets the default route table
func DefaultRouteTable() *RouteTable {
	_once.Do(func() {
		_routeTable = RouteTable{
			routes:  map[int][]*Route{},
			configs: map[string]map[string]any{},
		}
	})
	return &_routeTable
}

// Registeres a new route to the route table
func (rt RouteTable) Register(url *url.URL, conf map[string]any) {
	route := ParseRoute(url)
	len := len(route.routeParams)
	if _, ok := rt.configs[route.hash]; ok {
		return
	}
	rt.configs[route.hash] = conf
	_, ok := rt.routes[len]
	if !ok {
		rt.routes[len] = make([]*Route, 0)
	}
	rt.routes[len] = append(rt.routes[len], route)
}

// Finds the route template for a given URL
func (rt RouteTable) Find(url *url.URL) (string, error) {
	if len(rt.routes) == 0 {
		return "", NO_URL_REGISTERED
	}
	prt := ParseRoute(url)
	routes, ok := rt.routes[len(prt.routeParams)]
	if !ok {
		return "", HOST_NOT_REGISTERED
	}
	lrnk := 0
	var lrt *Route
	for _, url := range routes {
		rnk := RouteCompare(url, prt)
		if rnk != 0 {
			if rnk > lrnk {
				lrnk = rnk
				lrt = url
			}
		}
	}
	if lrnk == 0 {
		return "", NO_MATCH_FOUND
	}
	return lrt.hash, nil
}

// Gets configuration for a given hash
func (rt RouteTable) GetConfig(hash string) map[string]any {
	return rt.configs[hash]
}
