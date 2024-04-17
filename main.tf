terraform {
    required_providers {
      aws = {
        source = "hashicorp/aws"
        version = "~> 4.16"
      }
    }

    required_version = ">= 1.2.0"
}

provider "aws" {
  region = "us-east-1"
}

module "base-network" {
  source = "./modules/networking"

  cidr_block = "192.168.0.0/16"

  vpc_additional_tags = {
    vpc_tag1 = "tag1",
    vpc_tag2 = "tag2",
  }

  public_subnets = {
    first_public_subnet = {
      availability_zone = "us-east-1a"
      cidr_block        = "192.168.0.0/19"
    }
    second_public_subnet = {
      availability_zone = "us-east-1b"
      cidr_block        = "192.168.32.0/19"
    }
  }

  public_subnets_additional_tags = {
    public_subnet_tag1 = "tag1",
    public_subnet_tag2 = "tag2",
  }

  private_subnets = {
    first_private_subnet = {
      availability_zone = "us-east-1a"
      cidr_block        = "192.168.128.0/19"
    }
    second_private_subnet = {
      availability_zone = "us-east-1b"
      cidr_block        = "192.168.160.0/19"
    }
  }

  private_subnets_additional_tags = {
    private_subnet_tag1 = "tag1",
    private_subnet_tag2 = "tag2",
  }
}