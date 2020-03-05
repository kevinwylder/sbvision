package frontend

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
)

type frontend struct {
	files    http.Handler
	dist     http.Dir
	react    http.Dir
	reactDOM http.Dir
}

// ServeFrontend handles static files from the frontend folder
func ServeFrontend(dir string) (http.Handler, error) {
	distFolder := path.Join(dir, "dist")
	reactFolder := path.Join(dir, "node_modules", "react", "umd")
	reactDOMFolder := path.Join(dir, "node_modules", "react-dom", "umd")
	if _, err := os.Stat(distFolder); os.IsNotExist(err) {
		return nil, fmt.Errorf("\n\tCannot find dist folder in %s. Try `yarn build`", dir)
	}
	if _, err := os.Stat(reactFolder); os.IsNotExist(err) {
		return nil, fmt.Errorf("\n\tCannot find node_modules/react/umd folder in %s. Try `yarn`", dir)
	}
	if _, err := os.Stat(reactDOMFolder); os.IsNotExist(err) {
		return nil, fmt.Errorf("\n\tCannot find node_modules/react-dom/umd folder in %s. Try `yarn`", dir)
	}
	f := &frontend{
		dist:     http.Dir(distFolder),
		react:    http.Dir(reactFolder),
		reactDOM: http.Dir(reactDOMFolder),
	}
	f.files = http.FileServer(f)
	return f, nil
}

func (f *frontend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// catch and redirect all routing requests
	if strings.HasPrefix(r.URL.Path, "/video/") {
		r.URL.Path = "/"
	}
	switch r.URL.Path {
	case "/videos":
		r.URL.Path = "/"
	case "/rotations":
		r.URL.Path = "/"
	case "/explore":
		r.URL.Path = "/"
	case "/api-docs":
		r.URL.Path = "/"
	}
	f.files.ServeHTTP(w, r)
}

func (f *frontend) Open(path string) (http.File, error) {
	file, err := f.dist.Open(path)
	if err == nil {
		return file, nil
	}
	file, err = f.react.Open(path)
	if err == nil {
		return file, nil
	}
	return f.reactDOM.Open(path)
}
