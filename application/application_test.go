package application

import (
	"testing"
)

var testConfig = &Config{
	Log: struct {
		FilePath             string
		MaxAgeHours          uint
		MaxRotationMegabytes uint
	}{FilePath: "logs/application.log", MaxAgeHours: 3 * 24, MaxRotationMegabytes: 1 * 1024},
}

func TestApplication(t *testing.T) {
	_, err := New(
		WithConfig(testConfig),
	)

	if err != nil {
		panic(err)
	}
}
