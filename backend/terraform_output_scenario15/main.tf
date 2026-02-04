provider "aws" {
  region = "us-east-1"
}

resource "aws_s3_bucket" "my_bucket" {
  bucket = "my-bucket"
  tags = {
    Name  =  "my-bucket"
  }
}

resource "aws_iam_policy" "policy-lambda-1-s3-3" {
  description = "Policy for edge edge-lambda-1-s3-3"
  name        = "policy-lambda-1-s3-3"
  path        = "/"
  policy      = "{\"Statement\":[{\"Action\":[\"s3:GetObject\",\"s3:GetObjectVersion\",\"s3:ListBucket\",\"s3:ListBucketVersions\"],\"Effect\":\"Allow\",\"Resource\":\"*\"}],\"Version\":\"2012-10-17\"}"
}

resource "aws_lambda_function" "my_function" {
  filename      = "function.zip"
  function_name = "my-function"
  handler       = "index.handler"
  memory_size   = 128
  role          = "arn:aws:iam::123456789012:role/service-role/role"
  runtime       = "nodejs18.x"
  tags = {
    Name  =  "my-function"
  }
}

resource "aws_iam_role_policy_attachment" "attach-lambda-1-s3-3" {
  policy_arn = aws_iam_policy.policy-lambda-1-s3-3.arn
  role       = "dddd"
}

