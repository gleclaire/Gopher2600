// This file is part of Gopher2600.
//
// Gopher2600 is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Gopher2600 is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Gopher2600.  If not, see <https://www.gnu.org/licenses/>.
//
// *** NOTE: all historical versions of this file, as found in any
// git repository, are also covered by the licence, even when this
// notice is not present ***

package commandline

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jetsetilly/gopher2600/errors"
)

// Validate input string against command defintions
func (cmds Commands) Validate(input string) error {
	return cmds.ValidateTokens(TokeniseInput(input))
}

// ValidateTokens like Validate, but works on tokens rather than an input
// string
func (cmds Commands) ValidateTokens(tokens *Tokens) error {
	cmd, ok := tokens.Peek()
	if !ok {
		return nil
	}
	cmd = strings.ToUpper(cmd)

	for n := range cmds.cmds {
		if cmd == cmds.cmds[n].tag {

			err := cmds.cmds[n].validate(tokens, false)
			if err != nil {
				// preserve FormattedError type
				if _, ok := err.(errors.AtariError); ok {
					return err
				}
				return errors.New(errors.ValidationError, err)
			}

			// if we've reached this point and there are still oustanding
			// tokens in the queue then something has gone wrong.
			if tokens.Remaining() > 0 {
				arg, _ := tokens.Get()

				// special handling for help command
				if cmd == cmds.helpCommand {
					return errors.New(errors.ValidationError, fmt.Sprintf("no help for %s", strings.ToUpper(arg)))
				}

				return errors.New(errors.ValidationError, fmt.Sprintf("unrecognised argument (%s) for %s", arg, cmd))
			}

			return nil
		}
	}

	return fmt.Errorf("unrecognised command (%s)", cmd)
}

func (n *node) validate(tokens *Tokens, speculative bool) error {
	// we cannot do anything useful with a node with an empty tag, but if there
	// is a "next" node then we can move immediately to validation of that node
	// instead.
	//
	// empty tags like this, happen as a result of parsing nested groups
	//
	// a node with an empty tag but no next array (or a next array with too
	// many entries) is an illegal node and should not have been parsed
	if n.tag == "" {
		if n.next == nil || len(n.next) > 1 {
			return errors.New(errors.PanicError, "commandline validation", "illegal empty node")
		}

		// speculatively validate the next node. don't do anything with any
		// error just yet. if there is an error we need to validate against any
		// branches. if there is still no match we cam return the error then

		err := n.next[0].validate(tokens, true)
		match := err == nil

		if !match {
			for bi := range n.branch {
				tokens.Unget()
				if n.branch[bi].validate(tokens, true) == nil {
					match = true
					break // for loop
				}
			}
		}

		if match {
			return nil
		}

		return err
	}

	// make a note of the token queue before we proceed. we'll use this to
	// decide whether to proceed with any repeat loops.
	remainder := tokens.Remainder()

	// get the next token in the token queue
	//
	// in the event of there being no more tokens, then we need to consider
	// whether the current node is optional or not. if it's optional then the
	// validation has passed and we return with no error. if the node is not
	// optional then we return a meaningful and descriptive error.
	tok, ok := tokens.Get()
	if !ok {
		// we treat arguments in the root-group as though they are required
		if n.typ == nodeRequired || n.typ == nodeRoot {
			return fmt.Errorf("%s required", n.nodeVerbose())
		}
		return nil
	}

	// check the current token against the node's tag, using placeholder
	// matching if appropriate.
	//
	// to help we use two boolean variables: match and tentativeMatch
	//
	// match is used to indicate that there is a definite match.
	//
	// tentativeMatch meanwhile is used to indicate that there is a match but
	// there may be a better one for example, the word "foo" matches the %S
	// placeholder but if another branch expects the exact argument "foo" then
	// that would be a better match.

	match := false
	tentativeMatch := false

	switch n.tag {
	case "%N":
		// normalise hex notation and update token
		if tok[0] == '$' {
			tok = fmt.Sprintf("0x%s", tok[1:])
		}
		tokens.Update(tok)

		_, e := strconv.ParseInt(tok, 0, 32)
		match = e == nil

	case "%P":
		_, e := strconv.ParseFloat(tok, 32)
		match = e == nil

		// I originally thought that an error message describing how the
		// argument is "not a number" or "not a float" would be helpful but in
		// practice, it wasn't as useful as you might expect. for instance if
		// we had the template:
		//
		// WATCH (READ|WRITE) %N
		//
		// the command:
		//
		// WATCH ANY 0x80
		//
		// would result in an error message like "ANY is not a number", because
		// ANY does not match the optional group. I think this is misleading.
		//
		// with a bit of work we could craft the validation algorithm to notice
		// that "0x80" does match the %N argument and so ANY was supposed to be
		// an attempt at the optional argument, but that's a lot more work.
		// however, for now, we've opted to resond to all bad arguments with a
		// catch-all "unrecognised argument" message (see below).

	case "%S":
		tentativeMatch = true
		match = n.branch == nil

	case "%F":
		// not checking for file existence
		tentativeMatch = true
		match = n.branch == nil

	default:
		// case insensitive matching. n.tag should have been normalised
		// already.
		tok = strings.ToUpper(tok)
		match = tok == n.tag

		// update token with normalised string
		if match {
			tokens.Update(tok)
		}
	}

	// if input doesn't match this node we need to check branches. we may well
	// have a tentative match at this point but we need to put that to one side
	// until we've checked all other options.
	if !match && n.branch != nil {
		for bi := range n.branch {
			tokens.Unget()

			if n.branch[bi].validate(tokens, true) == nil {
				return nil
			}
		}

		// there's no explicit match in any of the matches. if we've
		// encountered a tentative match however, we can use that
		match = tentativeMatch
	}

	if !match {
		err := fmt.Errorf("unrecognised argument (%s)", tok)

		// there's still no match but the speculative flag means we were half
		// expecting it. return error without further consideration
		//
		// the fact that this is a speculative validation means that the error
		// may well be ignored; but that's not a decision to make here
		if speculative {
			return err
		}

		// if the node is not optional then failing to match is a definite
		// error. return the previously prepared error back to the caller
		if n.typ != nodeOptional {
			return err
		}

		// the node is optional so we can simply carry on to the "next" nodes.
		// however, because the current token did not match we'll need to
		// examine it again
		//
		// The Unget() function "pushes" the current token back onto the queue.
		tokens.Unget()
	}

	// check nodes that follow on from the current node
	for ni := range n.next {
		err := n.next[ni].validate(tokens, false)
		if err != nil {
			return err
		}
	}

	// no more nodes in the next array. move to the repeat node if there is one
	// and if the tokens queue has changed since the beginning of this
	// function.
	if n.repeat != nil && remainder != tokens.Remainder() {
		err := n.repeat.validate(tokens, false)
		if err != nil {
			return err
		}
	}

	return nil
}
