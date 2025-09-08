import boto3
from django.conf import settings
from botocore.exceptions import ClientError
import uuid
from datetime import datetime

class S3FileManager:
    """Manager for S3 file operations"""

    def __init__(self):
        """Initialize S3 client"""
        if settings.DEBUG:
            # Mock implementation for testing
            self.mock_files = {}
        else:
            # Real S3 implementation
            s3_config = {
                'service_name': 's3',
                'aws_access_key_id': settings.AWS_ACCESS_KEY_ID,
                'aws_secret_access_key': settings.AWS_SECRET_ACCESS_KEY,
                'region_name': settings.AWS_S3_REGION_NAME
            }
            
            if settings.AWS_S3_ENDPOINT_URL:
                s3_config['endpoint_url'] = settings.AWS_S3_ENDPOINT_URL
                
            self.s3_client = boto3.client(**s3_config)
            self.bucket_name = settings.AWS_STORAGE_BUCKET_NAME
            
            try:
                self.s3_client.head_bucket(Bucket=self.bucket_name)
            except self.s3_client.exceptions.NoSuchBucket:
                if settings.AWS_S3_ENDPOINT_URL:
                    # For localstack/testing
                    self.s3_client.create_bucket(Bucket=self.bucket_name)
                else:
                    # For real AWS
                    self.s3_client.create_bucket(
                        Bucket=self.bucket_name,
                        CreateBucketConfiguration={
                            'LocationConstraint': settings.AWS_S3_REGION_NAME
                        }
                    )
    
    def upload_file(self, file_obj, file_key):
        """Upload file to S3 with encryption - no fallbacks"""
        self.s3_client.upload_fileobj(
            file_obj,
            self.bucket_name,
            file_key,
            ExtraArgs={
                'ServerSideEncryption': 'AES256',
                'ContentType': file_obj.content_type
            }
        )
        return f"s3://{self.bucket_name}/{file_key}"
    
    def generate_download_url(self, file_key, expiration=3600):
        """Generate presigned URL for file download - no fallbacks"""
        return self.s3_client.generate_presigned_url(
            'get_object',
            Params={'Bucket': self.bucket_name, 'Key': file_key},
            ExpiresIn=expiration
        )
    
    def delete_file(self, file_key):
        """Delete file from S3"""
        try:
            self.s3_client.delete_object(Bucket=self.bucket_name, Key=file_key)
            return True
        except ClientError as e:
            print(f"Error deleting file: {e}")
            return False