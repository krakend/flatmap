package flatmap

import (
	"regexp"
	"strconv"
	"strings"
)

var defaultCollectionPatter = regexp.MustCompile(`\.\*\.`)

func newMap(t Tokenizer) (*Map, error) {
	sep := t.Separator()
	var hasWildcard *regexp.Regexp
	var err error
	if sep == "." {
		hasWildcard = defaultCollectionPatter
	} else {
		hasWildcard, err = regexp.Compile(sep + `\*` + sep)
	}
	if err != nil {
		return nil, err
	}
	return &Map{
		m:  make(map[string]interface{}),
		t:  t,
		re: hasWildcard,
	}, nil
}

// Map is a flatten map
type Map struct {
	m  map[string]interface{}
	t  Tokenizer
	re *regexp.Regexp
}

// Move makes changes in the flatten hierarchy moving contents from origin to newKey
func (m *Map) Move(original, newKey string) {
	if v, ok := m.m[original]; ok {
		m.m[newKey] = v
		delete(m.m, original)
		return
	}

	if m.re.MatchString(original) {
		m.moveSliceAttribute(original, newKey)
		return
	}

	sep := m.t.Separator()

	for k := range m.m {
		if !strings.HasPrefix(k, original) {
			continue
		}

		if k[len(original):len(original)+1] != sep {
			continue
		}

		m.m[newKey+sep+k[len(original)+1:]] = m.m[k]
		delete(m.m, k)
	}
}

// Del deletes a key out of the map with the given prefix
func (m *Map) Del(prefix string) {
	if _, ok := m.m[prefix]; ok {
		delete(m.m, prefix)
		return
	}

	if m.re.MatchString(prefix) {
		m.delSliceAttribute(prefix)
		return
	}

	sep := m.t.Separator()

	for k := range m.m {
		if !strings.HasPrefix(k, prefix) {
			continue
		}

		if k[len(prefix):len(prefix)+1] != sep {
			continue
		}

		delete(m.m, k)
	}
}

func (m *Map) delSliceAttribute(prefix string) {
	i := strings.Index(prefix, "*")
	sep := m.t.Separator()
	prefixRemainder := prefix[i+1:]
	recursive := strings.Index(prefixRemainder, "*") > -1

	for k := range m.m {
		if len(k) < i+2 {
			continue
		}

		if !strings.HasPrefix(k, prefix[:i]) {
			continue
		}

		if recursive {
			// TODO: avoid recursive calls by managing nested collections in a single key evaluation
			newPref := k[:i+1+strings.Index(k[i+1:], sep)] + prefixRemainder
			m.Del(newPref)
			continue
		}

		keyRemainder := k[i+1+strings.Index(k[i+1:], sep):]
		if keyRemainder == prefixRemainder {
			delete(m.m, k)
			continue
		}

		if !strings.HasPrefix(keyRemainder, prefixRemainder+sep) {
			continue
		}

		delete(m.m, k)
	}
}

func (m *Map) moveSliceAttribute(original, newKey string) {
	i := strings.Index(original, "*")
	sep := m.t.Separator()
	originalRemainder := original[i+1:]
	recursive := strings.Index(originalRemainder, "*") > -1

	newKeyOffset := strings.Index(newKey, "*")
	newKeyRemainder := newKey[newKeyOffset+1:]
	newKeyPrefix := newKey[:newKeyOffset]

	for k := range m.m {
		if len(k) <= i+2 {
			continue
		}

		if !strings.HasPrefix(k, original[:i]) {
			continue
		}

		remainder := k[i:]
		idLen := strings.Index(remainder, sep)
		cleanRemainder := k[i+idLen:]
		keyPrefix := newKeyPrefix + k[i:i+idLen]

		if recursive {
			// TODO: avoid recursive calls by managing nested collections in a single key evaluation
			m.Move(k[:i+idLen]+originalRemainder, keyPrefix+newKeyRemainder)
			continue
		}

		if cleanRemainder == originalRemainder[1:] {
			m.m[keyPrefix+newKeyRemainder] = m.m[k]
			delete(m.m, k)
			continue
		}

		rPrefix := originalRemainder[1:] + sep

		if cleanRemainder != sep+originalRemainder[1:] && !strings.HasPrefix(cleanRemainder, sep+rPrefix) {
			continue
		}

		m.m[keyPrefix+newKeyRemainder+cleanRemainder[len(rPrefix):]] = m.m[k]
		delete(m.m, k)
	}
}

// Expand expands the Map into a more complex structure. This is the reverse of the Flatten operation.
func (m *Map) Expand() map[string]interface{} {
	res := map[string]interface{}{}
	hasCollections := false
	for k, v := range m.m {
		ks := m.t.Keys(k)
		tr := res

		if ks[len(ks)-1] == "#" {
			hasCollections = true
		}
		for _, tk := range ks[:len(ks)-1] {
			trnew, ok := tr[tk]
			if !ok {
				trnew = make(map[string]interface{})
				tr[tk] = trnew
			}
			tr = trnew.(map[string]interface{})
		}
		tr[ks[len(ks)-1]] = v
	}

	if !hasCollections {
		return res
	}

	return m.expandNestedCollections(res).(map[string]interface{})
}

func (m *Map) expandNestedCollections(original map[string]interface{}) interface{} {
	for k, v := range original {
		if t, ok := v.(map[string]interface{}); ok {
			original[k] = m.expandNestedCollections(t)
		}
	}

	size, ok := original["#"]
	if !ok {
		return original
	}

	col := make([]interface{}, size.(int))
	for k := range col {
		col[k] = original[strconv.Itoa(k)]
	}
	return col
}
