dfago
=====

A simple implementation of Discrete Finite Automata written in Go. Uses TOML for DFA syntax. See example.dfa.

States in the DFA are listed as:

    [state.state_name]
    acccept = true | false
    transitions [
      ... array of transitions which are...
      ["s:exact_string_match", "destination_state_name"],
      ["r:.*regex_match", "destination_state_name"],
      ["f:lua_function:argument1", "destination_state_name"]
    ]


The top level [dfa] object will include the name of the starting state and a few optional fields described below.
Destination states must be valid names.  In a given state, the current input will be matched against
each transition until the first match.  Matches can be exact strings (prefixed with 's:'), 
regular expression matches (prefixed with 'r:'), or Lua function callse (prefixed with 'f:' and formatted as
'f:function\_name:argument\_string).  If no match is made, the DFA will remain in the 
current state by default.  When all inputs have been consumed, the accept value of the current state is 
returned.

Lua functions must match the signature:

    function (str1, str2)
      return bool
    end

These will be called with the argument embedded in the transition as str1 and the input as str2. The functions 
are pulled from defaultlib.lua and any files passed to dfa.AddLua(). If the function returns true, the 
transition is considered matched to the input and will be taken. Otherwise, the DFA will continue searching for
transition matches.

    [dfa]
    start = "start_state_name"
    hasNoMatch = false | true
    noMatch = "destination_state_if_no_match"
    globalPostTransitions = [
      ... array of transitions which are applied in every state 
      if no transitions in the current state are matched...
      ["s:exact_string_match", "destination_state_name"],
      ["r:.*regex_match", "destination_state_name"],
    ]
    
The rundfa binary in the rundfa folder can be used to directly run TOML formatted DFA files on the commant line.
It is a thin wrapper around the dfago library.

To build:

    cd rundfa
    go build
    
To run:

    ./rundfa <toml_file> <input0> <input1> ... <inputN>


A Note on Building
------------------

With the addition of Lua support (via [aarzilli's excellent golua](https://github.com/aarzilli/golua)), you will need 
to have the Lua shared library installed to build.  See his page for additional details. Once linked, this is not required
for running.
