terraform {
  backend "s3" {
    bucket     = "terraform-state-storage-586877430255"
    lock_table = "terraform-state-lock-586877430255"
    region     = "us-west-2"

    // THIS MUST BE UNIQUE
    key = "control-keys.tfstate"
  }
}

provider "aws" {
  region = "us-west-2"
}

data "aws_ssm_parameter" "eks_cluster_endpoint" {
  name = "/eks/av-cluster-endpoint"
}

provider "kubernetes" {
  host = data.aws_ssm_parameter.eks_cluster_endpoint.value
}

data "aws_ssm_parameter" "couch_address" {
  name = "/env/couch-address"
}

data "aws_ssm_parameter" "couch_username" {
  name = "/env/couch-username"
}

data "aws_ssm_parameter" "couch_password" {
  name = "/env/couch-password"
}

module "prd" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "control-keys"
  image          = "docker.pkg.github.com/byuoitav/control-keys/control-keys-dev"
  image_version  = "821030e"
  container_port = 8029
  repo_url       = "https://github.com/byuoitav/control-keys"

  // optional
  image_pull_secret = "github-docker-registry"
  public_urls       = ["control-keys.av.byu.edu"]
  container_env = {
    "DB_ADDRESS"       = data.aws_ssm_parameter.couch_address.value
    "DB_USERNAME"      = data.aws_ssm_parameter.couch_username.value
    "DB_PASSWORD"      = data.aws_ssm_parameter.couch_password.value
    "STOP_REPLICATION" = "true"
  }
  container_args = []
  ingress_annotations = {
    "nginx.ingress.kubernetes.io/whitelist-source-range" = "128.187.0.0/16"
  }
  health_check = false
}
