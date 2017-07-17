// Copyright © 2014-2017 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package router

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/zenazn/goji/web"
)

const (
	ABSOLUTE_URL  = 0
	ABSOLUTE_PATH = 1
	NETWORK_PATH  = 3
)

var PatternMatching = regexp.MustCompile("(:[a-zA-Z]*)")

type route struct {
	path    string
	params  []string
	renders []func(values url.Values) (string, error)
}

func (r *route) compile() {
	result := PatternMatching.FindAllStringSubmatchIndex(r.path, -1)

	renders := make([]func(values url.Values) (string, error), 0)

	currentIndex := 0
	for _, match := range result {

		if match[0] != currentIndex {
			prefix := r.path[currentIndex:match[0]]
			renders = append(renders, func(values url.Values) (string, error) {
				return prefix, nil
			})
		}

		name := r.path[match[0]+1 : match[3]]

		renders = append(renders, func(values url.Values) (string, error) {
			v := values.Get(name)
			values.Del(name)
			return v, nil
		})

		currentIndex = match[3]
	}

	if currentIndex < len(r.path) {
		final := r.path[currentIndex:]
		renders = append(renders, func(values url.Values) (string, error) {
			return final, nil
		})
	}

	r.renders = renders
}

func NewRouter(mux *web.Mux) *Router {
	if mux == nil {
		mux = web.New()
	}

	return &Router{
		Routes: make(map[string]*route),
		Mux:    mux,
	}
}

type Router struct {
	Routes map[string]*route
	Mux    *web.Mux
}

func (u *Router) GenerateUrl(name string, params url.Values, context *RequestContext) (string, error) {
	return u.generate(name, params, ABSOLUTE_URL, context)
}

func (u *Router) GeneratePath(name string, params url.Values) (string, error) {
	return u.generate(name, params, ABSOLUTE_PATH, nil)
}

func (u *Router) GenerateNet(name string, params url.Values) (string, error) {
	return u.generate(name, params, NETWORK_PATH, nil)
}

func (u *Router) Handle(name, pattern string, handler web.HandlerType) *Router {
	u.Mux.Handle(pattern, handler)
	u.addRoute(name, pattern)

	return u
}

func (u *Router) Get(name, pattern string, handler web.HandlerType) *Router {
	u.Mux.Get(pattern, handler)
	u.addRoute(name, pattern)

	return u
}

func (u *Router) Post(name, pattern string, handler web.HandlerType) *Router {
	u.Mux.Post(pattern, handler)
	u.addRoute(name, pattern)

	return u
}

func (u *Router) Put(name, pattern string, handler web.HandlerType) *Router {
	u.Mux.Put(pattern, handler)
	u.addRoute(name, pattern)

	return u
}

func (u *Router) Delete(name, pattern string, handler web.HandlerType) *Router {
	u.Mux.Delete(pattern, handler)
	u.addRoute(name, pattern)

	return u
}

func (u *Router) Head(name, pattern string, handler web.HandlerType) *Router {
	u.Mux.Head(pattern, handler)
	u.addRoute(name, pattern)

	return u
}

func (u *Router) Trace(name, pattern string, handler web.HandlerType) *Router {
	u.Mux.Trace(pattern, handler)
	u.addRoute(name, pattern)

	return u
}

func (u *Router) Patch(name, pattern string, handler web.HandlerType) *Router {
	u.Mux.Patch(pattern, handler)
	u.addRoute(name, pattern)

	return u
}

func (u *Router) Options(name, pattern string, handler web.HandlerType) *Router {
	u.Mux.Options(pattern, handler)
	u.addRoute(name, pattern)

	return u
}

func (u *Router) addRoute(name, pattern string) {
	u.Routes[name] = &route{
		path: pattern,
	}

	u.Routes[name].compile()
}

func (u *Router) generate(name string, params url.Values, ref int, context *RequestContext) (string, error) {
	queryString := url.Values{}

	if _, ok := u.Routes[name]; !ok {
		return "", errors.New(fmt.Sprintf("Route `%s` does not exist", name))
	}

	for name, values := range params {
		queryString[name] = values
	}

	path := ""
	for _, f := range u.Routes[name].renders {
		v, _ := f(queryString)
		path += v
	}

	sep := "?"

	for name, values := range queryString {
		for _, value := range values {
			path += sep + name + "=" + value

			sep = "&"
		}
	}

	if ref == ABSOLUTE_PATH {
		return path, nil
	}

	if ref == NETWORK_PATH {
		return "/" + path, nil
	}

	if context == nil {
		panic("Unable to generate absolute url without a context")
	}

	return context.Prefix + path, nil
}
