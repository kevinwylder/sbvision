package frontend

import (
	"fmt"
	"net/http"
	"os"
	"path"
)

type frontend struct {
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
	return http.FileServer(&frontend{
		dist:     http.Dir(distFolder),
		react:    http.Dir(reactFolder),
		reactDOM: http.Dir(reactDOMFolder),
	}), nil
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
