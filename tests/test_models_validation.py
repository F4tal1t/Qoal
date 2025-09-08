from django.test import TestCase
from django.contrib.auth.models import User
from django.core.exceptions import ValidationError
from conversions.models import ConversionJob, SupportedFormat
import uuid

class ModelValidationTests(TestCase):
    def setUp(self):
        self.user = User.objects.create_user(
            username='testuser',
            email='test@example.com',
            password='testpass123'
        )
        
        self.jpeg_format = SupportedFormat.objects.create(
            name='JPEG',
            extension='jpg',
            category='image',
            mime_type='image/jpeg'
        )

    def test_supported_format_unique_name(self):
        """Test that SupportedFormat names must be unique"""
        with self.assertRaises(Exception):  # IntegrityError
            SupportedFormat.objects.create(
                name='JPEG',  # Duplicate name
                extension='jpeg',
                category='image',
                mime_type='image/jpeg'
            )

    def test_conversion_job_uuid_generation(self):
        """Test that ConversionJob generates unique UUIDs"""
        job1 = ConversionJob.objects.create(
            user=self.user,
            original_filename='test1.jpg',
            file_size=1024,
            source_format=self.jpeg_format,
            target_format=self.jpeg_format,
            input_file_path='uploads/test1.jpg'
        )
        
        job2 = ConversionJob.objects.create(
            user=self.user,
            original_filename='test2.jpg',
            file_size=2048,
            source_format=self.jpeg_format,
            target_format=self.jpeg_format,
            input_file_path='uploads/test2.jpg'
        )
        
        self.assertNotEqual(job1.job_id, job2.job_id)
        self.assertIsInstance(job1.job_id, uuid.UUID)
        self.assertIsInstance(job2.job_id, uuid.UUID)

    def test_conversion_job_status_choices(self):
        """Test ConversionJob status field validation"""
        job = ConversionJob.objects.create(
            user=self.user,
            original_filename='test.jpg',
            file_size=1024,
            source_format=self.jpeg_format,
            target_format=self.jpeg_format,
            input_file_path='uploads/test.jpg'
        )
        
        # Test valid statuses
        valid_statuses = ['pending', 'processing', 'completed', 'failed']
        for status in valid_statuses:
            job.status = status
            job.save()  # Should not raise exception
            job.refresh_from_db()
            self.assertEqual(job.status, status)

    def test_supported_format_category_choices(self):
        """Test SupportedFormat category validation"""
        valid_categories = ['image', 'video', 'audio', 'document', 'archive']
        
        for category in valid_categories:
            format_obj = SupportedFormat.objects.create(
                name=f'TEST_{category.upper()}',
                extension='test',
                category=category,
                mime_type=f'test/{category}'
            )
            self.assertEqual(format_obj.category, category)

    def test_conversion_job_required_fields(self):
        """Test that required fields are enforced"""
        with self.assertRaises(Exception):  # IntegrityError for missing user
            ConversionJob.objects.create(
                original_filename='test.jpg',
                file_size=1024,
                source_format=self.jpeg_format,
                target_format=self.jpeg_format,
                input_file_path='uploads/test.jpg'
            )