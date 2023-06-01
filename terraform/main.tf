provider "aws" {
  region = "us-east-1"
}

resource "random_password" "webhook_secret" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "aws_secretsmanager_secret" "webhook_secret" {
  name        = "gha-instrumentor/webhook-secret"
  description = "Secure password secret"
}

resource "aws_secretsmanager_secret_version" "webhook_secret" {
  secret_id     = aws_secretsmanager_secret.webhook_secret.id
  secret_string = random_password.webhook_secret.result
}

terraform {
  required_providers {
    aws    = ">= 5.0.1"
    random = ">= 3.5.1"
  }
}

output "webhook_secret" {
  sensitive   = true
  value       = random_password.webhook_secret.result
  description = "The webhook secret for the github app"
}
