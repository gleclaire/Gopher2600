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

package commandline_test

import (
	"fmt"
	"testing"

	"github.com/jetsetilly/gopher2600/debugger/terminal/commandline"
)

func TestValidation_required(t *testing.T) {
	var cmds *commandline.Commands
	var err error

	cmds, err = commandline.ParseCommandTemplate([]string{"TEST [arg]"})
	if err != nil {
		t.Fatalf("%s", err)
	}

	err = cmds.Validate("TEST arg foo")
	if err == nil {
		t.Errorf("matches but shouldn't")
	} else {
		fmt.Println(err)
	}

	err = cmds.Validate("TEST arg")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	err = cmds.Validate("TEST")
	if err == nil {
		t.Errorf("matches but shouldn't")
	} else {
		fmt.Println(err)
	}
}

func TestValidation_optional(t *testing.T) {
	var cmds *commandline.Commands
	var err error

	cmds, err = commandline.ParseCommandTemplate([]string{"TEST (arg)"})
	if err != nil {
		t.Fatalf("%s", err)
	}

	err = cmds.Validate("TEST")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	err = cmds.Validate("TEST arg")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	err = cmds.Validate("TEST arg foo")
	if err == nil {
		t.Errorf("matches but shouldn't")
	} else {
		fmt.Println(err)
	}

	err = cmds.Validate("TEST foo")
	if err == nil {
		t.Errorf("matches but shouldn't")
	} else {
		fmt.Println(err)
	}
}

func TestValidation_branchesAndNumeric(t *testing.T) {
	var cmds *commandline.Commands
	var err error

	cmds, err = commandline.ParseCommandTemplate([]string{"TEST (arg [%N]|foo)"})
	if err != nil {
		t.Fatalf("%s", err)
	}

	err = cmds.Validate("TEST")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	err = cmds.Validate("TEST arg")
	if err == nil {
		t.Errorf("matches but shouldn't")
	} else {
		fmt.Println(err)
	}

	// numeric argument matching
	err = cmds.Validate("TEST arg 10")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	// failing a numeric argument match
	err = cmds.Validate("TEST arg bar")
	if err == nil {
		t.Errorf("matches but shouldn't")
	} else {
		fmt.Println(err)
	}

	// ---------------

	cmds, err = commandline.ParseCommandTemplate([]string{"TEST (arg|foo) %N"})
	if err != nil {
		t.Fatalf("%s", err)
	}

	err = cmds.Validate("TEST arg")
	if err == nil {
		t.Errorf("matches but shouldn't")
	} else {
		fmt.Println(err)
	}

	err = cmds.Validate("TEST arg 10")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	err = cmds.Validate("TEST 10")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}
}

func TestValidation_deepBranches(t *testing.T) {
	var cmds *commandline.Commands
	var err error

	// retry numeric argument matching but with an option for a specific string
	cmds, err = commandline.ParseCommandTemplate([]string{"TEST (arg [%N|bar]|foo)"})
	if err != nil {
		t.Fatalf("%s", err)
	}

	err = cmds.Validate("TEST arg bar")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	err = cmds.Validate("TEST arg foo")
	if err == nil {
		t.Errorf("matches but shouldn't")
	} else {
		fmt.Println(err)
	}
}

func TestValidation_tripleBranches(t *testing.T) {
	var cmds *commandline.Commands
	var err error

	cmds, err = commandline.ParseCommandTemplate([]string{"TEST (arg|foo|bar) wibble"})
	if err != nil {
		t.Fatalf("%s", err)
	}

	err = cmds.Validate("TEST foo wibble")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	err = cmds.Validate("TEST bar wibble")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	err = cmds.Validate("TEST wibble")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}
}

func TestValidation_doubleArgs(t *testing.T) {
	var cmds *commandline.Commands
	var err error

	cmds, err = commandline.ParseCommandTemplate([]string{"TEST (nug nog|egg|cream) (tug)"})
	if err != nil {
		t.Fatalf("%s", err)
	}

	err = cmds.Validate("TEST nug nog")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	err = cmds.Validate("TEST egg tug")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	err = cmds.Validate("TEST nug nog tug")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	// ---------------

	cmds, err = commandline.ParseCommandTemplate([]string{"TEST (egg|fog|nug nog|big) (tug)"})
	if err != nil {
		t.Fatalf("%s", err)
	}

	err = cmds.Validate("TEST nug nog")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	err = cmds.Validate("TEST fog tug")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	err = cmds.Validate("TEST nug nog tug")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}
}

func TestValidation_filenameFirstArg(t *testing.T) {
	var cmds *commandline.Commands
	var err error

	cmds, err = commandline.ParseCommandTemplate([]string{"TEST [%F|foo [wibble]|bar]"})
	if err != nil {
		t.Fatalf("%s", err)
	}

	err = cmds.Validate("TEST foo wibble")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}
}

func TestValidation_singluarOption(t *testing.T) {
	var cmds *commandline.Commands
	var err error

	cmds, err = commandline.ParseCommandTemplate([]string{"SCRIPT [RECORD (REGRESSION) [%S]|END|%F]"})
	if err != nil {
		t.Fatalf("%s", err)
	}

	err = cmds.Validate("SCRIPT foo")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	err = cmds.Validate("SCRIPT END")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	err = cmds.Validate("SCRIPT RECORD foo")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	err = cmds.Validate("SCRIPT RECORD REGRESSION foo")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	err = cmds.Validate("SCRIPT RECORD REGRESSION foo end")
	if err == nil {
		t.Errorf("matches but shouldn't")
	} else {
		fmt.Println(err)
	}
}

func TestValidation_nestedGroups(t *testing.T) {
	var cmds *commandline.Commands
	var err error

	cmds, err = commandline.ParseCommandTemplate([]string{"TEST [(foo)|bar]"})
	if err != nil {
		t.Fatalf("%s", err)
	}
	err = cmds.Validate("TEST foo")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}
	err = cmds.Validate("TEST bar")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}
	err = cmds.Validate("TEST wibble")
	if err == nil {
		t.Errorf("matches but shouldn't")
	} else {
		fmt.Println(err)
	}

	cmds, err = commandline.ParseCommandTemplate([]string{"TEST (foo|[bar|(baz|qux)]|wibble)"})
	if err != nil {
		t.Fatalf("%s", err)
	}
	err = cmds.Validate("TEST foo")
	if err != nil {
		t.Errorf("1 doesn't match but should: %s", err)
	}
	err = cmds.Validate("TEST wibble")
	if err != nil {
		t.Errorf("2 doesn't match but should: %s", err)
	}
	err = cmds.Validate("TEST bar")
	if err != nil {
		t.Errorf("3 doesn't match but should: %s", err)
	}
}

func TestValidation_repeatGroups(t *testing.T) {
	var cmds *commandline.Commands
	var err error

	cmds, err = commandline.ParseCommandTemplate([]string{"TEST {foo}"})
	if err != nil {
		t.Fatalf("%s", err)
	}
	err = cmds.Validate("TEST foo")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}
	err = cmds.Validate("TEST foo foo")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	cmds, err = commandline.ParseCommandTemplate([]string{"TEST {foo|bar|baz}"})
	if err != nil {
		t.Fatalf("%s", err)
	}
	err = cmds.Validate("TEST foo")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}
	err = cmds.Validate("TEST foo foo")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	err = cmds.Validate("TEST bar foo")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	err = cmds.Validate("TEST bar foo baz baz")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	cmds, err = commandline.ParseCommandTemplate([]string{"TEST [foo|bar {baz|qux}]"})
	if err != nil {
		t.Fatalf("%s", err)
	}
	err = cmds.Validate("TEST foo")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}
	err = cmds.Validate("TEST bar")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}
	err = cmds.Validate("TEST bar baz")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}
	err = cmds.Validate("TEST bar baz qux")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}

	err = cmds.Validate("TEST foo bar")
	if err == nil {
		t.Errorf("matches but shouldn't")
	} else {
		fmt.Println(err)
	}

	err = cmds.Validate("TEST bar baz bar")
	if err == nil {
		t.Errorf("matches but shouldn't")
	} else {
		fmt.Println(err)
	}

	err = cmds.Validate("TEST bar baz qux qux baz wibble")
	if err == nil {
		t.Errorf("matches but shouldn't")
	} else {
		fmt.Println(err)
	}

	// the following template doesn't make sense (yet?) and should fail

	cmds, err = commandline.ParseCommandTemplate([]string{"TEST {[foo]}"})
	if err != nil {
		t.Fatalf("%s", err)
	}
	err = cmds.Validate("TEST foo")
	if err == nil {
		t.Errorf("matches but shouldn't")
	} else {
		fmt.Println(err)
	}
	err = cmds.Validate("TEST foo foo")
	if err == nil {
		t.Errorf("matches but shouldn't")
	} else {
		fmt.Println(err)
	}
}

func TestValidation_foo(t *testing.T) {
	var cmds *commandline.Commands
	var err error

	cmds, err = commandline.ParseCommandTemplate([]string{"SYMBOL [%S (ALL|MIRRORS)|LIST]"})
	if err != nil {
		t.Fatalf("%s", err)
	}
	err = cmds.Validate("SYMBOL enabl")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}
}

func TestValidation_bar(t *testing.T) {
	var cmds *commandline.Commands
	var err error

	cmds, err = commandline.ParseCommandTemplate([]string{
		"LIST",
		"PRINT [%s]",
		"SORT (RISING|FALLING)",
	})
	if err != nil {
		t.Fatalf("%s", err)
	}
	err = cmds.Validate("list")
	if err != nil {
		t.Errorf("doesn't match but should: %s", err)
	}
}
