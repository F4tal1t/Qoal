from django.test import TestCase
from django.contrib.auth.models import User
from django.core.files.uploadedfile import SimpleUploadedFile
from rest_framework.test import APIClient
from conversions.models import SupportedFormat
from unittest.mock import patch, MagicMock
import io
import uuid

class FileManagementTests(TestCase):
    def setUp(self):
        self.client = APIClient()
        self.user = User.objects.create_user(
            username='testuser',
            email='test@example.com',
            password='testpass123'
        )
        self.client.force_authenticate(user=self.user)
        
        # Create test format
        self.jpeg_format = SupportedFormat.objects.create(
            name='JPEG',
            extension='jpg',
            category='image',
            mime_type='image/jpeg'
        )

    def create_test_file(self):
        """Create a test file for upload"""
        file_content = b"fake image content"
        return SimpleUploadedFile(
            "test.jpg",
            file_content,
            content_type="image/jpeg"
        )

    @patch('file_management.views.magic.from_buffer')
    @patch('file_management.views.S3FileManager')
    @patch('file_management.views.SupportedFormat.objects.get')
    @patch('conversions.models.ConversionJob.objects.create')
    @patch('file_management.tasks.process_conversion_job.delay')
    def test_file_upload_success(self, mock_process_job_delay, mock_create, mock_format_get, mock_s3, mock_magic):
        """Test successful file upload"""
        # Mock magic to return JPEG mime type
        mock_magic.return_value = 'image/jpeg'
        
        # Mock format lookup
        mock_source_format = MagicMock()
        mock_source_format.name = 'JPEG'
        mock_target_format = MagicMock()
        mock_target_format.name = 'PNG'
        mock_format_get.side_effect = lambda *args, **kwargs: (
            mock_source_format if kwargs.get('mime_type') == 'image/jpeg' 
            else mock_target_format
        )
        
        # Mock S3 upload
        mock_s3_instance = MagicMock()
        mock_s3_instance.upload_fileobj.return_value = None
        mock_s3_instance.generate_presigned_url.return_value = 'http://mock.url'
        mock_s3_instance.head_bucket.return_value = True
        mock_s3.return_value = mock_s3_instance
        
        # Mock file object
        mock_file = MagicMock()
        mock_file.content_type = 'image/jpeg'
        mock_file.name = 'test.jpg'
        mock_file.size = len(b"fake image content")
        
        # Mock job creation
        mock_job = MagicMock()
        mock_job.job_id = uuid.UUID('12345678123456781234567812345678')
        mock_job.status = 'pending'
        mock_job.original_filename = 'test.jpg'
        mock_job.source_format = mock_source_format
        mock_job.target_format = mock_target_format
        mock_create.return_value = mock_job
        
        # Mock celery task
        mock_process_job_delay.return_value = None
        
        test_file = self.create_test_file()
        
        response = self.client.post('/api/files/upload/', {
            'file': test_file,
            'target_format': 'PNG'
        }, format='multipart')
        
        print(f"DEBUG: Response status: {response.status_code}")  # Debug
        print(f"DEBUG: Response data: {response.data}")  # Debug
        
        self.assertEqual(response.status_code, 200)
        self.assertIn('job_id', response.data)
        self.assertIn('status', response.data)
        self.assertIn('source_format', response.data)
        self.assertIn('target_format', response.data)
        
        test_file = self.create_test_file()
        
        response = self.client.post('/api/files/upload/', {
            'file': test_file,
            'target_format': 'PNG'
        })
        
        self.assertEqual(response.status_code, 200)
        self.assertIn('job_id', response.data)
        self.assertEqual(response.data['source_format'], 'JPEG')

    def test_file_upload_no_file(self):
        """Test file upload without file"""
        response = self.client.post('/api/files/upload/', {
            'target_format': 'PNG'
        })
        
        self.assertEqual(response.status_code, 400)
        self.assertIn('error', response.data)

    @patch('file_management.views.magic.from_buffer')
    def test_unsupported_file_type(self, mock_magic):
        """Test upload of unsupported file type"""
        mock_magic.return_value = 'application/unknown'
        
        test_file = self.create_test_file()
        
        response = self.client.post('/api/files/upload/', {
            'file': test_file,
            'target_format': 'PNG'
        })
        
        self.assertEqual(response.status_code, 400)
        self.assertIn('error', response.data)