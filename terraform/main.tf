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

resource "aws_secretsmanager_secret" "notification_url" {
  name        = "gha-instrumentor/notification-url"
  description = "Notification URL for the lambda function"
}

resource "aws_secretsmanager_secret_version" "notification_url" {
  secret_id     = aws_secretsmanager_secret.notification_url.id
  secret_string = var.notification_url
}

resource "aws_iam_role" "lambda" {
  name               = var.lambda_function_name
  assume_role_policy = data.aws_iam_policy_document.lambda.json
}

resource "aws_iam_role_policy" "lambda" {
  name = aws_iam_role.lambda.name
  role = aws_iam_role.lambda.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "arn:aws:logs:*:*:*"
      },
      {
        Effect = "Allow",
        Action = [
          "xray:PutTraceSegments",
          "xray:PutTelemetryRecords",
          "xray:GetSamplingRules",
          "xray:GetSamplingTargets",
          "xray:GetSamplingStatisticSummaries"
        ]
        Resource = "*"
      },
      {
        Effect = "Allow",
        Action = [
          "secretsmanager:GetSecretValue",
          "xray:PutTelemetryRecords",
          "xray:GetSamplingRules",
          "xray:GetSamplingTargets",
          "xray:GetSamplingStatisticSummaries"
        ]
        Resource = [
          aws_secretsmanager_secret.webhook_secret.arn,
          aws_secretsmanager_secret.notification_url.arn,
        ]
      }
    ]
  })
}

resource "aws_lambda_function_url" "this" {
  function_name      = aws_lambda_function.this.function_name
  authorization_type = "NONE"
}

resource "aws_lambda_function" "this" {
  function_name    = var.lambda_function_name
  role             = aws_iam_role.lambda.arn
  handler          = "bootstrap"
  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
  runtime          = "go1.x"
  tracing_config {
    mode = "Active"
  }
  layers = [
    "arn:aws:lambda:us-east-1:901920570463:layer:aws-otel-collector-amd64-ver-0-62-1:1"
  ]
}

resource "aws_xray_group" "workflows" {
  insights_configuration {
    insights_enabled      = false
    notifications_enabled = false
  }
  group_name        = "GithubWorkflowJobs"
  filter_expression = "http.useragent = \"github-actions/WorkflowJob\""
}
