package main

import (
  "fmt"
  "os"
  "strings"
  "github.com/garyburd/redigo/redis"
)

func colorize(msg string) string {
  green := "%{[32m%}"
  yellow := "%{[1;33m%}"
  red := "%{[31m%}"
  reset := "%{[00m%}"

  if strings.Contains(msg, "success") {
    return green + "✓" + reset
  } else if strings.Contains(msg, "pending") {
    return yellow + "…" + reset
  } else if strings.Contains(msg, "error") {
    return red + "ⅹ" + reset
  } else {
    return msg
  }
}

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
    fmt.Printf("%s", colorize(fmt.Sprintf("%s", status)))
  }

  c.Do("PUBLISH", "hub-status", fmt.Sprintf("{ \"ref\": \"%s\", \"dir\": \"%s\" }", ref, dir))
}