terraform {
  required_version = ">= 1.6"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    # Populated via -backend-config or terraform.tfbackend file
    # bucket = "your-tf-state-bucket"
    # key    = "design-youtube/terraform.tfstate"
    # region = "us-east-1"
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "design-youtube"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  }
}
