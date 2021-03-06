package dfa

import (
  "fmt"
  "regexp"
  "io/ioutil"
  "strings"

  "github.com/BurntSushi/toml"
  "github.com/aarzilli/golua/lua"
)

type DFAConfig struct {
  Start                 string      `toml:"start"`
  HasNoMatchState       bool        `toml:"hasNoMatch"`
  NoMatchState          string      `toml:"noMatch"`
  GlobalPostTransitions [][2]string `toml:"globalPostTransitions"`
}

type DFAState struct {
  Accept      bool        `toml:"accept"`
  Transitions [][2]string `toml:"transitions"`
}

type DFA struct {
  L      *lua.State
  Config DFAConfig           `toml:"dfa"`
  States map[string]DFAState `toml:"state"`
}

var validTypes = map[string]bool{
  "s:":true, // string
  "r:":true, // regex
  "f:":true, // lua func
}

func ParseDFAFromFile(path string) (*DFA, error) {
  dfaStr,err := ioutil.ReadFile(path)
  if err != nil {
    return nil,err
  }

  return ParseDFA(string(dfaStr))
}

func ParseDFA(dfaToml string) (*DFA, error) {
  var dfa DFA
  if _, err := toml.Decode(dfaToml, &dfa); err != nil {
    return nil,err
  }


  if _, ok := dfa.States[dfa.Config.Start]; !ok {
    return nil,fmt.Errorf("Start state %v was not found in the DFA.", dfa.Config.Start)
  }

  var acceptFound = false
  var rejectFound = false

  for state,data := range dfa.States {
    acceptFound = acceptFound || data.Accept
    rejectFound = rejectFound || !data.Accept

    for i,trans := range data.Transitions {
      if _, ok := validTypes[trans[0][:2]]; !ok {
        return nil,fmt.Errorf("Transition type %v in transition %v from %v to %v was not found in the DFA.", trans[0], i, state, trans[1])
      }
      if _, ok := dfa.States[trans[1]]; !ok {
        return nil,fmt.Errorf("Destination state in transition %v from %v to %v was not found in the DFA.", i, state, trans[1])
      }
    }
  }

  if !acceptFound {
    fmt.Println("Warning: no accept state found in DFA.")
  }

  if !rejectFound {
    fmt.Println("Warning: no reject state found in DFA.")
  }

  for i,trans := range dfa.Config.GlobalPostTransitions {
    if _, ok := validTypes[trans[0][:2]]; !ok {
      return nil,fmt.Errorf("Transition type %v in global transition %v to %v was not found in the DFA.", trans[0], i, trans[1])
    }
    if _, ok := dfa.States[trans[1]]; !ok {
      return nil,fmt.Errorf("Destination state in global transition %v to %v was not found in the DFA.", i, trans[1])
    }
  }

  if _, ok := dfa.States[dfa.Config.NoMatchState]; dfa.Config.HasNoMatchState && !ok {
      return nil,fmt.Errorf("No Match State %v does not exist.", dfa.Config.NoMatchState)
  }

  dfa.L = lua.NewState()
  dfa.L.OpenLibs()
  dfa.L.DoFile("defaultlib.lua")

  return &dfa,nil
}

func (dfa *DFA) Close() {
  if dfa.L != nil {
    dfa.L.Close()
  }
}

func (dfa *DFA) AddLua(file string) {
  dfa.L.DoFile(file)
}

func (dfa *DFA) ApplyLua(funcname string, inputA string, inputB string) bool {
  dfa.L.GetField(lua.LUA_GLOBALSINDEX, funcname)
  dfa.L.PushString(inputA)
  dfa.L.PushString(inputB)
  dfa.L.Call(2,1)
  rtn := dfa.L.ToBoolean(-1)
  dfa.L.Pop(-1)
  return rtn
}

func (dfa *DFA) ApplyDFA(inputs []string) (bool, error) {
  var curState = dfa.Config.Start

  LoopTop:
  for _,input := range inputs {
    for _,trans := range dfa.States[curState].Transitions {
      switch(trans[0][:2]) {
      case "s:":
        if input == trans[0][2:] {
          curState = trans[1]
          continue LoopTop
        }
      case "r:":
        matched, err := regexp.MatchString(trans[0][2:], input)
        if err != nil {
          return false, err
        }
        if matched {
          curState = trans[1]
          continue LoopTop
        }
      case "f:":
        split := strings.Split(trans[0],":")
        if dfa.ApplyLua(split[1], split[2], input) {
          curState = trans[1]
          continue LoopTop
        }
      }
    }

    for _,trans := range dfa.Config.GlobalPostTransitions {
      switch(trans[0][:2]) {
      case "s:":
        if input == trans[0][2:] {
          curState = trans[1]
          continue LoopTop
        }
      case "r:":
        matched, err := regexp.MatchString(trans[0][2:], input)
        if err != nil {
          return false, err
        }
        if matched {
          curState = trans[1]
          continue LoopTop
        }
      case "f:":
        split := strings.Split(trans[0],":")
        if dfa.ApplyLua(split[1], split[2], input) {
          curState = trans[1]
          continue LoopTop
        }
      }
    }

    if dfa.Config.HasNoMatchState {
      curState = dfa.Config.NoMatchState
    }
  }

  return dfa.States[curState].Accept,nil
}
