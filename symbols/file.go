package symbols

import (
	"fmt"
	"gopher2600/errors"
	"gopher2600/hardware/memory/addresses"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// ReadSymbolsFile initialises a symbols table from the symbols file for the
// specified cartridge
//
// Table instance will always be valid even if error is returned. for example,
// if the symbols file cannot be opened the symbols file will still contain the
// canonical vcs symbols file
//
// currently, only symbols files generated by DASM are supported
func ReadSymbolsFile(cartridgeFilename string) (*Table, error) {
	table := new(Table)
	table.Locations = newTable()
	table.Read = newTable()
	table.Write = newTable()

	// prioritise symbols with reference symbols for the VCS.
	//
	// deferred function because we want to do this in all instances, even if
	// there is an error with the symbols file.
	defer func() {
		for k, v := range addresses.Read {
			table.Read.add(k, v, true)
		}
		for k, v := range addresses.Write {
			table.Write.add(k, v, true)
		}

		sort.Sort(table.Locations)
		sort.Sort(table.Read)
		sort.Sort(table.Write)

		// find max symbol width
		table.MaxLocationWidth = table.Locations.maxWidth
		if table.Read.maxWidth > table.Write.maxWidth {
			table.MaxSymbolWidth = table.Read.maxWidth
		} else {
			table.MaxSymbolWidth = table.Write.maxWidth
		}
	}()

	// if this is the empty cartridge then this error is expected. return
	// the empty symbol table
	if cartridgeFilename == "" {
		return table, nil
	}

	// try to open symbols file
	symFilename := cartridgeFilename
	ext := path.Ext(symFilename)

	// try to figure out the case of the file extension
	if ext == ".BIN" {
		symFilename = fmt.Sprintf("%s.SYM", symFilename[:len(symFilename)-len(ext)])
	} else {
		symFilename = fmt.Sprintf("%s.sym", symFilename[:len(symFilename)-len(ext)])
	}

	sf, err := os.Open(symFilename)
	if err != nil {
		return table, errors.New(errors.SymbolsFileUnavailable, cartridgeFilename)
	}
	defer func() {
		_ = sf.Close()
	}()

	sym, err := ioutil.ReadAll(sf)
	if err != nil {
		return nil, errors.New(errors.SymbolsFileError, err)
	}
	lines := strings.Split(string(sym), "\n")

	// find interesting lines in the symbols file and add to the Table
	// instance.
	for _, ln := range lines {
		// ignore uninteresting lines
		p := strings.Fields(ln)
		if len(p) < 2 || p[0] == "---" {
			continue // for loop
		}

		// get address
		address, err := strconv.ParseUint(p[1], 16, 16)
		if err != nil {
			continue // for loop
		}

		// get symbol
		symbol := p[0]

		// differentiate between location and other symbols. this is a little
		// heavy handed, but still, it's better than nothing.
		if unicode.IsDigit(rune(symbol[0])) {
			// if symbol begins with a number and a period then it is a location symbol
			i := strings.Index(symbol, ".")
			if i != -1 {
				table.Locations.add(uint16(address), symbol[i:], false)
			}
		} else {
			// every non-location symbols is both a read and write symbol.
			// compar to canonical vcs symbols which are specific to a read or
			// write context
			table.Read.add(uint16(address), symbol, false)
			table.Write.add(uint16(address), symbol, false)
		}
	}

	return table, nil
}
