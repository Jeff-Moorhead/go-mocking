package friends

import "testing"

func TestEquals(t *testing.T) {
	this := Friend{
		"Jeff",
		25,
		"Programmer",
	}
	other := Friend{
		"Jeff",
		25,
		"Programmer",
	}

	if !this.Equals(other) {
		t.Errorf("Equal values are not considered equal by Equals(): %v, %v", this, other)
	}

}

func TestUnEquals(t *testing.T) {
	this := Friend{
		"Jeff",
		25,
		"Programmer",
	}
	other := Friend{
		"Emily",
		25,
		"Teacher",
	}
	if this.Equals(other) {
		t.Errorf("Unequal values are considered equal by Equals(): %v, %v", this, other)
	}
}

func TestDelete(t *testing.T) {
	friends := []Friend{
		{
			Name:       "Jeff",
			Age:        25,
			Occupation: "Programmer",
		},
		{
			Name:       "Emily",
			Age:        25,
			Occupation: "Teacher",
		},
		{
			Name:       "Paul",
			Age:        59,
			Occupation: "Cabinet Maker",
		},
	}
	friendstore := FriendStore{
		friends,
		"foo",
	}
	expected := []Friend{
		{
			Name:       "Jeff",
			Age:        25,
			Occupation: "Programmer",
		},
		{
			Name:       "Paul",
			Age:        59,
			Occupation: "Cabinet Maker",
		},
	}

	res := friendstore.Delete("Emily")
	if !res {
		t.Error("Return value should be true, got false")
	}

	if len(friendstore.friends) != len(expected) {
		t.Errorf("Incorrect number of elements, expected %v, got %v", expected, friendstore.friends)
		return
	}

	for i, v := range friendstore.friends {
		if !v.Equals(expected[i]) {
			t.Errorf("Incorrect value: %v, expected %v", v, expected[i])
		}
	}
}
