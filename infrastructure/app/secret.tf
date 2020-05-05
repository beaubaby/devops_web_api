data "aws_iam_policy_document" "rds_kms_key_policy" {
  statement {
    sid       = "Enable IAM User Permissions"
    effect    = "Allow"
    resources = ["*"]
    actions   = ["kms:*"]

    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${var.account_id}:root"]
    }
  }

  statement {
    sid       = "Allow use of the key from tools account to copy rds snapshots"
    effect    = "Allow"
    resources = ["*"]

    actions = [
      "kms:Decrypt",
      "kms:DescribeKey",
    ]

    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${var.account_id}:root"]
    }
  }

  statement {
    sid       = "Allow attachment of persistent resources"
    effect    = "Allow"
    resources = ["*"]

    actions = [
      "kms:CreateGrant",
      "kms:ListGrants",
      "kms:RevokeGrant",
    ]

    condition {
      test     = "Bool"
      variable = "kms:GrantIsForAWSResource"
      values   = ["true"]
    }

    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${var.account_id}:root"]
    }
  }
}

resource "aws_kms_key" "loan_kms_key_rds" {
  tags = {
    Name        = "Loan Eligibility DB Encryption Key"
    App         = "Loan Eligibility"
    Component   = "DB"
    Environment = "${var.environment_name}"
  }

  policy = data.aws_iam_policy_document.rds_kms_key_policy.json
}

resource "aws_kms_alias" "loan_kms_alias_rds" {
  target_key_id = aws_kms_key.loan_kms_key_rds.id
  name = "alias/${var.environment_name}/loan-kms-alias-rds"
}