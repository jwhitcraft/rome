package build

import "strings"

// processBuildTag takes a tag and evaluates it into it's proper boolean value
func processBuildTag(tag string, flavors []string) bool {
	// first things first, check for &&
	tags := strings.Split(tag, "&&")
	ok := true
	for _, tag := range tags {

		var tagOk bool
		if strings.Contains(tag, "||") {
			// split on the ||
			var orOk bool
			orTags := strings.Split(tag, "||")
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
func getTagBooleanValue(tag string, flavors []string) bool {
	tag = strings.TrimSpace(tag)
	tagSep := getTagSeparator(tag)
	tagKey, tagVal := splitTag(tag, tagSep)

	var testValue []string
	// default the tag to be allowed, only change it something else is off
	switch tagKey {
	case "flav":
		testValue = flavors
	case "lic":
		testValue = License[tagKey]
	case "dep":
		testValue = []string{"os"}
	}
	if tagSep == "!=" {
		return notContains(testValue, tagVal)
	}

	return contains(testValue, tagVal)
}

// getTagSeparator is used to figure out how to compare a tag to the kind of values
func getTagSeparator(tag string) string {
	if strings.Contains(tag, "!=") {
		return "!="
	}

	return "="
}

// splitTag splits a tag into a key, value pair
func splitTag(eval, splitOn string) (key, val string) {
	splitVal := strings.Split(eval, splitOn)

	// if the value only has item in it after split, assume it's the value not the key
	if len(splitVal) == 1 {
		key = "flav"
		val = splitVal[0]
	} else {
		key = splitVal[0]
		val = splitVal[1]
	}

	return strings.TrimSpace(strings.ToLower(key)), strings.TrimSpace(strings.ToLower(val))
}

// notContains is the opposite of contains.
func notContains(slice []string, item string) bool {
	return !contains(slice, item)
}

// contains is how we figure out if a tag is in a list of options
func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
