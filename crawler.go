//go:generate protoc --go_out=. host-assistant-proto/listing.proto

package main

import (
  "fmt"
  "os"
  "github.com/streadway/amqp"
  "github.com/golang/protobuf/proto"
  "github.com/tjurkiewicz/airbnb-crawler/host-assistant-proto"
  "github.com/tjurkiewicz/airbnb-api-client"
)

const BASE_URL string = "https://www.airbnb.com/api/v1/listings/"

func ReadListing(id int64) {
  cli := client.AirBNB{ApiKey: os.Getenv("AIRBNB_KEY")}
  listingResponse, errorResponse, err := cli.ReadListing(fmt.Sprintf("%d", id))
  failOnError(err, "airbnb.readlisting")

  fmt.Println(listingResponse, errorResponse)
}

func main() {
  amqpConnection, err := amqp.Dial(os.Getenv("AMQP_URL"))
  failOnError(err, "amqp.dial")
  defer amqpConnection.Close()

  ch, err := amqpConnection.Channel()
  failOnError(err, "amqp.channel")
  defer ch.Close()

  messages, err := ch.Consume(os.Getenv("AMQP_QUEUE"), "", true, false, false, false, nil)
  failOnError(err, "amqp.queue.consume")

  for message := range messages {
    request := &listing.ListingRequest{}
    err = proto.Unmarshal(message.Body, request)
    failOnError(err, "proto.unmarshal")

    go ReadListing(request.Id)
  }
}

