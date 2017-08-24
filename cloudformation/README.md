# CloudFormation templates for Cloudfrontbeat

These CloudFormation templates has been created to help you to provision your AWS environment. The templates will automatically create the SQS queue and an S3 bucket for the CloudFront logs. The S3 bucket will notify the queue when a new object / file has been created inside the bucket.

## How to use the templates

The templates has to first be uploaded to an S3 bucket and then you need to create a new CloudFormation stack from the AWS console.

1. Create a bucket for the templates, you can use the [AWS CLI](https://aws.amazon.com/cli/) for that. E.g. `aws s3api create-bucket --bucket my-cloudfrontbeat --region eu-west-1 --create-bucket-configuration LocationConstraint=eu-west-1` Note: S3 buckets are unique and you might want to use some other AWS region.
2. Upload the templates to the bucket with the following command `aws s3 cp . s3://my-cloudfrontbeat/cloudfrontbeat/ --recursive`
3. Navigate to CloudFormation section in your AWS console and click in "create new stack".
4. Select "Specify an Amazon S3 template URL" on the "Select Template" page and insert the following url `https://my-cloudfrontbeat.s3.amazonaws.com/cloudfrontbeat/main.yaml` Note: you need to change the bucket part of the URL.
5. Fill Stack Name and all parameters and continue the wizard till the end.
