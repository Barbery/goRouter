package goRouter

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// http method
const (
	CONNECT = "CONNECT"
	DELETE  = "DELETE"
	GET     = "GET"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
	POST    = "POST"
	PUT     = "PUT"
	TRACE   = "TRACE"
)

//mime-types
const (
	applicationJson = "application/json"
	applicationXml  = "application/xml"
	textXml         = "text/xml"
)

type Mux struct {
	beforeMatch   http.HandlerFunc
	afterMatch    http.HandlerFunc
	beforeExecute http.HandlerFunc
	afterExecute  http.HandlerFunc
	routes        map[string]map[string][]*route
}

type route struct {
	pattern *regexp.Regexp
	params  []string
	handler http.HandlerFunc
}

var muxInstance = &Mux{
	beforeMatch:   func(rw http.ResponseWriter, req *http.Request) {},
	afterMatch:    func(rw http.ResponseWriter, req *http.Request) {},
	beforeExecute: func(rw http.ResponseWriter, req *http.Request) {},
	afterExecute:  func(rw http.ResponseWriter, req *http.Request) {},
	routes:        make(map[string]map[string][]*route),
}

var config = map[string]string{
	//路由匹配规则
	"matchReg": `^%s$`,
	//默认参数匹配规则
	"defaultParamsReg": `([^/]+)`,
	//查找参数规则
	"findParamsReg": `(:\w+)`,
	//处理没有带正则规则的参数规则
	"processParamsReg": `:\w+`,
	//处理带正则规则的参数规则
	"processParamsWithReg": `:\w+(\(.*?\))`,
}

func init() {
}

func GetMuxInstance() *Mux {
	return muxInstance
}

func (m *Mux) Get(pattern string, handler http.HandlerFunc) {
	m.AddRoute(pattern, handler, GET)
}

func (m *Mux) Post(pattern string, handler http.HandlerFunc) {
	m.AddRoute(pattern, handler, POST)
}

func (m *Mux) Put(pattern string, handler http.HandlerFunc) {
	m.AddRoute(pattern, handler, PUT)
}

func (m *Mux) Delete(pattern string, handler http.HandlerFunc) {
	m.AddRoute(pattern, handler, DELETE)
}

func (m *Mux) AddRoute(pattern string, handler http.HandlerFunc, method string) {
	pattern = strings.TrimRight(pattern, `/`)
	parts := strings.Split(pattern, `/`)
	var prefix []string
	for _, part := range parts {
		if strings.Index(part, ":") == -1 {
			prefix = append(prefix, part)
		} else {
			break
		}
	}

	//找出所有需要匹配的参数
	findParamReg := regexp.MustCompile(config["findParamsReg"])
	params := findParamReg.FindAllString(pattern, -1)

	//先处理带正则规则限定的参数
	replaceReg := regexp.MustCompile(config["processParamsWithReg"])
	pattern = replaceReg.ReplaceAllString(pattern, "$1")

	//没有正则限定的参数，使用默认正则规则来匹配
	replaceReg = regexp.MustCompile(config["processParamsReg"])
	pattern = replaceReg.ReplaceAllString(pattern, config["defaultParamsReg"])

	regex := regexp.MustCompile(fmt.Sprintf(config["matchReg"], pattern))

	if _, exist := m.routes[method]; !exist {
		m.routes[method] = map[string][]*route{}
	}

	prefixUrl := strings.Join(prefix, `/`)
	if _, exist := m.routes[method][prefixUrl]; !exist {
		m.routes[method][prefixUrl] = []*route{}
	}

	m.routes[method][prefixUrl] = append(m.routes[method][prefixUrl], &route{
		pattern: regex,
		handler: handler,
		params:  params,
	})
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.beforeMatch(w, r)

	handler, isMatch := m.match(strings.TrimRight(r.URL.Path, `/`), r)
	if !isMatch {
		http.NotFound(w, r)
		return
	}
	m.afterMatch(w, r)
	m.beforeExecute(w, r)
	handler(w, r)
	m.afterExecute(w, r)
}

func (m *Mux) match(requestPath string, r *http.Request) (http.HandlerFunc, bool) {
	paths := strings.Split(requestPath, `/`)

	if _, ok := m.routes[r.Method]; !ok {
		return nil, false
	}

	i := len(paths)
	var path string
	for ; i > 0; i-- {
		path = strings.Join(paths[:i], `/`)
		if routes, ok := m.routes[r.Method][path]; ok {
			for _, route := range routes {
				if !route.pattern.MatchString(requestPath) {
					continue
				}

				//whether need to match the parameters
				if len(route.params) > 0 {
					matches := route.pattern.FindStringSubmatch(requestPath)
					if len(matches) < 2 || len(matches[1:]) != len(route.params) {
						// panic("Parameters do not match")
						return nil, false
					}

					values := r.URL.Query()
					for i, match := range matches[1:] {
						values.Add(route.params[i], match)
					}

					//reassemble query params and add to RawQuery
					r.URL.RawQuery = url.Values(values).Encode() + "&" + r.URL.RawQuery
				}

				return route.handler, true
			}
		}
	}

	return nil, false
}
