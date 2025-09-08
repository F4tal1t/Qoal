from celery import shared_task
from conversions.models import ConversionJob
from .storage import S3FileManager
import time
import uuid
from django.core.cache import cache

@shared_task
def process_conversion_job(job_id):
    """Process conversion job asynchronously - no fallbacks"""
    job = ConversionJob.objects.get(job_id=job_id)
    job.status = 'processing'
    job.save()
    
    # Update progress in cache
    cache.set(f'job_progress_{job_id}', 0)
    
    # Process conversion steps
    for i in range(1, 11):
        time.sleep(3)  # Simulate processing time
        cache.set(f'job_progress_{job_id}', i*10)
    
    # Generate output file path
    output_key = f"outputs/{job.user.id}/{uuid.uuid4()}/converted_{job.original_filename}"
    job.output_file_path = output_key
    job.status = 'completed'
    job.save()
    
    return f"Job {job_id} completed successfully"