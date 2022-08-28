package main

import (
	cs "gitlab.com/cronolabs/cable/internal/proxy"
	web "gitlab.com/cronolabs/cable/internal/server"
)

func main() {
	cs.InitTracker()
	web.RunServer()
}
