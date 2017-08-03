package main

import (
  "log"
)

func failOnError(err error, key string) {
  if err != nil {
    log.Fatalf("%s: %s", key, err)
  }
}

