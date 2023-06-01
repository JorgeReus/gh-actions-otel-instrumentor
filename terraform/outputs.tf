output "webhook_secret" {
  sensitive   = true
  value       = random_password.webhook_secret.result
  description = "The webhook secret for the github app"
}

output "url" {
  value = aws_lambda_function_url.this.function_url
}
