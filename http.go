package hubbard

import "http"
import "os"

func getParameter(req *http.Request, name string) string {
	values := req.Form[name]
	if len(values) == 0 {
		panic(os.NewError("missing parameter: " + name))
	}
	if len(values) > 1 {
		panic(os.NewError("too many values for parameter: " + name))
	}
	return values[0]
}
