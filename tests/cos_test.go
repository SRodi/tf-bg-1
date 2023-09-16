package test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/awserr"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

const (
	authEndpoint      = "https://iam.cloud.ibm.com/identity/token"
	serviceEndpoint   = "https://s3.eu-gb.cloud-object-storage.appdomain.cloud"
	testObjectKey     = "testKey1"
	testObjectContent = "some testing random text"
)

func createClient(apiKey, serviceInstanceID string) *s3.S3 {

	conf := aws.NewConfig().
		WithRegion("us-standard").
		WithEndpoint(serviceEndpoint).
		WithCredentials(ibmiam.NewStaticCredentials(aws.NewConfig(), authEndpoint, apiKey, serviceInstanceID)).
		WithS3ForcePathStyle(true)

	clientSession := session.Must(session.NewSession())

	client := s3.New(clientSession, conf)

	return client
}

func testListBuckets(t *testing.T, apiKey, serviceInstanceID, bucketName string) {

	client := createClient(apiKey, serviceInstanceID)

	listBucketOutput, err := client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Panic(err)
		t.Fail()
	}

	// test bucketName
	assert.Equal(t, bucketName, *listBucketOutput.Buckets[0].Name)

	// test we have only one bucket
	assert.Equal(t, 1, len(listBucketOutput.Buckets))

}

func testUploadObject(t *testing.T, apiKey, serviceInstanceID, bucketName string) {

	client := createClient(apiKey, serviceInstanceID)

	content := bytes.NewReader([]byte(testObjectContent))

	result, err := client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(testObjectKey),
		Body:   content,
	})
	if err != nil {
		log.Panic(err)
		t.Fail()
	}

	// test ETag exist, meaning the object is uploaded
	assert.NotEmpty(t, result.ETag)
}

func testGetObjectContent(t *testing.T, apiKey, serviceInstanceID, bucketName string) {

	client := createClient(apiKey, serviceInstanceID)

	res, err := client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(testObjectKey),
	})
	if err != nil {
		log.Panic(err)
		t.Fail()
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		log.Panic(err)
		t.Fail()
	}
	// test object content is what was uploaded
	assert.Equal(t, testObjectContent, string(content))
}

func testDeleteObject(t *testing.T, apiKey, serviceInstanceID, bucketName string) {

	client := createClient(apiKey, serviceInstanceID)

	_, err := client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(testObjectKey),
	})
	if err != nil {
		if awsError, ok := err.(awserr.Error); ok {
			// if there is an error we expect a specific failure
			// in this case AccessDenied as the service key
			// does not have permission to delete the object
			assert.Equal(t, "AccessDenied", awsError.Code())
		} else {
			log.Panic(err)
			t.Fail()
		}

	}
}

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

	keyObjectWriter := terraform.OutputRequired(t, opts, "key_object_writer")
	keyReader := terraform.OutputRequired(t, opts, "key_reader")
	serviceInstanceId := terraform.OutputRequired(t, opts, "service_instance_id")
	bucketName := terraform.OutputRequired(t, opts, "bucket_name")

	fmt.Println(serviceInstanceId, bucketName)

	// test the outputs for serviceInstanceId and bucketName
	// if this test is successful we are confident Terraform outputs are
	// correctly configured and they allow a client to connect to the COS instance
	// Note: this test uses the provisioning IBM Cloud apiKey that we use to
	// create the Terraform resources, this is stored as an env variable
	testListBuckets(t, os.Getenv("IC_API_KEY"), serviceInstanceId, bucketName)

	// test keyObjectWriter, if this test pass
	// we are confident Terraform output is valid and we have defined
	// the correct Write permission for this service key
	testUploadObject(t, keyObjectWriter, serviceInstanceId, bucketName)

	// test keyReader, if this test pass
	// we are confident Terraform output is valid and we have defined
	// the correct Read permission for this service key
	testGetObjectContent(t, keyReader, serviceInstanceId, bucketName)

	// test delete function with both keyObjectWriter and keyReader
	// expect failure as we did not set Delete permission to keys
	testDeleteObject(t, keyObjectWriter, serviceInstanceId, bucketName)
	testDeleteObject(t, keyReader, serviceInstanceId, bucketName)

	// test object can be deleted with provisioning apiKey
	// expect success
	testDeleteObject(t, os.Getenv("IC_API_KEY"), serviceInstanceId, bucketName)
}
