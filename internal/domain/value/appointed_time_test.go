package value

import "testing"

func TestNewAppointedTime(t *testing.T) {
	// DOMAIN ERRORS
	// a.	value below minimum
	if _, err := NewAppointedTime(5); err == nil {
		t.Fatal("failed to create appointed time, expected value >= 10")
	}

	// b. value above maximum
	if _, err := NewAppointedTime(525601); err == nil {
		t.Fatal("failed to create appointed time, expected value <= 525600")
	}

	// VALID SCENARIOS
	// a.	value in the accepted minimum domain
	if _, err := NewAppointedTime(10); err != nil {
		t.Fatal("failed to create appointed time, expected nil error")
	}

	// b.	value in the accepted maximum domain
	if _, err := NewAppointedTime(525600); err != nil {
		t.Fatal("failed to create appointed time, expected nil error")
	}
}
