package database

import "gopher2600/errors"

// SelectAll entries in the database. onSelect can be nil.
//
// returns last matched entry in selection or an error with the last entry
// matched before the error occurred.
func (db Session) SelectAll(onSelect func(Entry) (bool, error)) (Entry, error) {
	var entry Entry

	if onSelect == nil {
		onSelect = func(_ Entry) (bool, error) { return true, nil }
	}

	keyList := db.SortedKeyList()

	for k := range keyList {
		entry := db.entries[keyList[k]]
		cont, err := onSelect(entry)
		if err != nil {
			return entry, err
		}
		if !cont {
			break // for loop
		}
	}

	return entry, nil
}

// SelectKeys matches entries with the specified key(s). keys can be singular.
// if list of keys is empty then all keys are matched (SelectAll() maybe more
// appropriate in that case). onSelect can be nil.
//
// returns last matched entry in selection or an error with the last entry
// matched before the error occurred.
func (db Session) SelectKeys(onSelect func(Entry) (bool, error), keys ...int) (Entry, error) {
	var entry Entry

	if onSelect == nil {
		onSelect = func(_ Entry) (bool, error) { return true, nil }
	}

	keyList := keys
	if len(keys) == 0 {
		keyList = db.SortedKeyList()
	}

	for i := range keyList {
		entry = db.entries[keyList[i]]
		cont, err := onSelect(entry)
		if err != nil {
			return entry, err
		}
		if !cont {
			break // for loop
		}
	}

	if entry == nil {
		return nil, errors.New(errors.DatabaseSelectEmpty)
	}

	return entry, nil
}
