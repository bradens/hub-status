package main

import (
  "fmt"
  "os"
  "github.com/garyburd/redigo/redis"
)

func main() {
  ref := os.Args[1]
  dir := os.Args[2]

  redisPool := redis.NewPool(func() (redis.Conn, error) {
    c, err := redis.Dial("tcp", "localhost:6379")
    if err != nil {
      return nil, err
    }
    return c, err
  }, 10)

  defer redisPool.Close()

  c := redisPool.Get()
  defer c.Close()

  key := fmt.Sprintf("hub-status:%s", ref)
  status, _ := c.Do("GET", key)

  if status != nil {
    fmt.Printf("%s", status)
  }

  c.Do("PUBLISH", "hub-status", fmt.Sprintf("{ \"ref\": \"%s\", \"dir\": \"%s\" }", ref, dir))
}