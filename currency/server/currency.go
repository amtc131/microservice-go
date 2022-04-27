package server

import (
	"context"
	"io"
	"time"

	"github.com/amtc131/microservice-go/currency/data"
	protos "github.com/amtc131/microservice-go/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Currency struct {
	rates         *data.ExchangeRates
	log           hclog.Logger
	subscriptions map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
}

func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	c := &Currency{r, l, make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest)}
	go c.handleUpdates()

	return c
}

func (c *Currency) handleUpdates() {
	ru := c.rates.MonitorRates(5 * time.Second)
	for range ru {
		c.log.Info("Got update rates")
		for k, v := range c.subscriptions {
			for _, rr := range v {
				r, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
				if err != nil {
					c.log.Error("Unable to get update rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}

				err = k.Send(
					&protos.StreamingRateResponse{
						Message: &protos.StreamingRateResponse_RateReponse{
							RateResponse: &protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: r},
						},
					},
				)
				if err != nil {
					c.log.Error("Unable to send updated rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}
			}
		}
	}

}

// GetRate implements the CurrencyServer GetRate method and returns the currency exchange rate
// for the two given currencies
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

	if rr.Base == rr.Destination {
		err := status.Newf(
			codes.InvalidArgument,
			"Base currency %s can not be the same as the destination currency %s",
			rr.Base.String(),
			rr.Destination.String(),
		)

		err, wde := err.WithDetails(rr)
		if wde != nil {
			return nil, wde
		}

		return nil, err.Err()
	}

	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: rate}, nil
}

// SubscribeRates  implements the gRPC bidirection streaming method for the server
func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {

	for {
		rr, err := src.Recv()
		if err == io.EOF {
			c.log.Info("Client hat closed connection")
			break
		}
		if err != nil {
			c.log.Error("Unable to read from client", "error", err)
			return err
		}

		c.log.Info("Handle client request", "request", rr)
		rrs, ok := c.subscriptions[src]
		if !ok {
			rrs = []*protos.RateRequest{}
		}

		//check that subscribtion does not exist
		var validationError *status.Status
		for _, v := range rrs {
			if v.Base == rr.Base && v.Destination == rr.Destination {

				// subscription exists return errors
				s := status.Newf(
					codes.AlreadyExists,
					"Unable to subscribe for currency as subscription already exists")

				// add the original request as metadata
				s, err = s.WithDetails(err)
				if err != nil {
					c.log.Error("Unable to add metadate to error", "error", err)
					break
				}

				break

			}
		}

		//if validation error return error and continue
		if validationError != nil {
			src.Send(
				&protos.StreamingRateResponse{
					Message: &protos.StreamingRateResponse_Error{
						Error: s.Proto(),
					},
				},
			)
			continue
		}

		// all ok
		rrs = append(rrs, rr)
		c.subscriptions[src] = rrs
	}

	return nil
}
