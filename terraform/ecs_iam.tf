# ---------------------------------------------
# ECSタスク - IAMロール
# ---------------------------------------------
# ECSタスクがAWSリソースとやり取りするために必要なIAMロールを作成します。
# このロールは、ECSタスクがAWSのサービス(例：CloudWatch LogsやECR)と通信するために使用されます。
resource "aws_iam_role" "ecs_task_execution_role" {
  name = "${var.project}-${var.environment}-ecs-iam-role"

  # ECSタスクがこのロールを引き受ける(AssumeRole)することを許可するポリシーを指定します。
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Effect = "Allow",
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

# ---------------------------------------------
# ECSタスク - IAMロールポリシー(基本)
# ---------------------------------------------
# Amazonが提供するECSタスク実行用の基本的なポリシーをECSタスク実行ロールにアタッチします。
# このポリシーには、ECSタスクが実行に必要な基本的な権限が含まれています。
resource "aws_iam_role_policy_attachment" "ecs_task_execution_role_basic_policy" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# ---------------------------------------------
# ECSタスク - IAMロールポリシー(カスタム)
# ---------------------------------------------
# ECSタスクが実行時に必要とする権限を持つIAMポリシーを定義します。
resource "aws_iam_role_policy" "ecs_task_execution_role_custom_policy" {
  name = "${var.project}-${var.environment}-ecs-task-execution-iam-role-custom-policy"
  role = aws_iam_role.ecs_task_execution_role.id

  # ポリシーの内容を指定します。ECSタスクがCloudWatch Logsにログを送信したり、
  # ECRからコンテナイメージを取得するための権限を定義しています。
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      # CloudWatch Logsにログを送信するための権限を定義します。
      {
        Effect = "Allow",
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        Resource = "*",
      },
      # ECSタスクがECRからコンテナイメージを取得するための権限を定義します。
      {
        Effect = "Allow",
        Action = [
          "ecr:GetAuthorizationToken",
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage"
        ],
        Resource = "*",
      },
      # シークレットマネージャーからシークレットを取得するための権限を定義します。
      {
        Effect = "Allow",
        Action = [
          "secretsmanager:GetSecretValue",
          "secretsmanager:DescribeSecret"
        ],
        Resource = [
          "${aws_secretsmanager_secret.echo_env.arn}",
          "${aws_secretsmanager_secret_version.echo_env.arn}"
        ]
      }
    ]
  })
}
