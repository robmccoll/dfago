package main

import (
  "fmt"
  "os"

  "github.com/robmccoll/dfago"
)

func main() {
  if len(os.Args) < 3 {
    fmt.Println("Usage: %v <dfatomlfile> input0 input1 ... inputn", os.Args[0])
    return
  }

  dfa,err := dfa.ParseDFAFromFile(os.Args[1])

  if err != nil {
    fmt.Println(err.Error())
    return
  }

  result,err := dfa.ApplyDFA(os.Args[2:])

  if err != nil {
    fmt.Println(err.Error())
    return
  }

  fmt.Println(result)
  if result {
    os.Exit(1)
  } else {
    os.Exit(0)
  }
}
