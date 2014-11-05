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
    ]


The top level [dfa] object will include the name of the starting state.
Destination states must be valid names.  In a given state, the current input will be matched against
each transition until the first match.  Matches can either be exact strings (prefixed with 's:') or
regular expression matches (prefixed with 'r:').  If no match is made, the DFA will remain in the 
current state.  When all inputs have been consumed, the accept value of the current state is 
returned.

