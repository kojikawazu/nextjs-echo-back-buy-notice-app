# ---------------------------------------------
# Security Group
# ---------------------------------------------
# ALB Security Group
resource "aws_security_group" "alb_sg" {
  name        = "${var.project}-${var.environment}-alb-sg"
  description = "ALB security group"
  vpc_id      = aws_vpc.vpc.id

  tags = {
    Name    = "${var.project}-${var.environment}-alb-sg"
    Project = var.project
    Env     = var.environment
  }
}

resource "aws_security_group_rule" "alb_in_http" {
  security_group_id = aws_security_group.alb_sg.id
  type              = "ingress"
  protocol          = "tcp"
  from_port         = var.http_port
  to_port           = var.http_port
  cidr_blocks       = [var.default_route]
}

resource "aws_security_group_rule" "alb_in_https" {
  security_group_id = aws_security_group.alb_sg.id
  type              = "ingress"
  protocol          = "tcp"
  from_port         = var.https_port
  to_port           = var.https_port
  cidr_blocks       = [var.default_route]
}

resource "aws_security_group_rule" "alb_out_internet" {
  security_group_id = aws_security_group.alb_sg.id
  type              = "egress"
  protocol          = "-1"
  from_port         = 0
  to_port           = 0
  cidr_blocks       = [var.default_route]
}

# ECS Security Group
resource "aws_security_group" "ecs_sg" {
  name        = "${var.project}-${var.environment}-ecs-sg"
  description = "ECS security group"
  vpc_id      = aws_vpc.vpc.id

  tags = {
    Name    = "${var.project}-${var.environment}-ecs-sg"
    Project = var.project
    Env     = var.environment
  }
}

resource "aws_security_group_rule" "ecs_in_alb" {
  security_group_id        = aws_security_group.ecs_sg.id
  type                     = "ingress"
  protocol                 = "tcp"
  from_port                = var.ecs_port
  to_port                  = var.ecs_port
  source_security_group_id = aws_security_group.alb_sg.id
}

resource "aws_security_group_rule" "ecs_out_internet" {
  security_group_id = aws_security_group.ecs_sg.id
  type              = "egress"
  protocol          = "-1"
  from_port         = 0
  to_port           = 0
  cidr_blocks       = [var.default_route]
}

resource "aws_security_group_rule" "ecs_out_supabase" {
  security_group_id = aws_security_group.ecs_sg.id
  type              = "egress"
  protocol          = "tcp"
  from_port         = var.db_port
  to_port           = var.db_port
  cidr_blocks       = [var.default_route]
}

# ElastiCache Security Group
resource "aws_security_group" "redis_sg" {
  name        = "${var.project}-${var.environment}-redis-sg"
  description = "Redis security group"
  vpc_id      = aws_vpc.vpc.id

  tags = {
    Name    = "${var.project}-${var.environment}-redis-sg"
    Project = var.project
    Env     = var.environment
  }
}

resource "aws_security_group_rule" "redis_in_ecs" {
  security_group_id        = aws_security_group.redis_sg.id
  type                     = "ingress"
  protocol                 = "tcp"
  from_port                = var.redis_port
  to_port                  = var.redis_port
  source_security_group_id = aws_security_group.ecs_sg.id
}
