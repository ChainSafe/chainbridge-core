package chain

import (
	"testing"
)

func TestValidateConfig(t *testing.T) {
	var id uint8 = 1
	valid := GeneralChainConfig{
		Name:     "chain",
		Id:       &id,
		Endpoint: "endpoint",
		From:     "0x0",
	}

	missingEndpoint := GeneralChainConfig{
		Name:     "chain",
		Id:       &id,
		Endpoint: "",
		From:     "0x0",
	}

	missingName := GeneralChainConfig{
		Name:     "",
		Id:       &id,
		Endpoint: "endpoint",
		From:     "0x0",
	}

	missingId := GeneralChainConfig{
		Name:     "chain",
		Endpoint: "endpoint",
		From:     "0x0",
	}

	missingFrom := GeneralChainConfig{
		Name:     "chain",
		Id:       &id,
		Endpoint: "endpoint",
		From:     "",
	}

	err := valid.Validate()
	if err != nil {
		t.Fatal(err)
	}

	err = missingEndpoint.Validate()
	if err == nil {
		t.Fatalf("must require endpoint field, %v", err)
	}

	err = missingName.Validate()
	if err == nil {
		t.Fatal("must require name field")
	}

	err = missingId.Validate()
	if err == nil {
		t.Fatalf("must require domain id field, %v", err)
	}

	err = missingFrom.Validate()
	if err == nil {
		t.Fatalf("must require from field, %v", err)
	}
}
