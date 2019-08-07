package common

import (
	"testing"
)
type Profile struct {
	Name string
	Age  int
}
func TestGetBytes(t *testing.T)  {
	// --- encode ---
	profile := &Profile{
		Name: "Roman",
		Age:  30,
	}

	bts, err := GetBytes(profile)
	if err != nil {
		t.Fatal("encode error:",err.Error())
		return
	}
	// --- decode ---
	var decodedProfile Profile
	err = GetInterface(bts, &decodedProfile)
	if err != nil {
		t.Fatal("encode error:",err.Error())
		return
	}
	if decodedProfile.Name != "Roman" {
		t.Fatal("encode error:",err.Error())
		return
	}
	if decodedProfile.Age != 30 {
		t.Fatal("encode error:",err.Error())
		return
	}

}
