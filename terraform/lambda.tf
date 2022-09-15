resource "aws_lambda_function" "lambda" {
  function_name = "MoviesLambda"

  s3_bucket = aws_s3_bucket.lambda_deploy_bucket.id
  s3_key    = aws_s3_object.lambda_deploy_bucket_object.key

  runtime = "go1.x"
  handler = "movies_lambda"

  source_code_hash = data.archive_file.lambda_zip.output_base64sha256

  role = aws_iam_role.lambda_execution_role.arn
}

resource "aws_lambda_permission" "allow_cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda.function_name
  principal     = "events.amazonaws.com"
}

resource "aws_lambda_function_url" "function" {
    function_name      = aws_lambda_function.lambda.function_name
    authorization_type = "NONE"
}
