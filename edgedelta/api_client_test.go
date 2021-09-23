package edgedelta

import (
	"flag"
	"os"
	"testing"
)

var (
	apiKey   = flag.String("api_key", "", "API auth key")
	confID   = flag.String("conf_id", "", "Unique configuration ID")
	orgID    = flag.String("org_id", "", "Unique organization ID")
	confPath = flag.String("conf_path", "", "Path to the new config file")
)

func TestGetConfigWithID(t *testing.T) {
	if *orgID == "" {
		t.Error("org_id is not specified")
		return
	}
	if *apiKey == "" {
		t.Error("api_key is not specified")
		return
	}
	if *confID == "" {
		t.Error("conf_id is not specified")
		return
	}

	cli := ConfigAPIClient{
		OrgID:      *orgID,
		apiKey:     *apiKey,
		APIBaseURL: "https://api.staging.edgedelta.com",
	}

	confObject, err := cli.getConfigWithID(*confID)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%v", confObject)
}

func TestUpdateConfigWithID(t *testing.T) {
	if *orgID == "" {
		t.Error("org_id is not specified")
		return
	}
	if *confID == "" {
		t.Error("conf_id is not specified")
		return
	}
	if *confPath == "" {
		t.Error("conf_path is not specified")
		return
	}

	cli := ConfigAPIClient{
		OrgID:      *orgID,
		APIBaseURL: "https://api.staging.edgedelta.com",
	}

	confDataRaw, err := os.ReadFile(*confPath)
	if err != nil {
		t.Error(err)
		return
	}

	confData := Config{
		Content: string(confDataRaw[:]),
	}

	confObject, err := cli.updateConfigWithID(*confID, confData)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%v", confObject)
}
