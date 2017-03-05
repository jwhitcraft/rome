package build

import "bytes"

var (
	andByte      = []byte("&&")
	orByte       = []byte("||")
	equalByte    = []byte("=")
	notEqualByte = []byte("!=")
	flavByte     = []byte("flav")
)

// processBuildTag takes a tag and evaluates it into it's proper boolean value
func processBuildTag(tag []byte, flavors []string) bool {
	// first things first, check for &&
	//tags := strings.Split(tag, "&&")
	tags := bytes.Split(tag, andByte)
	ok := true
	for _, tag := range tags {

		var tagOk bool
		if bytes.Contains(tag, orByte) {
			// split on the ||
			var orOk bool
			orTags := bytes.Split(tag, orByte)
			for _, orTag := range orTags {
				orOk = orOk || getTagBooleanValue(orTag, flavors)
			}
			tagOk = orOk
		} else {
			tagOk = getTagBooleanValue(tag, flavors)
		}

		ok = ok && tagOk
	}

	return ok
}

// getTagBooleanValue get the actual bool value of a section of the tag
func getTagBooleanValue(tag []byte, flavors []string) bool {
	tag = bytes.TrimSpace(tag)
	tagSep := getTagSeparator(tag)
	tagKey, tagVal := splitTag(tag, tagSep)

	var testValue []string
	// default the tag to be allowed, only change it something else is off
	switch string(tagKey) {
	case "flav":
		testValue = flavors
	case "lic":
		testValue = License[string(tagKey)]
	case "dep":
		testValue = []string{"os"}
	}
	if bytes.Equal(tagSep, notEqualByte) {
		return notContains(testValue, tagVal)
	}

	return contains(testValue, tagVal)
}

// getTagSeparator is used to figure out how to compare a tag to the kind of values
func getTagSeparator(tag []byte) []byte {
	if bytes.Contains(tag, notEqualByte) {
		return notEqualByte
	}

	return equalByte
}

// splitTag splits a tag into a key, value pair
func splitTag(eval, splitOn []byte) (key, val []byte) {
	splitVal := bytes.Split(eval, splitOn)

	// if the value only has item in it after split, assume it's the value not the key
	if len(splitVal) == 1 {
		key = flavByte
		val = splitVal[0]
	} else {
		key = splitVal[0]
		val = splitVal[1]
	}

	return bytes.TrimSpace(bytes.ToLower(key)), bytes.TrimSpace(bytes.ToLower(val))
}

// notContains is the opposite of contains.
func notContains(slice []string, item []byte) bool {
	return !contains(slice, item)
}

// contains is how we figure out if a tag is in a list of options
func contains(slice []string, item []byte) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[string(item)]
	return ok
}
