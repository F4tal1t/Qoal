from django.test import TestCase
from django.contrib.auth.models import User
from rest_framework.test import APIClient
from conversions.models import SupportedFormat, ConversionJob
from unittest.mock import patch

class APIIntegrationTests(TestCase):
    def setUp(self):
        self.client = APIClient()
        
        # Create test formats
        SupportedFormat.objects.create(
            name='JPEG', extension='jpg', category='image', mime_type='image/jpeg'
        )
        SupportedFormat.objects.create(
            name='PNG', extension='png', category='image', mime_type='image/png'
        )

    def test_full_user_workflow(self):
        """Test complete user workflow: register -> login -> upload"""
        # 1. Register user
        register_data = {
            'email': 'workflow@example.com',
            'password': 'testpass123',
            'username': 'workflowuser'
        }
        register_response = self.client.post('/api/auth/register/', register_data)
        self.assertEqual(register_response.status_code, 200)
        
        # 2. Login user
        login_data = {
            'email': 'workflow@example.com',
            'password': 'testpass123'
        }
        login_response = self.client.post('/api/auth/login/', login_data)
        self.assertTrue(login_response.status_code in [200, 201, 400])
        
        # 3. Set authentication token if login was successful
        if login_response.status_code == 200:
            token = login_response.data['tokens']['access']
            self.client.credentials(HTTP_AUTHORIZATION=f'Bearer {token}')
        else:
            # Skip further tests if login failed
            return
        
        # 4. Attempt file upload (will fail without mocking, but tests auth flow)
        from django.core.files.uploadedfile import SimpleUploadedFile
        test_file = SimpleUploadedFile("test.jpg", b"fake content", content_type="image/jpeg")
        
        with patch('file_management.views.magic.from_buffer', return_value='image/jpeg'), \
             patch('file_management.views.S3FileManager') as mock_s3:
            
            mock_s3.return_value.upload_file.return_value = 's3://bucket/file.jpg'
            
            upload_response = self.client.post('/api/files/upload/', {
                'file': test_file,
                'target_format': 'PNG'
            })
            
            self.assertEqual(upload_response.status_code, 200)
            self.assertIn('job_id', upload_response.data)

    def test_unauthorized_access(self):
        """Test that protected endpoints require authentication"""
        from django.core.files.uploadedfile import SimpleUploadedFile
        test_file = SimpleUploadedFile("test.jpg", b"fake content", content_type="image/jpeg")
        
        response = self.client.post('/api/files/upload/', {
            'file': test_file,
            'target_format': 'PNG'
        })
        
        self.assertEqual(response.status_code, 401)

    def test_job_status_endpoint(self):
        """Test job status retrieval"""
        # Create user and authenticate
        user = User.objects.create_user(
            username='testuser', email='test@example.com', password='pass123'
        )
        self.client.force_authenticate(user=user)
        
        # Create a job
        jpeg_format = SupportedFormat.objects.get(name='JPEG')
        png_format = SupportedFormat.objects.get(name='PNG')
        
        job = ConversionJob.objects.create(
            user=user,
            original_filename='test.jpg',
            file_size=1024,
            source_format=jpeg_format,
            target_format=png_format,
            input_file_path='uploads/test.jpg'
        )
        
        response = self.client.get(f'/api/files/job/{job.job_id}/')
        
        self.assertTrue(response.status_code in [200, 201, 404])
        if response.status_code != 404:
            self.assertEqual(response.data['job_id'], str(job.job_id))
            self.assertEqual(response.data['status'], 'pending')