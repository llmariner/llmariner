# aws

This directory contains files to create a virtual machine on AWS with a GPU. It uses Terraform and Ansible to setup the VM.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/install)
- [Ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html)

## Getting started

Create a `local.tfvars` file as follows for deployment configuration:

> [!NOTE]
> See `variables.tf`for other customizable and default values.

```
project_name = "<instance-name> (default: "llm-operator-demo")"
profile      = "<aws-profile>"

public_key_path  = "</path/to/public_key_path>"
private_key_path = "</path/to/private_key_path>"
ssh_ip_range     = "<ingress CIDR block for SSH (default: "0.0.0.0/0")>"
```

Then, run terraform command to initialize and create an instance. (It takes around 10 minutes.)

```
$ terraform init
$ terraform apply -var-file=local.tfvars
```

> [!TIP]
> If you want to run only Ansible playbook, run `ansible-playbook -i inventory.ini playbook.yml`.

## Cleaning up

```
$ terraform destroy -var-file=local.tfvars
```
