package browser

import (
	"b2t_helpdesk/httpserver"
	"b2t_helpdesk/injector"
)

func StartHTTPServer(di *injector.Injector) {
	outbondserver := httpserver.NewServer("0.0.0.0:8879", "Outbound Listener")
	initRoute(outbondserver.Router, di)
	outbondserver.Run()
	di.WG.Done()
}
