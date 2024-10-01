# ---------------------------------------------
# ECS - タスク定義
# ---------------------------------------------
resource "aws_ecs_task_definition" "ecs_task_def" {
  # タスク定義のファミリー名。タスク定義の論理的なグループを形成します。
  family = var.ecs_container_name
  # タスクがAWSサービスとやりとりするために使用するIAMロールのARN。
  execution_role_arn = aws_iam_role.ecs_task_execution_role.arn
  # タスク定義がサポートする起動タイプ。Fargateはサーバレスコンピュートエンジンです。
  requires_compatibilities = ["FARGATE"]
  # タスクに割り当てるCPUの量（単位はvCPU）。
  cpu = var.ecs_cpu
  # タスクに割り当てるメモリの量（単位はMiB）。
  memory = var.ecs_memory
  # タスクのネットワークモード。awsvpcモードでは、各タスクにElastic Network Interfaceが割り当てられます。
  network_mode = "awsvpc"

  depends_on = [aws_elasticache_replication_group.redis]

  # コンテナ定義。タスクが起動するコンテナの設定をJSON形式で記述します。
  container_definitions = jsonencode([
    {
      name  = "${var.ecs_container_name}",
      image = "${var.ecr_repository_name}",
      portMappings = [
        {
          containerPort = var.ecs_port,
          hostPort      = var.ecs_port,
          protocol      = "tcp"
        },
      ],
      environment = [
        {
          name  = "REDIS_URL",
          value = "${aws_elasticache_replication_group.redis.primary_endpoint_address}:${var.redis_port}"
        }
      ],
      logConfiguration = {
        logDriver = "awslogs",
        options = {
          awslogs-group         = aws_cloudwatch_log_group.ecs_log_group.name,
          awslogs-region        = "${var.region}",
          awslogs-stream-prefix = "ecs"
        }
      },
      secrets = [
        {
          name      = "ALLOWED_ORIGINS",
          valueFrom = "${data.aws_secretsmanager_secret.echo_env.arn}:ALLOWED_ORIGINS::"
        },
        {
          name      = "PORT",
          valueFrom = "${data.aws_secretsmanager_secret.echo_env.arn}:PORT::"
        },
        {
          name      = "SUPABASE_URL",
          valueFrom = "${data.aws_secretsmanager_secret.echo_env.arn}:SUPABASE_URL::"
        },
        {
          name      = "JWT_SECRET_KEY",
          valueFrom = "${data.aws_secretsmanager_secret.echo_env.arn}:JWT_SECRET_KEY::"
        }
      ]
    }
  ])
}

# ---------------------------------------------
# ECS - クラスタ
# ---------------------------------------------
resource "aws_ecs_cluster" "ecs_cluster" {
  name = "${var.project}-${var.environment}-ecs-cluster"
}

# ---------------------------------------------
# ECS - サービス
# ---------------------------------------------
resource "aws_ecs_service" "ecs_service" {
  name = "${var.project}-${var.environment}-ecs-service"
  # このサービスが属するECSクラスタのID。
  cluster = aws_ecs_cluster.ecs_cluster.id
  # サービスで使用する起動タイプ。ここではFargateを指定しています。
  launch_type = "FARGATE"
  # サービスが維持すべきタスクの希望数。
  desired_count = var.ecs_desired_count
  # このサービスで使用するタスク定義のARN。
  task_definition = aws_ecs_task_definition.ecs_task_def.arn

  # ネットワーク設定。サービスによって起動されるタスクが使用するVPC設定を指定します。
  network_configuration {
    subnets = [
      aws_subnet.private_subnet_1a.id,
      aws_subnet.private_subnet_1c.id
    ]
    security_groups = [
      aws_security_group.ecs_sg.id
    ]
    assign_public_ip = false
  }

  # ロードバランサー設定。サービスによって起動されるタスクをロードバランサーに登録
  load_balancer {
    target_group_arn = aws_lb_target_group.alb_target_group.arn
    container_name   = var.ecs_container_name
    container_port   = var.ecs_port
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.alb_target_group_ws.arn
    container_name   = var.ecs_container_name
    container_port   = var.ecs_ws_port
  }
}
