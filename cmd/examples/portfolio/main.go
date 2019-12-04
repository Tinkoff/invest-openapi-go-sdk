package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
)

func main() {
	token := flag.String("token", "", "OpenApi Token")

	flag.Parse()

	logger := log.New(os.Stdout, "[invest-openapi-go-sdk]", log.LstdFlags)

	restClient := sdk.NewRestClient(*token)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	portfolio, err := restClient.Portfolio(ctx)
	if err != nil {
		logger.Fatalln(err)
	}

	for _, position := range portfolio.Positions {
		instrument, err := restClient.SearchInstrumentByFIGI(ctx, position.FIGI)

		if err != nil {
			logger.Printf("Cannot find position in /market/search/figi=%s", position.FIGI)
			break
		}

		logger.Printf("Position %s (%s): %f", position.FIGI, instrument.Name, position.Balance)
	}
}
