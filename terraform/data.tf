data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = "../src/bin/bootstrap"
  output_path = "lambda.zip"
}

data "aws_iam_policy_document" "lambda" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}
