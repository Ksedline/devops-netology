output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "caller_user" {
  value = data.aws_caller_identity.current.user_id
}

output "caller_arn" {
  value = data.aws_caller_identity.current.arn
}

output "region" {
  value = data.aws_region.current.name
}
# Приватный IP ec2 инстансы
output "private_ip" {
  value = aws_instance.netology_test.*.private_ip
}

output "subnet_id" {
  value = aws_instance.netology_test.*.subnet_id
}
