# ---------------------------------------------
# ECSタスク用 - CloudWatch ロググループ
# ---------------------------------------------
# ECSタスクからのログを収集するためのCloudWatchロググループを作成します。
# このロググループは、タスクの実行ログを保存し、監視とデバッグに役立てます。
resource "aws_cloudwatch_log_group" "ecs_log_group" {
  name = "/ecs/${var.project}-${var.environment}"
}
