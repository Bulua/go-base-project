package routereg

import "regexp"

var paramPattern = regexp.MustCompile(`\{[^}]+\}`)

// Route is a single HTTP endpoint to be synced into gbp_api_resources.
type Route struct {
	Method string
	Path   string
	Group  string
	Desc   string
}

var registered []Route

// Add records a route. Path params like {id}, {dictCode} are normalized to {id}
// so they match the format stored in gbp_api_resources and used by the permission middleware.
func Add(method, path, group, desc string) {
	registered = append(registered, Route{
		Method: method,
		Path:   paramPattern.ReplaceAllString(path, "{id}"),
		Group:  group,
		Desc:   desc,
	})
}

// All returns all registered routes.
func All() []Route { return registered }
