from django.test import TestCase
from django.contrib.auth.models import User
from conversions.models import ConversionJob, SupportedFormat
from rest_framework.test import APIClient
import uuid

class ConversionTests(TestCase):
    def setUp(self):
        self.client = APIClient()
        self.user = User.objects.create_user(
            username='testuser',
            email='test@example.com',
            password='testpass123'
        )
        self.client.force_authenticate(user=self.user)
        
        # Create test formats
        self.jpeg_format = SupportedFormat.objects.create(
            name='JPEG',
            extension='jpg',
            category='image',
            mime_type='image/jpeg'
        )
        self.png_format = SupportedFormat.objects.create(
            name='PNG',
            extension='png',
            category='image',
            mime_type='image/png'
        )

    def test_conversion_job_creation(self):
        """Test ConversionJob model creation"""
        job = ConversionJob.objects.create(
            user=self.user,
            original_filename='test.jpg',
            file_size=1024,
            source_format=self.jpeg_format,
            target_format=self.png_format,
            input_file_path='uploads/test.jpg'
        )
        
        self.assertIsInstance(job.job_id, uuid.UUID)
        self.assertEqual(job.status, 'pending')
        self.assertEqual(job.user, self.user)

    def test_supported_format_model(self):
        """Test SupportedFormat model"""
        self.assertEqual(self.jpeg_format.name, 'JPEG')
        self.assertEqual(self.jpeg_format.category, 'image')
        self.assertEqual(str(self.jpeg_format), 'JPEG (image)')

    def test_conversion_job_str_representation(self):
        """Test ConversionJob string representation"""
        job = ConversionJob.objects.create(
            user=self.user,
            original_filename='test.jpg',
            file_size=1024,
            source_format=self.jpeg_format,
            target_format=self.png_format,
            input_file_path='uploads/test.jpg'
        )
        
        expected_str = f"Job {job.job_id}: test.jpg"
        self.assertEqual(str(job), expected_str)