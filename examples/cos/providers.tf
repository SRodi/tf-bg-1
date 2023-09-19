terraform {
  required_version = ">=1.3.0, <1.6"
  required_providers {
    ibm = {
      source  = "IBM-Cloud/ibm"
      version = "1.57.0"
    }
  }
}