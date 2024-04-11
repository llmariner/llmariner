provider "aws" {
  region  = var.region
  profile = var.profile
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = [var.image_name]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_key_pair" "ssh" {
  key_name   = "key_pair-${var.project_name}"
  public_key = file(var.public_key_path)
}

resource "aws_security_group" "allow_ssh" {
  name        = "allow_ssh-${var.project_name}"
  description = "Allow ssh traffic on port 22 from the specified IP addresses"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [var.ssh_ip_range]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "server" {
  ami                    = var.ami == "" ? data.aws_ami.ubuntu.id : var.ami
  instance_type          = var.instance_type
  key_name               = aws_key_pair.ssh.id
  vpc_security_group_ids = [aws_security_group.allow_ssh.id]

  tags = {
    Name = var.project_name
  }

  root_block_device {
    volume_size = var.volume_size
    volume_type = var.volume_type
  }
}

resource "local_file" "inventory" {
  depends_on = [aws_instance.server]
  content = templatefile("inventory.tftpl",
    {
      server           = aws_instance.server
      private_key_path = var.private_key_path
    }
  )
  filename = "./inventory.ini"
}

resource "null_resource" "provisioning" {
  depends_on = [local_file.inventory]
  provisioner "remote-exec" {
    connection {
      host        = aws_instance.server.public_ip
      user        = var.instance_user
      private_key = file(var.private_key_path)
    }
    inline = ["echo 'server is ready'"]
  }
  provisioner "local-exec" {
    command = "ansible-playbook -i inventory.ini playbook.yml"
  }
}
