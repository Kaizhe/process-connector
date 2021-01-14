package main

import (
	"github.com/kaizhe/proc-connector/connector"
	"github.com/kaizhe/proc-connector/enricher"
	"github.com/kaizhe/proc-connector/pkg/types"
)

func main() {
	messageChan := make(chan types.Message, 1)

	pc := connector.NewProcessConnector()

	e := enricher.NewEnricher()

	go e.Enrich(messageChan)

	err := pc.Listen(messageChan)
	if err != nil {
		panic(err)
	}
}
