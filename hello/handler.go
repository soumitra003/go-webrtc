package hello

import "net/http"

func (h *ModuleHello) sayHelloHandler(res http.ResponseWriter, req *http.Request) {
	_, _ = res.Write([]byte("Hello World"))
}
