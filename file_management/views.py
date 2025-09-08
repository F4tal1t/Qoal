from rest_framework.decorators import api_view, parser_classes
from rest_framework.parsers import MultiPartParser
from rest_framework.response import Response
from rest_framework import status
# from .storage import S3FileManager
from conversions.models import ConversionJob, SupportedFormat
from .tasks import process_conversion_job
import magic
import uuid
from datetime import datetime
from django.core.cache import cache
from .storage import S3FileManager

@api_view(['POST'])
@parser_classes([MultiPartParser])
def upload_file(request):
    """Upload file to S3 and create conversion job"""
    uploaded_file = request.FILES.get('file')
    target_format = request.data.get('target_format')
    
    if not uploaded_file:
        return Response({'error': 'No file provided'}, status=400)
    
    # Validate file type
    file_content = uploaded_file.read(1024)
    uploaded_file.seek(0)
    mime_type = magic.from_buffer(file_content, mime=True)
    
    try:
        source_format = SupportedFormat.objects.get(mime_type=mime_type)
    except SupportedFormat.DoesNotExist:
        return Response({'error': f'Unsupported file type: {mime_type}'}, status=400)
    
    # Get target format
    try:
        target_format_obj = SupportedFormat.objects.get(name=target_format)
    except SupportedFormat.DoesNotExist:
        target_format_obj = source_format  # Default to same format
    
    try:
        # Upload to S3
        s3_manager = S3FileManager()
        file_key = f"uploads/{request.user.id}/{uuid.uuid4()}/{uploaded_file.name}"
        print(f"DEBUG: Attempting to upload file {uploaded_file.name} to {file_key}")  # Debug
        try:
            s3_path = s3_manager.upload_file(uploaded_file, file_key)
            print(f"DEBUG: Successfully uploaded file to {s3_path}")  # Debug
        except Exception as upload_error:
            print(f"DEBUG: S3 upload failed: {str(upload_error)}")  # Debug
            raise
        
        # Create conversion job
        print(f"DEBUG: Creating conversion job for file {uploaded_file.name}")  # Debug
        job = ConversionJob.objects.create(
            user=request.user,
            original_filename=uploaded_file.name,
            file_size=uploaded_file.size,
            source_format=source_format,
            target_format=target_format_obj,
            input_file_path=file_key,
            status='pending'
        )
        
        # Queue conversion job
        try:
            process_conversion_job.delay(str(job.job_id))
        except Exception as e:
            # Clean up uploaded file if queuing fails
            s3_manager.delete_file(file_key)
            job.delete()
            return Response({'error': f'Failed to queue job: {str(e)}'}, status=500)
            
    except Exception as e:
        return Response({'error': f'File upload failed: {str(e)}'}, status=500)
    
    return Response({
        'job_id': str(job.job_id),
        'status': job.status,
        'original_filename': job.original_filename,
        'source_format': source_format.name,
        'target_format': target_format_obj.name,
        'estimated_time': 30  # seconds
    })

@api_view(['GET'])
def job_status(request, job_id):
    """Get conversion job status"""
    try:
        job = ConversionJob.objects.get(job_id=job_id, user=request.user)
        
        response_data = {
            'job_id': str(job.job_id),
            'status': job.status,
            'progress': cache.get(f'job_progress_{job.job_id}', 0),
            'original_filename': job.original_filename,
            'source_format': job.source_format.name,
            'target_format': job.target_format.name,
        }
        
        # Add download URL if completed
        if job.status == 'completed' and job.output_file_path:
            s3_manager = S3FileManager()
            response_data['download_url'] = s3_manager.generate_download_url(job.output_file_path)
        
        return Response(response_data)
        
    except ConversionJob.DoesNotExist:
        return Response({'error': 'Job not found'}, status=404)