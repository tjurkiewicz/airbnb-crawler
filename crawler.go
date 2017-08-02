package main

import (
  "fmt"
  "log"
  "os"
  "github.com/dghubble/sling"
  "github.com/streadway/amqp"
)

const BASE_URL string = "https://www.airbnb.com/api/v1/listings/"

type Params struct {
  Key string `url:"key,omitempty"`
}

type ErrorResponse struct {
  ErrorCode    int    `json:"error_code"`
  Error        string `json:"error"`
  ErrorMessage string `json:"error_message"`
}

type ListingResponse struct {
  Listing struct {
    Id             int64   `json:"id"`
    City           string  `json:"city"`
    UserId         int64   `json:"user_id"`
    Latitude       float64 `json:"lat"`
    Longitude      float64 `json:"lng"`
    Bathrooms      float64 `json:"bathrooms"`
    Bedrooms       float64 `json:"bedrooms"`
    Beds           float64 `json:"beds"`
    PersonCapacity int     `json:"person_capacity"`
    CountryCode    string  `json:"country_code"`
  } `json:"listing"`
}

func failOnError(err error, key string) {
  if err != nil {
    log.Fatalf("%s: %s", key, err)
  }
}


func ReadListing(id string) {
  params := &Params{Key: os.Getenv("AIRBNB_KEY")}
  listingResponse := new(ListingResponse)
  errorResponse := new(ErrorResponse)

  _, err := sling.New().Base(BASE_URL).Path(id).QueryStruct(params).Receive(listingResponse, errorResponse)
  failOnError(err, "sling.receive")

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
    id := string(message.Body[:])
    go ReadListing(id) 
  }
}

