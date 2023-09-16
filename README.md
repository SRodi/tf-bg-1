## Build on IBM Cloud with production grade Terraform
As part of a production grade Terraform deployment, testing plays an integral part. This tutorial will demonstrate how to write and effectively test a Terraform module before deploying to a production environment.

## Why testing Infrastructure as Code
Testing Infrastructure as Code (IaC) is as essential as testing any other software code. As part of a Terraform module, multiple resource outputs are normally used and consumed by other resources part of the Terraform stack for a cloud environment. If the outputs are not correct, we could break entire production systems. This not only occurs at resource creation, but also when updating the Terraform configuration for the same resources.

For a production grade Terraform deployment, IaC testing should be performed in automation with dedicated CI pipelines.

## Tutorial Goals
In this tutorial we will perform the following:

1. Create an IBM Cloud COS instance and bucket
1. Create a resource key with Object Writer permissions to the bucket
1. Create a resource key with Reader permissions to the bucket
1. Validate Terraform resources are correctly configured
1. Test COS bucket creation
1. Test writer key can upload an object to the bucket
1. Test reader key can read an object stored in the bucket
1. Test writer and reader key cannot delete an object stored in the bucket

## High level plan
We will create the infrastructure with Terraform in a modularised approach, for the purpose of the tutorial we create one single module. Then we will test the module by creating an example which instantiate the module implementation. Finally we will perform tests using Terratest, which is a Go testing library for terraform, and also IBM Cloud COS SDK for Go, to test with a real client and real resources.

## Prerequisites
The following are requirements for this tutorial, please make sure you have an IBM Cloud account with and apiKey available which has full permissions for COS service. In addition you need to have Terraform, Go and Git installed on your local machine.
1. IBM Cloud account
1. Terraform CLI
1. Go
1. Git

The versions I have currently installed on my environment are as follows:
```bash
❯ terraform version
Terraform v1.3.9
on darwin_amd64

❯ go version
go version go1.20.6 darwin/amd64

❯ git version
git version 2.39.2 (Apple Git-143)
```

## References
* [IBM Cloud docs: cloud-object-storage-using-go](https://cloud.ibm.com/docs/cloud-object-storage?topic=cloud-object-storage-using-go)
* [Tutorial GitHub repo](https://github.com/SRodi/tf-bg-1)
* [Terraform providers registry: IBM-Cloud](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest)

### Quick-start
Export the required environment variables:
```bash
export IC_API_KEY=yourSuperSecretApiKey
export IC_REGION=eu-gb
```

Run the test:
```bash
cd tests
go test -run TestInfraCOSExample
```

Example output:
```terraform
❯ go test -run TestInfraCOSExample
TestInfraCOSExample 2023-09-16T10:17:57+01:00 retry.go:91: terraform [init -upgrade=false]
TestInfraCOSExample 2023-09-16T10:17:57+01:00 logger.go:66: Running command terraform with args [init -upgrade=false]
TestInfraCOSExample 2023-09-16T10:17:57+01:00 logger.go:66: Initializing modules...
TestInfraCOSExample 2023-09-16T10:17:57+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:17:57+01:00 logger.go:66: Initializing the backend...
TestInfraCOSExample 2023-09-16T10:17:58+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:17:58+01:00 logger.go:66: Initializing provider plugins...
TestInfraCOSExample 2023-09-16T10:17:58+01:00 logger.go:66: - Reusing previous version of ibm-cloud/ibm from the dependency lock file
TestInfraCOSExample 2023-09-16T10:17:59+01:00 logger.go:66: - Using previously-installed ibm-cloud/ibm v1.57.0
TestInfraCOSExample 2023-09-16T10:17:59+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:17:59+01:00 logger.go:66: Terraform has been successfully initialized!
TestInfraCOSExample 2023-09-16T10:17:59+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:17:59+01:00 logger.go:66: You may now begin working with Terraform. Try running "terraform plan" to see
TestInfraCOSExample 2023-09-16T10:17:59+01:00 logger.go:66: any changes that are required for your infrastructure. All Terraform commands
TestInfraCOSExample 2023-09-16T10:17:59+01:00 logger.go:66: should now work.
TestInfraCOSExample 2023-09-16T10:17:59+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:17:59+01:00 logger.go:66: If you ever set or change modules or backend configuration for Terraform,
TestInfraCOSExample 2023-09-16T10:17:59+01:00 logger.go:66: rerun this command to reinitialize your working directory. If you forget, other
TestInfraCOSExample 2023-09-16T10:17:59+01:00 logger.go:66: commands will detect it and remind you to do so if necessary.
TestInfraCOSExample 2023-09-16T10:17:59+01:00 retry.go:91: terraform [apply -input=false -auto-approve -lock=false]
TestInfraCOSExample 2023-09-16T10:17:59+01:00 logger.go:66: Running command terraform with args [apply -input=false -auto-approve -lock=false]
TestInfraCOSExample 2023-09-16T10:18:02+01:00 logger.go:66: module.cos.data.ibm_resource_group.resource_group_default: Reading...
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66: module.cos.data.ibm_resource_group.resource_group_default: Read complete after 2s [id=4bef339d68b14906b5008348875d13db]
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66: Terraform used the selected providers to generate the following execution
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66: plan. Resource actions are indicated with the following symbols:
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:   + create
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66: Terraform will perform the following actions:
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:   # module.cos.ibm_cos_bucket.bucket_test will be created
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:   + resource "ibm_cos_bucket" "bucket_test" {
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + bucket_name          = "cos-bucket-tf-test"
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + crn                  = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + endpoint_type        = "public"
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + force_delete         = true
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + id                   = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + region_location      = "eu-gb"
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + resource_instance_id = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + s3_endpoint_direct   = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + s3_endpoint_private  = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + s3_endpoint_public   = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + storage_class        = "standard"
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:     }
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:   # module.cos.ibm_resource_instance.resource_instance_cos_test will be created
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:   + resource "ibm_resource_instance" "resource_instance_cos_test" {
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + account_id              = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + allow_cleanup           = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + created_at              = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + created_by              = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + crn                     = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + dashboard_url           = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + deleted_at              = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + deleted_by              = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + extensions              = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + guid                    = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + id                      = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + last_operation          = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + location                = "global"
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + locked                  = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + name                    = "cos-instance-test"
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + plan                    = "lite"
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + plan_history            = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + resource_aliases_url    = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + resource_bindings_url   = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + resource_controller_url = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + resource_crn            = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + resource_group_crn      = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + resource_group_id       = "4bef339d68b14906b5008348875d13db"
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + resource_group_name     = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + resource_id             = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + resource_keys_url       = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + resource_name           = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + resource_plan_id        = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + resource_status         = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + restored_at             = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + restored_by             = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + scheduled_reclaim_at    = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + scheduled_reclaim_by    = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + service                 = "cloud-object-storage"
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + service_endpoints       = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + state                   = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + status                  = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + sub_type                = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + tags                    = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + target_crn              = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + type                    = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + update_at               = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + update_by               = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:     }
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:   # module.cos.ibm_resource_key.key_cos_object_writer will be created
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:   + resource "ibm_resource_key" "key_cos_object_writer" {
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + account_id            = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + created_at            = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + created_by            = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + credentials           = (sensitive value)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + credentials_json      = (sensitive value)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + crn                   = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + deleted_at            = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + deleted_by            = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + guid                  = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + iam_compatible        = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + id                    = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + name                  = "cos-key-object-writer-test"
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + resource_group_id     = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + resource_instance_id  = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + resource_instance_url = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + role                  = "Object Writer"
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + source_crn            = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + state                 = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + status                = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + updated_at            = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + updated_by            = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + url                   = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:     }
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:   # module.cos.ibm_resource_key.key_cos_reader will be created
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:   + resource "ibm_resource_key" "key_cos_reader" {
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + account_id            = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + created_at            = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + created_by            = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + credentials           = (sensitive value)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + credentials_json      = (sensitive value)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + crn                   = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + deleted_at            = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + deleted_by            = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + guid                  = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + iam_compatible        = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + id                    = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + name                  = "cos-key-reader-test"
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + resource_group_id     = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + resource_instance_id  = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + resource_instance_url = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + role                  = "Reader"
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + source_crn            = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + state                 = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + status                = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + updated_at            = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + updated_by            = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:       + url                   = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:     }
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66: Plan: 4 to add, 0 to change, 0 to destroy.
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66: Changes to Outputs:
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:   + bucket_name         = "cos-bucket-tf-test"
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:   + key_object_writer   = (sensitive value)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:   + key_reader          = (sensitive value)
TestInfraCOSExample 2023-09-16T10:18:05+01:00 logger.go:66:   + service_instance_id = (known after apply)
TestInfraCOSExample 2023-09-16T10:18:07+01:00 logger.go:66: module.cos.ibm_resource_instance.resource_instance_cos_test: Creating...
TestInfraCOSExample 2023-09-16T10:18:17+01:00 logger.go:66: module.cos.ibm_resource_instance.resource_instance_cos_test: Still creating... [10s elapsed]
TestInfraCOSExample 2023-09-16T10:18:25+01:00 logger.go:66: module.cos.ibm_resource_instance.resource_instance_cos_test: Creation complete after 18s [id=crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc::]
TestInfraCOSExample 2023-09-16T10:18:25+01:00 logger.go:66: module.cos.ibm_resource_key.key_cos_object_writer: Creating...
TestInfraCOSExample 2023-09-16T10:18:25+01:00 logger.go:66: module.cos.ibm_resource_key.key_cos_reader: Creating...
TestInfraCOSExample 2023-09-16T10:18:25+01:00 logger.go:66: module.cos.ibm_cos_bucket.bucket_test: Creating...
TestInfraCOSExample 2023-09-16T10:18:29+01:00 logger.go:66: module.cos.ibm_resource_key.key_cos_object_writer: Creation complete after 4s [id=crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc:resource-key:2cc846ff-c393-4e6c-8e06-d18845cbb79b]
TestInfraCOSExample 2023-09-16T10:18:29+01:00 logger.go:66: module.cos.ibm_resource_key.key_cos_reader: Creation complete after 4s [id=crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc:resource-key:9e7aec27-3d9c-45a8-b031-ce4e0aa985e9]
TestInfraCOSExample 2023-09-16T10:18:31+01:00 logger.go:66: module.cos.ibm_cos_bucket.bucket_test: Creation complete after 7s [id=crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc:bucket:cos-bucket-tf-test:meta:rl:eu-gb:public]
TestInfraCOSExample 2023-09-16T10:18:31+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:18:31+01:00 logger.go:66: Apply complete! Resources: 4 added, 0 changed, 0 destroyed.
TestInfraCOSExample 2023-09-16T10:18:31+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:18:31+01:00 logger.go:66: Outputs:
TestInfraCOSExample 2023-09-16T10:18:31+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:18:31+01:00 logger.go:66: bucket_name = "cos-bucket-tf-test"
TestInfraCOSExample 2023-09-16T10:18:31+01:00 logger.go:66: key_object_writer = <sensitive>
TestInfraCOSExample 2023-09-16T10:18:31+01:00 logger.go:66: key_reader = <sensitive>
TestInfraCOSExample 2023-09-16T10:18:31+01:00 logger.go:66: service_instance_id = "7756104a-c094-4834-8e13-ee07be2518bc"
TestInfraCOSExample 2023-09-16T10:18:31+01:00 retry.go:91: terraform [output -no-color -json key_object_writer]
TestInfraCOSExample 2023-09-16T10:18:31+01:00 logger.go:66: Running command terraform with args [output -no-color -json key_object_writer]
TestInfraCOSExample 2023-09-16T10:18:32+01:00 logger.go:66: "tExcGJ04mBoMRTYht1U9QqF3vXxSgFXkzQxl8_XuafKa"
TestInfraCOSExample 2023-09-16T10:18:32+01:00 retry.go:91: terraform [output -no-color -json key_reader]
TestInfraCOSExample 2023-09-16T10:18:32+01:00 logger.go:66: Running command terraform with args [output -no-color -json key_reader]
TestInfraCOSExample 2023-09-16T10:18:32+01:00 logger.go:66: "3gHC7NiPmtORqLWQneh0XGeF69Dp9LzS82ogrFKHSbF4"
TestInfraCOSExample 2023-09-16T10:18:32+01:00 retry.go:91: terraform [output -no-color -json service_instance_id]
TestInfraCOSExample 2023-09-16T10:18:32+01:00 logger.go:66: Running command terraform with args [output -no-color -json service_instance_id]
TestInfraCOSExample 2023-09-16T10:18:32+01:00 logger.go:66: "7756104a-c094-4834-8e13-ee07be2518bc"
TestInfraCOSExample 2023-09-16T10:18:32+01:00 retry.go:91: terraform [output -no-color -json bucket_name]
TestInfraCOSExample 2023-09-16T10:18:32+01:00 logger.go:66: Running command terraform with args [output -no-color -json bucket_name]
TestInfraCOSExample 2023-09-16T10:18:33+01:00 logger.go:66: "cos-bucket-tf-test"
TestInfraCOSExample 2023-09-16T10:18:36+01:00 retry.go:91: terraform [destroy -auto-approve -input=false -lock=false]
TestInfraCOSExample 2023-09-16T10:18:36+01:00 logger.go:66: Running command terraform with args [destroy -auto-approve -input=false -lock=false]
TestInfraCOSExample 2023-09-16T10:18:38+01:00 logger.go:66: module.cos.data.ibm_resource_group.resource_group_default: Reading...
TestInfraCOSExample 2023-09-16T10:18:39+01:00 logger.go:66: module.cos.data.ibm_resource_group.resource_group_default: Read complete after 1s [id=4bef339d68b14906b5008348875d13db]
TestInfraCOSExample 2023-09-16T10:18:39+01:00 logger.go:66: module.cos.ibm_resource_instance.resource_instance_cos_test: Refreshing state... [id=crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc::]
TestInfraCOSExample 2023-09-16T10:18:41+01:00 logger.go:66: module.cos.ibm_resource_key.key_cos_reader: Refreshing state... [id=crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc:resource-key:9e7aec27-3d9c-45a8-b031-ce4e0aa985e9]
TestInfraCOSExample 2023-09-16T10:18:41+01:00 logger.go:66: module.cos.ibm_resource_key.key_cos_object_writer: Refreshing state... [id=crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc:resource-key:2cc846ff-c393-4e6c-8e06-d18845cbb79b]
TestInfraCOSExample 2023-09-16T10:18:41+01:00 logger.go:66: module.cos.ibm_cos_bucket.bucket_test: Refreshing state... [id=crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc:bucket:cos-bucket-tf-test:meta:rl:eu-gb:public]
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66: Terraform used the selected providers to generate the following execution
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66: plan. Resource actions are indicated with the following symbols:
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:   - destroy
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66: Terraform will perform the following actions:
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:   # module.cos.ibm_cos_bucket.bucket_test will be destroyed
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:   - resource "ibm_cos_bucket" "bucket_test" {
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - allowed_ip           = [] -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - bucket_name          = "cos-bucket-tf-test" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - crn                  = "crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc:bucket:cos-bucket-tf-test" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - endpoint_type        = "public" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - force_delete         = true -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - hard_quota           = 0 -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - id                   = "crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc:bucket:cos-bucket-tf-test:meta:rl:eu-gb:public" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - region_location      = "eu-gb" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - resource_instance_id = "crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc::" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - s3_endpoint_direct   = "s3.direct.eu-gb.cloud-object-storage.appdomain.cloud" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - s3_endpoint_private  = "s3.private.eu-gb.cloud-object-storage.appdomain.cloud" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - s3_endpoint_public   = "s3.eu-gb.cloud-object-storage.appdomain.cloud" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - storage_class        = "standard" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:     }
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:   # module.cos.ibm_resource_instance.resource_instance_cos_test will be destroyed
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:   - resource "ibm_resource_instance" "resource_instance_cos_test" {
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - account_id              = "7f57936089fa4570adc362a257bb04a0" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - allow_cleanup           = false -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - created_at              = "2023-09-16T09:18:09.775Z" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - created_by              = "iam-ServiceId-46fc18c6-40fd-4c4c-bde3-736b3d6597dc" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - crn                     = "crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc::" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - dashboard_url           = "https://cloud.ibm.com/objectstorage/crn%3Av1%3Abluemix%3Apublic%3Acloud-object-storage%3Aglobal%3Aa%2F7f57936089fa4570adc362a257bb04a0%3A7756104a-c094-4834-8e13-ee07be2518bc%3A%3A" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - extensions              = {} -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - guid                    = "7756104a-c094-4834-8e13-ee07be2518bc" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - id                      = "crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc::" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - last_operation          = {
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:           - "async"       = "false"
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:           - "cancelable"  = "false"
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:           - "description" = "Completed create instance operation"
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:           - "poll"        = "false"
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:           - "state"       = "succeeded"
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:           - "type"        = "create"
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:         } -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - location                = "global" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - locked                  = false -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - name                    = "cos-instance-test" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - plan                    = "lite" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - plan_history            = [
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:           - {
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:               - resource_plan_id = "2fdf0c08-2d32-4f46-84b5-32e0c92fffd8"
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:               - start_date       = "2023-09-16T09:18:09.775Z"
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:             },
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:         ] -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - resource_aliases_url    = "/v2/resource_instances/7756104a-c094-4834-8e13-ee07be2518bc/resource_aliases" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - resource_bindings_url   = "/v2/resource_instances/7756104a-c094-4834-8e13-ee07be2518bc/resource_bindings" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - resource_controller_url = "https://cloud.ibm.com/services/" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - resource_crn            = "crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc::" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - resource_group_crn      = "crn:v1:bluemix:public:resource-controller::a/7f57936089fa4570adc362a257bb04a0::resource-group:4bef339d68b14906b5008348875d13db" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - resource_group_id       = "4bef339d68b14906b5008348875d13db" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - resource_group_name     = "crn:v1:bluemix:public:resource-controller::a/7f57936089fa4570adc362a257bb04a0::resource-group:4bef339d68b14906b5008348875d13db" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - resource_id             = "dff97f5c-bc5e-4455-b470-411c3edbe49c" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - resource_keys_url       = "/v2/resource_instances/7756104a-c094-4834-8e13-ee07be2518bc/resource_keys" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - resource_name           = "cos-instance-test" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - resource_plan_id        = "2fdf0c08-2d32-4f46-84b5-32e0c92fffd8" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - resource_status         = "active" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - service                 = "cloud-object-storage" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - state                   = "active" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - status                  = "active" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - tags                    = [] -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - target_crn              = "crn:v1:bluemix:public:globalcatalog::::deployment:2fdf0c08-2d32-4f46-84b5-32e0c92fffd8%3Aglobal" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - type                    = "service_instance" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - update_at               = "2023-09-16T09:18:13.166Z" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:     }
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:   # module.cos.ibm_resource_key.key_cos_object_writer will be destroyed
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:   - resource "ibm_resource_key" "key_cos_object_writer" {
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - account_id            = "7f57936089fa4570adc362a257bb04a0" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - created_at            = "2023-09-16T09:18:27.825Z" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - created_by            = "iam-ServiceId-46fc18c6-40fd-4c4c-bde3-736b3d6597dc" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - credentials           = (sensitive value)
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - credentials_json      = (sensitive value)
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - crn                   = "crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc:resource-key:2cc846ff-c393-4e6c-8e06-d18845cbb79b" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - guid                  = "2cc846ff-c393-4e6c-8e06-d18845cbb79b" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - iam_compatible        = true -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - id                    = "crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc:resource-key:2cc846ff-c393-4e6c-8e06-d18845cbb79b" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - name                  = "cos-key-object-writer-test" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - resource_group_id     = "4bef339d68b14906b5008348875d13db" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - resource_instance_id  = "crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc::" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - resource_instance_url = "/v2/resource_instances/7756104a-c094-4834-8e13-ee07be2518bc" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - role                  = "Object Writer" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - source_crn            = "crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc::" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - state                 = "active" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - status                = "active" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - updated_at            = "2023-09-16T09:18:27.825Z" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - url                   = "/v2/resource_keys/2cc846ff-c393-4e6c-8e06-d18845cbb79b" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:     }
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:   # module.cos.ibm_resource_key.key_cos_reader will be destroyed
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:   - resource "ibm_resource_key" "key_cos_reader" {
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - account_id            = "7f57936089fa4570adc362a257bb04a0" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - created_at            = "2023-09-16T09:18:27.938Z" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - created_by            = "iam-ServiceId-46fc18c6-40fd-4c4c-bde3-736b3d6597dc" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - credentials           = (sensitive value)
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - credentials_json      = (sensitive value)
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - crn                   = "crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc:resource-key:9e7aec27-3d9c-45a8-b031-ce4e0aa985e9" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - guid                  = "9e7aec27-3d9c-45a8-b031-ce4e0aa985e9" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - iam_compatible        = true -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - id                    = "crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc:resource-key:9e7aec27-3d9c-45a8-b031-ce4e0aa985e9" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - name                  = "cos-key-reader-test" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - resource_group_id     = "4bef339d68b14906b5008348875d13db" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - resource_instance_id  = "crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc::" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - resource_instance_url = "/v2/resource_instances/7756104a-c094-4834-8e13-ee07be2518bc" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - role                  = "Reader" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - source_crn            = "crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc::" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - state                 = "active" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - status                = "active" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - updated_at            = "2023-09-16T09:18:27.938Z" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:       - url                   = "/v2/resource_keys/9e7aec27-3d9c-45a8-b031-ce4e0aa985e9" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:     }
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66: Plan: 0 to add, 0 to change, 4 to destroy.
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66: Changes to Outputs:
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:   - bucket_name         = "cos-bucket-tf-test" -> null
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:   - key_object_writer   = (sensitive value)
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:   - key_reader          = (sensitive value)
TestInfraCOSExample 2023-09-16T10:18:46+01:00 logger.go:66:   - service_instance_id = "7756104a-c094-4834-8e13-ee07be2518bc" -> null
TestInfraCOSExample 2023-09-16T10:18:47+01:00 logger.go:66: module.cos.ibm_resource_key.key_cos_object_writer: Destroying... [id=crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc:resource-key:2cc846ff-c393-4e6c-8e06-d18845cbb79b]
TestInfraCOSExample 2023-09-16T10:18:47+01:00 logger.go:66: module.cos.ibm_resource_key.key_cos_reader: Destroying... [id=crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc:resource-key:9e7aec27-3d9c-45a8-b031-ce4e0aa985e9]
TestInfraCOSExample 2023-09-16T10:18:47+01:00 logger.go:66: module.cos.ibm_cos_bucket.bucket_test: Destroying... [id=crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc:bucket:cos-bucket-tf-test:meta:rl:eu-gb:public]
TestInfraCOSExample 2023-09-16T10:18:48+01:00 logger.go:66: module.cos.ibm_cos_bucket.bucket_test: Destruction complete after 1s
TestInfraCOSExample 2023-09-16T10:18:49+01:00 logger.go:66: module.cos.ibm_resource_key.key_cos_reader: Destruction complete after 2s
TestInfraCOSExample 2023-09-16T10:18:50+01:00 logger.go:66: module.cos.ibm_resource_key.key_cos_object_writer: Destruction complete after 2s
TestInfraCOSExample 2023-09-16T10:18:50+01:00 logger.go:66: module.cos.ibm_resource_instance.resource_instance_cos_test: Destroying... [id=crn:v1:bluemix:public:cloud-object-storage:global:a/7f57936089fa4570adc362a257bb04a0:7756104a-c094-4834-8e13-ee07be2518bc::]
TestInfraCOSExample 2023-09-16T10:19:00+01:00 logger.go:66: module.cos.ibm_resource_instance.resource_instance_cos_test: Still destroying... [id=crn:v1:bluemix:public:cloud-object-stor...7756104a-c094-4834-8e13-ee07be2518bc::, 10s elapsed]
TestInfraCOSExample 2023-09-16T10:19:04+01:00 logger.go:66: module.cos.ibm_resource_instance.resource_instance_cos_test: Destruction complete after 14s
TestInfraCOSExample 2023-09-16T10:19:04+01:00 logger.go:66: 
TestInfraCOSExample 2023-09-16T10:19:04+01:00 logger.go:66: Destroy complete! Resources: 4 destroyed.
TestInfraCOSExample 2023-09-16T10:19:04+01:00 logger.go:66: 
PASS
ok      github.com/SRodi/tf-bg-1        67.734s
```



