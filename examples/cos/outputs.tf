output "key_object_writer" {
  value     = module.cos.key_object_writer
  sensitive = true
}

output "key_reader" {
  value     = module.cos.key_reader
  sensitive = true
}

output "service_instance_id" {
  value = module.cos.service_instance_id
}

output "bucket_name" {
  value = module.cos.bucket_name
}