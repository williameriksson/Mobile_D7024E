package d7024e_test

import (
	"Mobile_D7024E/d7024e"
	"testing"
)

func TestNewKademliaID(t *testing.T) {
	str := "FFFFFFFF00000000000000000000000000000000"
	id := d7024e.NewKademliaID(str)
	to_str := id.String()

	if to_str != str {
		t.Fail()
	}
}

func TestLess(t *testing.T) {
	id1 := d7024e.NewKademliaID("FFFFFFFF000A00550000FFF0000AAAA000000000")
	id2 := d7024e.NewKademliaID("FFFFFFFF000A005500F000000000000000000000")

	if !id1.Less(id2) {
		t.Fail()
	}
	if id2.Less(id1) {
		t.Fail()
	}
}

func TestEquals(t *testing.T) {
	id1 := d7024e.NewKademliaID("FABC12754930476DEA648FE6A5C8E76A68690126")
	id2 := d7024e.NewKademliaID("FABC12754930476DEA648FE6A5C8E76A68690126")

	if !id1.Equals(id2) {
		t.Fail()
	}
}

func TestDistance(t *testing.T) {
	id1 := d7024e.NewKademliaID("FABC12754930476DEA648FE6A5C8E76A68690126")
	id2 := d7024e.NewKademliaID("FABC12754930476DEA648FE6A5C8E76A68690126")
	id3 := d7024e.NewKademliaID("F74529AB9839C734DF98A9B7590A967F0C0BA082")

	dist1 := d7024e.NewKademliaID("0000000000000000000000000000000000000000")
	dist2 := d7024e.NewKademliaID("0DF93BDED109805935FC2651FCC271156462A1A4")

	if id1.CalcDistance(id2).String() != dist1.String() {
		t.Fail()
	}

	if id1.CalcDistance(id3).String() != dist2.String() {
		t.Fail()
	}
}

func TestNewRandomKademliaID(t *testing.T) {
	id1 := d7024e.NewRandomKademliaID()
	id2 := d7024e.NewRandomKademliaID()

	/* Extremely low chance of failing by chance */ 
	if id1.Equals(id2) {
		t.Fail()
	}

	if len(id1) != len(id2) && len(id1) != d7024e.IDLength {
		t.Fail()
	}
}