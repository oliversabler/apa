package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hej1 := &String{Value: "Hej Världen"}
	hej2 := &String{Value: "Hej Världen"}
	diff1 := &String{Value: "Mitt namn är johnny"}
	diff2 := &String{Value: "Mitt namn är johnny"}

	if hej1.HashKey() != hej2.HashKey() {
		t.Errorf("strings with the same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with the same content have different hash keys")
	}

	if hej1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have the same hash keys")
	}
}
