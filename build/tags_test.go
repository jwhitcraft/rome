package build

import (
	"bytes"
	"testing"
)

func TestPrivateProcessBuildTag(t *testing.T) {

	var flav []byte = []byte("flav=ent && flav!=dev")

	if !processBuildTag(flav, Flavors["ent"]) {
		t.Errorf("Expected Value to be true, but got false for %s", flav)
	}

	flav = []byte("flav=ent")
	if !processBuildTag(flav, Flavors["ent"]) {
		t.Errorf("Expected Value to be true, but got false for %s", flav)
	}

	flav = []byte("flav=pro && flav!=ent")
	if processBuildTag(flav, Flavors["ent"]) {
		t.Errorf("Expected Value to be false, but got true for %s", flav)
	}
	if !processBuildTag(flav, Flavors["pro"]) {
		t.Errorf("Expected Value to be true, but got false for %s", flav)
	}

	flav = []byte("flav=corp")
	if processBuildTag(flav, Flavors["ent"]) {
		t.Errorf("Expected Value to be false, but got true for %s", flav)
	}

	flav = []byte("flav = ent")
	if !processBuildTag(flav, Flavors["ent"]) {
		t.Errorf("Expected Value to be true, but got false for %s", flav)
	}

	flav = []byte("flav=pro || flav=com")
	if !processBuildTag(flav, Flavors["ent"]) {
		t.Errorf("Expected Value to be true, but got false for %s", flav)
	}
}

func TestPrivateSplitTag(t *testing.T) {
	var key []byte
	var val []byte

	key, val = splitTag([]byte("pro"), equalByte)

	if !bytes.Equal(key, []byte("flav")) || !bytes.Equal(val, []byte("pro")) {
		t.Errorf("Expected flav and pro but got %s and %s", key, val)
	}

	key, val = splitTag([]byte("lic=sub"), equalByte)

	if !bytes.Equal(key, []byte("lic")) || !bytes.Equal(val, []byte("sub")) {
		t.Errorf("Expected lice and sub but got %s and %s", key, val)
	}
}
