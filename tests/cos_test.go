package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

const (
	authEndpoint    = "https://iam.cloud.ibm.com/identity/token"
	serviceEndpoint = "https://s3.gb.cloud-object-storage.appdomain.cloud"
)

func TestInfraCOSExample(t *testing.T) {
	t.Parallel()

	opts := &terraform.Options{
		TerraformDir: "../examples/cos",

		Vars: map[string]interface{}{},
	}

	// clean up at the end of the test
	defer terraform.Destroy(t, opts)

	terraform.Init(t, opts)
	terraform.Apply(t, opts)

	// region := os.Getenv("IC_REGION")
	// apiKey := os.Getenv("IC_API_KEY")

	keyWriter := terraform.OutputRequired(t, opts, "key_writer")
	keyReader := terraform.OutputRequired(t, opts, "key_reader")
	serviceInstanceId := terraform.OutputRequired(t, opts, "service_instance_id")
}
