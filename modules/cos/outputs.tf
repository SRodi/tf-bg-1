output "key_writer"{
    value = ibm_resource_key.key_cos_writer.credentials["apikey"]
    sensitive = true
}

output "key_reader"{
    value = ibm_resource_key.key_cos_reader.credentials["apikey"]
    sensitive = true
}

output "service_instance_id"{
    value = ibm_resource_instance.resource_instance_cos_test.guid
}

output "bucket_name"{
    value = ibm_cos_bucket.bucket_test.bucket_name
}