resource "aws_ecr_repository" "service" {
  name = "loan-eligibility-service"
}

resource "aws_ecr_repository_policy" "service" {
  policy     = data.aws_iam_policy_document.multi_account_access.json
  repository = aws_ecr_repository.service.name
}

resource "aws_ecr_lifecycle_policy" "service" {
  repository = aws_ecr_repository.service.name
  policy     = file("ecr-lifecycle-policy.json")
}

data "aws_iam_policy_document" "multi_account_access" {
  statement {
    sid    = "AllowPullingFromDevAndProd"
    effect = "Allow"

    actions = [
      "ecr:GetDownloadUrlForLayer",
      "ecr:BatchGetImage",
      "ecr:BatchCheckLayerAvailability",
    ]

    principals {
      type = "AWS"

      identifiers = [
        "arn:aws:iam::${local.dev_account_id}:root",
        "arn:aws:iam::${local.prod_account_id}:root",
      ]
    }
  }

  statement {
    sid    = "AllowPushingAndPullingFromTools"
    effect = "Allow"

    actions = [
      "ecr:UploadLayerPart",
      "ecr:PutImage",
      "ecr:InitiateLayerUpload",
      "ecr:GetDownloadUrlForLayer",
      "ecr:CompleteLayerUpload",
      "ecr:BatchGetImage",
      "ecr:BatchCheckLayerAvailability",
    ]

    principals {
      type = "AWS"

      identifiers = [
        "arn:aws:iam::${local.tools_account_id}:root",
      ]
    }
  }
}
