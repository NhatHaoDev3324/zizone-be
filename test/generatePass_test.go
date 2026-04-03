package test

import (
	"testing"

	"github.com/NhatHaoDev3324/zizone-be/utils"
)

func TestMain(t *testing.T) {
	email := "nguyennhathao.cm2k4@gmail.com"

	password, err := utils.SendPassword(email, "Nguyễn Văn A")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if password == "" {
		t.Errorf("Password should not be empty")
	}

	t.Logf("Generated password: %s", password)
}
