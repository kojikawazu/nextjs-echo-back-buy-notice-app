# ---------------------------------------------
# Certificate
# ---------------------------------------------
# for tokyo region
resource "aws_acm_certificate" "tokyo_cert" {
  domain_name       = "*.${var.domain}"
  validation_method = "DNS"

  tags = {
    Name    = "${var.project}-${var.environment}-wildcard-sslcert"
    Project = var.project
    Env     = var.environment
  }

  lifecycle {
    create_before_destroy = true
  }

  depends_on = [
    aws_route53_zone.route53_zone
  ]
}

# ---------------------------------------------
# Route53 - acm dns resolve
# ---------------------------------------------
# ACMがSSL/TLS証明書の検証（所有権確認）を行うときに必要となる、
# ACMからの受け取ったCNAMEレコードをRoute53に自動設定するためのリソース
resource "aws_route53_record" "route53_acm_dns_resolve" {
  for_each = {
    for dvo in aws_acm_certificate.tokyo_cert.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      type   = dvo.resource_record_type
      record = dvo.resource_record_value
    }
  }

  allow_overwrite = true
  zone_id         = aws_route53_zone.route53_zone.id
  name            = each.value.name
  type            = each.value.type
  ttl             = 600
  records         = [each.value.record]
}

# ---------------------------------------------
# Certificate - validate
# ---------------------------------------------
# 証明書検証（DNS検証）が適切に検証されたことを確認するリソース
# 検証が完了するまでは　Terraform apply　のプロセスは保留となり完了しない
resource "aws_acm_certificate_validation" "cert_valid" {
  certificate_arn         = aws_acm_certificate.tokyo_cert.arn
  validation_record_fqdns = [for record in aws_route53_record.route53_acm_dns_resolve : record.fqdn]
}
