package d7024e_test

import (
	"Mobile_D7024E/d7024e"
	"testing"
)

func TestKademliaInstantiation(t *testing.T) {
	kad1 := d7024e.NewKademlia()
	// kad2 := d7024e.NewKademlia()
	// kad3 := d7024e.NewKademlia()
	// kad4 := d7024e.NewKademlia()

	if kad1.LookupCount != 0 || kad1.RoutingTable != nil {
		t.Fail()
	}
}
