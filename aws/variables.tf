variable "project_name" {
  type    = string
  default = "llm-operator-demo"
}

variable "region" {
  type    = string
  default = "us-east-2"
}

variable "profile" {
  type    = string
  default = ""
}

variable "image_name" {
  type    = string
  default = "ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"
}

variable "instance_user" {
  type    = string
  default = "ubuntu"
}

variable "ami" {
  type    = string
  default = ""
}

variable "instance_type" {
  type    = string
  default = "g5.4xlarge"
}

variable "volume_size" {
  type    = number
  default = 500
}

variable "volume_type" {
  type    = string
  default = "gp3"
}

variable "ssh_ip_range" {
  type    = string
  default = "0.0.0.0/0"
}

variable "public_key_path" {
  type    = string
  default = ""
}

variable "private_key_path" {
  type    = string
  default = ""
}
