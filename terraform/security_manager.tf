# ---------------------------------------------
# Secret Manager - 環境変数
# ---------------------------------------------
resource "aws_secretsmanager_secret" "echo_env" {
  name = "echo-reservation-env-secret"
}

resource "aws_secretsmanager_secret_version" "echo_env" {
  secret_id = aws_secretsmanager_secret.echo_env.id
  secret_string = jsonencode({
    ALLOWED_ORIGINS = "${var.allowed_cors_address}",
    PORT            = "${var.ecs_port}",
    SUPABASE_URL    = "${var.supabase_url}",
    JWT_SECRET_KEY  = "${var.jwt_secret_key}",
  })

  depends_on = [aws_elasticache_replication_group.redis]
}

data "aws_secretsmanager_secret" "echo_env" {
  name = aws_secretsmanager_secret.echo_env.name
}

data "aws_secretsmanager_secret_version" "echo_env" {
  secret_id = data.aws_secretsmanager_secret.echo_env.id
}
