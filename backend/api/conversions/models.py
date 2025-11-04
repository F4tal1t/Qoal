from django.db import models
from django.contrib.auth.models import User
import uuid

class SupportedFormat(models.Model):
    CATEGORY_CHOICES = [
        ('image', 'Image'), ('video', 'Video'), ('audio', 'Audio'),
        ('document', 'Document'), ('archive', 'Archive'),
    ]
    
    name = models.CharField(max_length=50, unique=True)
    extension = models.CharField(max_length=10)
    category = models.CharField(max_length=20, choices=CATEGORY_CHOICES)
    mime_type = models.CharField(max_length=100)
    icon_filename = models.CharField(max_length=100, blank=True)
    
    def __str__(self):
        return f"{self.name} ({self.category})"

class ConversionJob(models.Model):
    STATUS_CHOICES = [
        ('pending', 'Pending'), ('processing', 'Processing'),
        ('completed', 'Completed'), ('failed', 'Failed'),
    ]
    
    job_id = models.UUIDField(default=uuid.uuid4, unique=True)
    user = models.ForeignKey(User, on_delete=models.CASCADE)
    original_filename = models.CharField(max_length=255)
    file_size = models.BigIntegerField()
    
    source_format = models.ForeignKey(SupportedFormat, on_delete=models.CASCADE, related_name='source_jobs')
    target_format = models.ForeignKey(SupportedFormat, on_delete=models.CASCADE, related_name='target_jobs')
    
    status = models.CharField(max_length=20, choices=STATUS_CHOICES, default='pending')
    input_file_path = models.CharField(max_length=500)
    output_file_path = models.CharField(max_length=500, blank=True)
    
    created_at = models.DateTimeField(auto_now_add=True)
    
    def __str__(self):
        return f"Job {self.job_id}: {self.original_filename}"