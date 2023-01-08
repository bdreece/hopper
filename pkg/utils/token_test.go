package utils_test

import (
	"testing"

	"github.com/bdreece/hopper/pkg/utils"
)

func TestToken1(t *testing.T) {
	secret := "secret123"
	token, _, err := utils.CreateToken("sub", "aud", secret)
	if err != nil {
		t.Errorf("An error occurred: %s\n", err.Error())
	}

	sub, aud, iss, err := utils.DecodeToken(*token, secret)
	if err != nil {
		t.Errorf("An error occurred: %s\n", err.Error())
	}

	if *sub != "sub" {
		t.Errorf("%s != sub\n", *sub)
	}

	if *aud != "aud" {
		t.Errorf("%s != aud\n", *aud)
	}

	if *iss != "hopper" {
		t.Errorf("%s != hopper\n", *iss)
	}
}
