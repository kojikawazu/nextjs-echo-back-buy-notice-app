# ---------------------------------------------
# Variables
# ---------------------------------------------
variable "project" {
  type = string
}

variable "environment" {
  type = string
}

variable "region" {
  type = string
}

variable "vpc_address" {
  type = string
}

variable "public_1a_address" {
  type = string
}

variable "public_1c_address" {
  type = string
}

variable "private_1a_address" {
  type = string
}

variable "private_1c_address" {
  type = string
}

variable "default_route" {
  type = string
}

variable "domain" {
  type = string
}

variable "http_port" {
  type = number
}

variable "https_port" {
  type = number
}

variable "ecs_port" {
  type = number
}

variable "ecs_ws_port" {
  type = number
}

variable "redis_port" {
  type = number
}

variable "db_port" {
  type = number
}

variable "ecs_ws_health_port" {
  type = number
}

variable "elb_health_url" {
  type = string
}

variable "elb_ws_health_url" {
  type = string
}

variable "ecs_container_name" {
  type = string
}

variable "ecr_repository_name" {
  type = string
}

variable "ecs_cpu" {
  type = number
}

variable "ecs_memory" {
  type = number
}

variable "ecs_desired_count" {
  type = number
}

variable "supabase_url" {
  type = string
}

variable "allowed_cors_address" {
  type = string
}

variable "jwt_secret_key" {
  type = string
}

