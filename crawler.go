package main

import (
  "fmt"
  "os"
  "github.com/dghubble/sling"
  "github.com/streadway/amqp"
  "github.com/tjurkiewicz/airbnb-crawler/host-assistant-proto" 
)

const BASE_URL string = "https://www.airbnb.com/api/v1/listings/"

func ReadListing(id string) {
  params := &Params{Key: os.Getenv("AIRBNB_KEY")}
  listingResponse := new(ListingResponse)
  errorResponse := new(ErrorResponse)

  _, err := sling.New().Base(BASE_URL).Path(id).QueryStruct(params).Receive(listingResponse, errorResponse)
  failOnError(err, "sling.receive")

  fmt.Println(listingResponse, errorResponse)

  _ = listing.ListingRequest{Id: 1}
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
    id := string(message.Body[:])
    go ReadListing(id) 
  }
}

