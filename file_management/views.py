from rest_framework.decorators import api_view, parser_classes
from rest_framework.parsers import MultiPartParser
from rest_framework.response import Response
from conversions.models import ConversionJob, SupportedFormat
import magic
import uuid

@api_view(['POST'])
@parser_classes([MultiPartParser])
def upload_file(request):
    uploaded_file = request.FILES.get('file')
    if not uploaded_file:
        return Response({'error': 'No file provided'}, status=400)
    
    # Validate file type
    mime_type = magic.from_buffer(uploaded_file.read(1024), mime=True)
    uploaded_file.seek(0)
    
    # Check if format is supported
    try:
        source_format = SupportedFormat.objects.get(mime_type=mime_type)
    except SupportedFormat.DoesNotExist:
        return Response({'error': f'Unsupported file type: {mime_type}'}, status=400)
    
    # Create conversion job
    job = ConversionJob.objects.create(
        user=request.user,
        original_filename=uploaded_file.name,
        file_size=uploaded_file.size,
        source_format=source_format,
        target_format=source_format,  # Default to same format
        input_file_path=f"uploads/{uuid.uuid4()}/{uploaded_file.name}"
    )
    
    return Response({
        'job_id': str(job.job_id),
        'original_filename': job.original_filename,
        'source_format': source_format.name,
        'supported_targets': list(SupportedFormat.objects.filter(
            category=source_format.category
        ).values_list('name', flat=True))
    })