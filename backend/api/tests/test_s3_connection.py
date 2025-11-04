import os
import boto3
from django.conf import settings

# Load environment variables
from dotenv import load_dotenv
load_dotenv()

try:
    # Initialize S3 client
    s3 = boto3.client(
        's3',
        region_name=os.getenv('AWS_REGION'),
        aws_access_key_id=os.getenv('AWS_ACCESS_KEY_ID'),
        aws_secret_access_key=os.getenv('AWS_SECRET_ACCESS_KEY')
    )
    
    # Test connection
    response = s3.list_buckets()
    print("Connection successful! Available buckets:")
    for bucket in response['Buckets']:
        print("- " + bucket['Name'])
    
except Exception as e:
    print("Error connecting to S3: " + str(e))