from rest_framework.decorators import api_view, parser_classes, permission_classes
from rest_framework.parsers import MultiPartParser
from rest_framework.response import Response
from rest_framework import status
from rest_framework.permissions import IsAuthenticated
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
@permission_classes([IsAuthenticated])
def upload_file(request):
    """Upload file and create conversion job per plan.json requirements"""
    uploaded_file = request.FILES.get('file')
    target_format = request.data.get('target_format')
    conversion_type = request.data.get('conversion_type')  # image, audio, video, document, archive
    quality_preset = request.data.get('quality_preset', 'standard')
    
    if not uploaded_file:
        return Response({'error': 'No file provided'}, status=400)
    
    # Validate file type based on plan.json supported conversions
    file_ext = uploaded_file.name.split('.')[-1].lower()
    
    # Map file extensions to conversion categories per plan.json
    conversion_categories = {
        'image': ['jpg', 'jpeg', 'png', 'webp', 'heic', 'bmp', 'tiff', 'svg'],
        'audio': ['mp3', 'wav', 'flac', 'aac', 'm4a', 'ogg', 'aiff'],
        'video': ['mp4', 'avi', 'mov', 'wmv', 'webm', 'mkv', 'flv'],
        'document': ['pdf', 'docx', 'xlsx', 'pptx', 'csv', 'txt', 'odt', 'rtf'],
        'archive': ['zip', '7z', 'rar', 'tar', 'gz']
    }
    
    detected_category = None
    for category, extensions in conversion_categories.items():
        if file_ext in extensions:
            detected_category = category
            break
    
    if not detected_category:
        return Response({'error': f'Unsupported file type: {file_ext}'}, status=400)
    
    # Create job ID and prepare for Go microservice processing
    job_id = str(uuid.uuid4())
    
    return Response({
        'job_id': job_id,
        'status': 'pending',
        'original_filename': uploaded_file.name,
        'file_size': uploaded_file.size,
        'source_format': file_ext,
        'target_format': target_format or file_ext,
        'conversion_category': detected_category,
        'quality_preset': quality_preset,
        'user_id': request.user.id,
        'message': f'{detected_category.title()} conversion job created'
    })

@api_view(['GET'])
@permission_classes([IsAuthenticated])
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