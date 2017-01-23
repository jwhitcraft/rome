package build

import "testing"

func TestProcessBuildTag(t *testing.T) {

	var flav string = "flav=ent && flav!=dev"

	if !ProcessBuildTag(flav, Flavors["ent"]) {
		t.Errorf("Expected Value to be true, but got false for %s", flav)
	}

	flav = "flav=ent"
	if !ProcessBuildTag(flav, Flavors["ent"]) {
		t.Errorf("Expected Value to be true, but got false for %s", flav)
	}

	flav = "flav=pro && flav!=ent"
	if ProcessBuildTag(flav, Flavors["ent"]) {
		t.Errorf("Expected Value to be false, but got true for %s", flav)
	}
	if !ProcessBuildTag(flav, Flavors["pro"]) {
		t.Errorf("Expected Value to be true, but got false for %s", flav)
	}

	flav = "flav=corp"
	if ProcessBuildTag(flav, Flavors["ent"]) {
		t.Errorf("Expected Value to be false, but got true for %s", flav)
	}

}