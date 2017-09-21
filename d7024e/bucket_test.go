package d7024e_test

import (
	"Mobile_D7024E/d7024e"
	"testing"
)

const test_length int = 10

func TestAddNodeAndLen(t *testing.T) {
	bucket := d7024e.NewBucket()
	for i := 1; i <= test_length; i++ {
		bucket.AddNode(d7024e.NewNode(d7024e.NewRandomKademliaID(), "localhost:9000"))
		if bucket.Len() != i {
			t.Fail()
		}
	}
}
/*
func TestGetNodeAndCalcDistance(t *testing.T) {
	bucket := d7024e.NewBucket()
	for i := 1; i <= test_length; i++ {
		bucket.AddNode(d7024e.NewNode(d7024e.NewRandomKademliaID(), "localhost:9000"))
	}

	target := d7024e.NewRandomKademliaID()



}
*/
