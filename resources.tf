resource "nomadutility_acl_bootstrap" "acl" {}

output "token" {
  value = nomadutility_acl_bootstrap.acl
}
