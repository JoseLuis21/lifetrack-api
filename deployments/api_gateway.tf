resource "aws_api_gateway_rest_api" "lifeTrack" {
  name = "LifeTrack API"
  description = "Neutrino LifeTrack API"
}

// -- ACM Cert, API Gateway Domain name & Route53 record --

/*
- ISSUE= Currently, we have multiple certs with same domain within the same region
- Most recent strategy doesnt works because non *.api. alternate domain
data "aws_acm_certificate" "cert" {
  domain = "damascus-engineering.com"
  most_recent = true
  statuses = ["ISSUED"]
  tags = {
    "Name": "api"
  }
}*/

resource "aws_api_gateway_domain_name" "domain" {
  certificate_arn = "arn:aws:acm:us-east-1:824699638576:certificate/dcd0e41b-406a-4d0a-a366-4db92d47c012"
  domain_name = "lifetrack.api.damascus-engineering.com"
  security_policy = "TLS_1_2"
  tags = {
    Environment: "prod",
    Name: "neutrino-lifetrack"
  }
}

resource "aws_api_gateway_base_path_mapping" "map" {
  api_id = aws_api_gateway_rest_api.lifeTrack.id
  stage_name = aws_api_gateway_deployment.deploy.stage_name
  domain_name = aws_api_gateway_domain_name.domain.domain_name
  base_path = "live"
}

data "aws_route53_zone" "primary" {
  name = "damascus-engineering.com"
  private_zone = false
}

resource "aws_route53_record" "domain-53" {
  name = aws_api_gateway_domain_name.domain.domain_name
  type = "A"
  zone_id = data.aws_route53_zone.primary.id

  alias {
    evaluate_target_health = true
    name = aws_api_gateway_domain_name.domain.cloudfront_domain_name
    zone_id = aws_api_gateway_domain_name.domain.cloudfront_zone_id
  }
}

// -- API Gateway Proxy config --

resource "aws_api_gateway_resource" "category" {
  rest_api_id = aws_api_gateway_rest_api.lifeTrack.id
  parent_id = aws_api_gateway_rest_api.lifeTrack.root_resource_id
  path_part = "category"
}

resource "aws_api_gateway_method" "category_post" {
  rest_api_id = aws_api_gateway_rest_api.lifeTrack.id
  resource_id = aws_api_gateway_resource.category.id
  http_method = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "lambda" {
  rest_api_id = aws_api_gateway_rest_api.lifeTrack.id
  resource_id = aws_api_gateway_method.category_post.resource_id
  http_method = aws_api_gateway_method.category_post.http_method

  integration_http_method = "POST"
  type = "AWS_PROXY"
  uri = aws_lambda_function.add-category.invoke_arn
}

/*
INFO: Avoid using proxy to keep serverless ecosystem

resource "aws_api_gateway_resource" "category_proxy" {
  rest_api_id = aws_api_gateway_rest_api.lifeTrack.id
  parent_id = aws_api_gateway_resource.category.id
  path_part = "{proxy+}"
}
*/

/*
resource "aws_api_gateway_method" "proxy_root" {
  rest_api_id = aws_api_gateway_rest_api.lifeTrack.id
  resource_id = aws_api_gateway_rest_api.lifeTrack.root_resource_id
  http_method = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "lambda_root" {
  rest_api_id = aws_api_gateway_rest_api.lifeTrack.id
  resource_id = aws_api_gateway_method.proxy_root.resource_id
  http_method = aws_api_gateway_method.proxy_root.http_method

  integration_http_method = "POST"
  type = "AWS_PROXY"
  uri = aws_lambda_function.add-category.invoke_arn
}
*/
resource "aws_api_gateway_deployment" "deploy" {
  depends_on = [
    aws_api_gateway_integration.lambda
  ]
  rest_api_id = aws_api_gateway_rest_api.lifeTrack.id
  stage_name = "live"
  stage_description = "Production stage"
}

output "base_url" {
  value = aws_api_gateway_deployment.deploy.invoke_url
}
