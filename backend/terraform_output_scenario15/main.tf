provider "aws" {
  region = "us-east-1"
}

resource "aws_s3_bucket" "fcf89be3_a1c1_4f77_9c24_d854805022ef" {
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

resource "aws_lambda_function" "r_6d972cbc_932c_4b57_9fdf_be323cf31f5f" {
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

