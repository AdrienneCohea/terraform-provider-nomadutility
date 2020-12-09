terraform {
  required_version = ">= 0.14.0"

  required_providers {
    nomadutility = {
      source  = "local/adriennecohea/nomadutility"
      version = "0.0.3"
    }
  }
}
