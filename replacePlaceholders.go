package slf

import "regexp"

// PlaceholdersRegex contains rules to find placeholders inside string
var PlaceholdersRegex = regexp.MustCompile(":[0-9a-zA-Z\\-_]+")

// ReplacePlaceholders func replaces placeholders inside source string using
// provided params
func ReplacePlaceholders(source string, params []Param, useBrackets bool) string {
	if len(params) == 0 {
		return source
	}

	// Building map
	mp := make(map[string]Param, len(params))
	for _, p := range params {
		mp[p.GetKey()] = p
	}

	return PlaceholdersRegex.ReplaceAllStringFunc(source, func(x string) string {
		key := x[1:]
		if v, ok := mp[key]; ok {
			value := v.String()
			if v.GetRaw() == nil {
				value = "nil"
			}
			if useBrackets {
				value = "[" + value + "]"
			}

			return value
		}
		return "<!" + x + ">"
	})
}
