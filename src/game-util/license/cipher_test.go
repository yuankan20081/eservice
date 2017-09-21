package license

import "testing"

func TestCiper(t *testing.T) {
	s := "this is a test message"
	encrypted, err := RsaEncrypt([]byte(s))
	if err != nil {
		t.Error(err)
	}

	origin, err := RsaDecrypt(encrypted)
	if err != nil {
		t.Error(err)
	}

	if s != string(origin) {
		t.Error("failed get the original data")
	}
}
