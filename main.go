package main

import "github.com/kaizhe/proc-connector/connector"

func main() {
	pc := connector.NewProcessConnector()

	err := pc.Listen()

	if err != nil {
		panic(err)
	}
}

