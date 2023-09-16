data "ibm_resource_group" "resource_group_default" {
  name = "Default"
}

resource "ibm_resource_instance" "resource_instance_cos_test" {
  name              = "cos-instance-test"
  service           = "cloud-object-storage"
  plan              = "lite"
  location          = "global"
  resource_group_id = data.ibm_resource_group.resource_group_default.id
}

resource "ibm_cos_bucket" "bucket_test" {
  bucket_name          = "cos-bucket-tf-test"
  resource_instance_id = ibm_resource_instance.resource_instance_cos_test.id
  region_location      = "eu-gb"
  storage_class        = "standard"
}

resource "ibm_resource_key" "key_cos_object_writer" {
  name                 = "cos-key-object-writer-test"
  role                 = "Object Writer"
  resource_instance_id = ibm_resource_instance.resource_instance_cos_test.id
}

resource "ibm_resource_key" "key_cos_reader" {
  name                 = "cos-key-reader-test"
  role                 = "Reader"
  resource_instance_id = ibm_resource_instance.resource_instance_cos_test.id
}

