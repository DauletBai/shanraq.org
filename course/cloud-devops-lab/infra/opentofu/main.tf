terraform {
  required_version = ">= 1.9.0"
}

variable "server_name" {
  type        = string
  description = "A stable name for the CloudLab server."
  default     = "cloudlab-01"

  validation {
    condition     = can(regex("^[a-z0-9-]+$", var.server_name))
    error_message = "Use lowercase letters, digits, and hyphens."
  }
}

variable "public_ip" {
  type        = string
  description = "The public IPv4 address created with your chosen provider."
}

variable "ssh_user" {
  type        = string
  description = "The non-root account used by Ansible."
  default     = "deploy"
}

resource "terraform_data" "server_record" {
  input = {
    name     = var.server_name
    address  = var.public_ip
    ssh_user = var.ssh_user
  }

  lifecycle {
    precondition {
      condition     = can(cidrhost("${var.public_ip}/32", 0))
      error_message = "public_ip must be a valid IPv4 address."
    }
  }
}

output "ansible_inventory_line" {
  value = "${terraform_data.server_record.output.name} ansible_host=${terraform_data.server_record.output.address} ansible_user=${terraform_data.server_record.output.ssh_user}"
}
