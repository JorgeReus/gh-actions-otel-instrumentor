variable "notification_url" {
  type        = string
  description = "The notification URL of the lambda function"
  sensitive   = true
}

variable "lambda_function_name" {
  description = "The lambda function name"
  type        = string
  default     = "gh-workflows-instrumentor"
}
