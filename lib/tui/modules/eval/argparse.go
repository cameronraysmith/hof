package eval

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hofstadter-io/hof/lib/tui"
)

// this function helps parse args and context into richer information
// so that we can handle different kinds of input coming from different places
// and then have more consistent input to the components consuming the inputs
func enrichContext(context map[string]any) (map[string]any) {
	// tui.Log("trace", fmt.Sprintf("enrichContext.BEG: %# v", context))

	// ensure non-nil context
	if context == nil {
		context = make(map[string]any)
	}

	// extract and record original args
	args := []string{}
	if _args, ok := context["args"]; ok {
		args = _args.([]string)
	}
	context["orig-args"] = args

	//// so we can special case create a default element for naked eval
	// hadEval := false
	if len(args) > 0 && args[0] == "eval" {
		// hadEval = true
		args = args[1:]
	}

	// TODO, look to see if there is a ' -- '
	// and if exists, always parse here until it is reached
	// this will allow a user to be more exact and get better error messages

	// find any modifiers
	for len(args) > 0 {
		tok := args[0]
		switch tok {

		//
		// top-level eval commands
		//
		case 
			"preview",
			"save",
			"list",
			"show",
			"load":
			context["action"] = tok 

		//
		// actions
		//
		// let's start using some `cmd.sub` syntax for these
	 	case
			"insert",
			"ins":
			context["action"] = "insert"
		case
			"delete",
			"del":
			context["action"] = "delete"

		// item update
		case
			"reload",
			"refresh":
			context["action"] = "reload"
		case
			"update":
			context["action"] = "update"  // probably the default?

		// this should probably be the new default
		case 
			"push",
			"export",
			"write",
			"set.value",
			"set.scope",
			// "set.text",

			"set.name",
			"set.item.name",
			"set.panel.name":
			context["action"] = tok 

		case
			"connect",
			"C", "conn":
			context["action"] = "connect"

		// not handled yet
		case "nav.left", "nl":
			context["action"] = "nav.left"
		case "nav.up", "nu":
			context["action"] = "nav.up"
		case "nav.down", "nd":
			context["action"] = "nav.down"
		case "nav.right", "nr":
			context["action"] = "nav.right"

		//
		// items / widgets / pane
		//

		// default when nothing
		case "help":
			context["item"] = "help"
			// context["action"] = "insert"  // probably the default?
		// dual-pane eval'r (default when doing eval things)
		case "play":
			context["item"] = "play"
			context["action"] = "update"  // probably the default?
		// value viewer
		case "view":
			context["item"] = "view"
			context["action"] = "update"  // probably the default?

		// should this be handled lower too?
		// we might want a more general 
		// sources?
		case "sh", "bash":
			context["source"] = "bash"

		//
		// what's acted upon, location within
		//
		
		// cue item action targets
		case "S", "scope":
			context["target"] = "scope"
		case "V", "value":
			context["target"] = "value"

		// default, the current focused item
		case "s", "self":
			context["where"] = "self"

		// eventually, which dashboard
		// here, there, new, extend

		// panel relative
		case "h", "head", "f", "first":
			context["where"] = "head"
		case "l", "last", "t", "tail", "e", "end":
			context["where"] = "last"
		case "b", "before", "p", "prev":
			context["where"] = "before"
		case "a", "after", "n", "next":
			context["where"] = "next"

		// for within, to be passed down to current panel
		case "K":
			context["where"] = "top"
		case "k":
			context["where"] = "above"
		case "j":
			context["where"] = "below"
		case "J":
			context["where"] = "bottom"

		// default is the browser
		// setup the play with the eval as scope

		// these are more options, should be handled lower down, with item local parser
		//case "C", "cue":
		//  context["mode"] = "cue"
		//case "yaml":
		//  context["mode"] = "yaml"
		//case "json":
		//  context["mode"] = "json"


		// UPDATE, these probably become ALT-... keys
		// connect (start and end modes)
		// set colors while in progress
		// or also support passing all at once
		// also support mouse clicking

		// user intention, end of args
		case "--":
			goto argsDone

		default:
			// is it an int?
			if i, err := strconv.Atoi(tok); err == nil {
				context["where"] = "index"
				// maybe make this a list, to support x,y row,col or more indexes
				// or maybe make them CSV each, like pos ~ 1,2
				// but maybe easier as a slice of ints that the function can interpret as it wishes, ordered like
				// we'd have a hard time assigning them to different keys here
				context["index"] = i
			} else {
				// break doens't break here, but we want to stop processing args here
				goto argsDone
			}
		}
		args = args[1:]
	}
argsDone:

	// handle some special cases
	if len(args) > 0 && strings.HasPrefix(args[0], "http") {
		context["source"] = "http"
		//context["from"] = args[0]
		//args = args[1:]
	}

	// make sure we update the context args
	context["args"] = args

	tui.Log("trace", fmt.Sprintf("enrichContext.END: %# v", context))
	return context
}
