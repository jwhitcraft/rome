package build

import "testing"

func TestPrivateProcessBuildTag(t *testing.T) {

	var flav string = "flav=ent && flav!=dev"

	if !processBuildTag(flav, Flavors["ent"]) {
		t.Errorf("Expected Value to be true, but got false for %s", flav)
	}

	flav = "flav=ent"
	if !processBuildTag(flav, Flavors["ent"]) {
		t.Errorf("Expected Value to be true, but got false for %s", flav)
	}

	flav = "flav=pro && flav!=ent"
	if processBuildTag(flav, Flavors["ent"]) {
		t.Errorf("Expected Value to be false, but got true for %s", flav)
	}
	if !processBuildTag(flav, Flavors["pro"]) {
		t.Errorf("Expected Value to be true, but got false for %s", flav)
	}

	flav = "flav=corp"
	if processBuildTag(flav, Flavors["ent"]) {
		t.Errorf("Expected Value to be false, but got true for %s", flav)
	}

	flav = "flav = ent"
	if !processBuildTag(flav, Flavors["ent"]) {
		t.Errorf("Expected Value to be true, but got false for %s", flav)
	}

	flav = "flav=pro || flav=com"
	if !processBuildTag(flav, Flavors["ent"]) {
		t.Errorf("Expected Value to be true, but got false for %s", flav)
	}
}

func TestPrivateSplitTag(t *testing.T) {
	var key string
	var val string

	key, val = splitTag("pro", "=")

	if key != "flav" || val != "pro" {
		t.Errorf("Expected flav and pro but got %s and %s", key, val)
	}

	key, val = splitTag("lic=sub", "=")

	if key != "lic" || val != "sub" {
		t.Errorf("Expected lice and sub but got %s and %s", key, val)
	}
}
