# ---------------------------------------------
# VPC
# ---------------------------------------------
resource "aws_vpc" "vpc" {
  cidr_block                       = var.vpc_address
  instance_tenancy                 = "default"
  enable_dns_support               = true
  enable_dns_hostnames             = true
  assign_generated_ipv6_cidr_block = false

  tags = {
    Name    = "${var.project}-${var.environment}-vpc"
    Project = var.project
    Env     = var.environment
  }
}

# ---------------------------------------------
# Subnet
# ---------------------------------------------
# public 1a
resource "aws_subnet" "public_subnet_1a" {
  vpc_id                  = aws_vpc.vpc.id
  availability_zone       = "${var.region}a"
  cidr_block              = var.public_1a_address
  map_public_ip_on_launch = true

  tags = {
    Name    = "${var.project}-${var.environment}-pub-subnet-1a"
    Project = var.project
    Env     = var.environment
    Type    = "public"
  }
}

# public 1c
resource "aws_subnet" "public_subnet_1c" {
  vpc_id                  = aws_vpc.vpc.id
  availability_zone       = "${var.region}c"
  cidr_block              = var.public_1c_address
  map_public_ip_on_launch = true

  tags = {
    Name    = "${var.project}-${var.environment}-pub-subnet-1c"
    Project = var.project
    Env     = var.environment
    Type    = "public"
  }
}

# private 1a
resource "aws_subnet" "private_subnet_1a" {
  vpc_id            = aws_vpc.vpc.id
  availability_zone = "${var.region}a"
  cidr_block        = var.private_1a_address

  tags = {
    Name    = "${var.project}-${var.environment}-pri-subnet-1a"
    Project = var.project
    Env     = var.environment
    Type    = "private"
  }
}

# private 1c
resource "aws_subnet" "private_subnet_1c" {
  vpc_id            = aws_vpc.vpc.id
  availability_zone = "${var.region}c"
  cidr_block        = var.private_1c_address

  tags = {
    Name    = "${var.project}-${var.environment}-pri-subnet-1c"
    Project = var.project
    Env     = var.environment
    Type    = "private"
  }
}

# ---------------------------------------------
# Route table
# ---------------------------------------------
resource "aws_route_table" "public_rt" {
  vpc_id = aws_vpc.vpc.id

  tags = {
    Name    = "${var.project}-${var.environment}-pub-rt"
    Project = var.project
    Env     = var.environment
    Type    = "public"
  }
}

resource "aws_route_table_association" "public_rt_1a" {
  route_table_id = aws_route_table.public_rt.id
  subnet_id      = aws_subnet.public_subnet_1a.id
}

resource "aws_route_table_association" "public_rt_1c" {
  route_table_id = aws_route_table.public_rt.id
  subnet_id      = aws_subnet.public_subnet_1c.id
}

resource "aws_route_table" "private_rt_1a" {
  vpc_id = aws_vpc.vpc.id

  tags = {
    Name    = "${var.project}-${var.environment}-pri-rt-1a"
    Project = var.project
    Env     = var.environment
    Type    = "private"
  }
}

resource "aws_route_table_association" "private_rt_1a" {
  route_table_id = aws_route_table.private_rt_1a.id
  subnet_id      = aws_subnet.private_subnet_1a.id
}

resource "aws_route" "private_rt_nat_gw_1a" {
  route_table_id         = aws_route_table.private_rt_1a.id
  destination_cidr_block = var.default_route
  nat_gateway_id         = aws_nat_gateway.nat_gw.id
}

resource "aws_route_table" "private_rt_1c" {
  vpc_id = aws_vpc.vpc.id

  tags = {
    Name    = "${var.project}-${var.environment}-pri-rt-1c"
    Project = var.project
    Env     = var.environment
    Type    = "private"
  }
}

resource "aws_route_table_association" "private_rt_1c" {
  route_table_id = aws_route_table.private_rt_1c.id
  subnet_id      = aws_subnet.private_subnet_1c.id
}

resource "aws_route" "private_rt_nat_gw_1c" {
  route_table_id         = aws_route_table.private_rt_1c.id
  destination_cidr_block = var.default_route
  nat_gateway_id         = aws_nat_gateway.nat_gw_1c.id
}

# ---------------------------------------------
# Internet Gateway
# ---------------------------------------------
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.vpc.id

  tags = {
    Name    = "${var.project}-${var.environment}-igw"
    Project = var.project
    Env     = var.environment
  }
}

resource "aws_route" "public_rt_igw_r" {
  route_table_id         = aws_route_table.public_rt.id
  destination_cidr_block = var.default_route
  gateway_id             = aws_internet_gateway.igw.id
}

# ---------------------------------------------
# NAT Gateway
# ---------------------------------------------
# NAT Gateway 1a
resource "aws_nat_gateway" "nat_gw" {
  allocation_id     = aws_eip.nat_eip.id
  subnet_id         = aws_subnet.public_subnet_1a.id
  connectivity_type = "public"

  tags = {
    Name    = "${var.project}-${var.environment}-nat-gw"
    Project = var.project
    Env     = var.environment
  }
}

# NAT Gateway 1c
resource "aws_nat_gateway" "nat_gw_1c" {
  allocation_id     = aws_eip.nat_eip_1c.id
  subnet_id         = aws_subnet.public_subnet_1c.id
  connectivity_type = "public"

  tags = {
    Name    = "${var.project}-${var.environment}-nat-gw-1c"
    Project = var.project
    Env     = var.environment
  }
}

# ---------------------------------------------
# Elastic IP（EIP）
# ---------------------------------------------
# NAT Gateway EIP
resource "aws_eip" "nat_eip" {
  tags = {
    Name    = "${var.project}-${var.environment}-nat-eip"
    Project = var.project
    Env     = var.environment
  }
}

# NAT Gateway EIP 1c
resource "aws_eip" "nat_eip_1c" {
  tags = {
    Name    = "${var.project}-${var.environment}-nat-eip-1c"
    Project = var.project
    Env     = var.environment
  }
}
