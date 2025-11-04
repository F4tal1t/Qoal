from rest_framework.decorators import api_view, permission_classes
from rest_framework.permissions import IsAuthenticated, AllowAny
from rest_framework.response import Response
from rest_framework import status
from django.core.cache import cache
from django.utils import timezone
import json
import uuid
from datetime import datetime, timedelta
import time

# Guest Conversion System
def get_client_ip(request):
    """Get client IP address"""
    x_forwarded_for = request.META.get('HTTP_X_FORWARDED_FOR')
    if x_forwarded_for:
        ip = x_forwarded_for.split(',')[0]
    else:
        ip = request.META.get('REMOTE_ADDR')
    return ip

def check_guest_conversion_limit(request):
    """Check if guest user can convert (3 conversions per IP per day)"""
    if request.user.is_authenticated:
        return True, 0  # Authenticated users have no limit
    
    client_ip = get_client_ip(request)
    cache_key = f"guest_conversions:{client_ip}"
    
    # Get current conversion count
    conversion_data = cache.get(cache_key, {'count': 0, 'date': timezone.now().date().isoformat()})
    
    # Reset count if new day
    today = timezone.now().date().isoformat()
    if conversion_data['date'] != today:
        conversion_data = {'count': 0, 'date': today}
    
    remaining = max(0, 3 - conversion_data['count'])
    can_convert = remaining > 0
    
    return can_convert, remaining

def increment_guest_conversion(request):
    """Increment guest conversion count"""
    if request.user.is_authenticated:
        return
    
    client_ip = get_client_ip(request)
    cache_key = f"guest_conversions:{client_ip}"
    
    conversion_data = cache.get(cache_key, {'count': 0, 'date': timezone.now().date().isoformat()})
    today = timezone.now().date().isoformat()
    
    if conversion_data['date'] != today:
        conversion_data = {'count': 1, 'date': today}
    else:
        conversion_data['count'] += 1
    
    # Cache for 24 hours
    cache.set(cache_key, conversion_data, timeout=86400)

@api_view(['POST'])
@permission_classes([AllowAny])  # Allow guest conversions
def upload_file_for_conversion(request):
    """Upload file for conversion - supports guest users"""
    try:
        # Check guest conversion limit
        can_convert, remaining = check_guest_conversion_limit(request)
        
        if not can_convert:
            return Response({
                'error': 'Guest conversion limit reached',
                'message': 'You have used your 3 free conversions today. Please register for unlimited conversions.',
                'remaining_conversions': 0,
                'requires_registration': True
            }, status=429)
        
        # Process file upload
        uploaded_file = request.FILES.get('file')
        target_format = request.data.get('target_format')
        
        if not uploaded_file or not target_format:
            return Response({'error': 'File and target_format required'}, status=400)
        
        # Generate job ID
        job_id = str(uuid.uuid4())
        
        # Save file locally (simulate S3 upload)
        import os
        from django.conf import settings
        
        # Create uploads directory if it doesn't exist
        upload_dir = os.path.join(settings.BASE_DIR, 'uploads')
        os.makedirs(upload_dir, exist_ok=True)
        
        # Save uploaded file
        file_path = os.path.join(upload_dir, f"{job_id}_{uploaded_file.name}")
        with open(file_path, 'wb+') as destination:
            for chunk in uploaded_file.chunks():
                destination.write(chunk)
        
        # Store job info in cache
        job_info = {
            'job_id': job_id,
            'user_id': request.user.id if request.user.is_authenticated else 'guest',
            'original_filename': uploaded_file.name,
            'source_format': uploaded_file.name.split('.')[-1].lower(),
            'target_format': target_format.lower(),
            'file_size': uploaded_file.size,
            'input_file_path': file_path,
            'status': 'queued',
            'created_at': time.time()
        }
        
        cache.set(f'job_info:{job_id}', job_info, timeout=3600)
        
        # Increment guest conversion count
        increment_guest_conversion(request)
        
        # Get updated remaining count
        _, new_remaining = check_guest_conversion_limit(request)
        
        # Return job data
        job_data = {
            'job_id': job_id,
            'status': 'queued',
            'original_filename': uploaded_file.name,
            'source_format': uploaded_file.name.split('.')[-1].upper(),
            'target_format': target_format.upper(),
            'file_size': uploaded_file.size,
            'is_guest_conversion': not request.user.is_authenticated,
            'remaining_guest_conversions': new_remaining if not request.user.is_authenticated else None
        }
        
        return Response(job_data, status=201)
        
    except Exception as e:
        return Response({'error': str(e)}, status=500)

@api_view(['GET'])
@permission_classes([AllowAny])
def guest_conversion_status(request):
    """Check guest conversion status"""
    can_convert, remaining = check_guest_conversion_limit(request)
    
    return Response({
        'can_convert': can_convert,
        'remaining_conversions': remaining,
        'is_authenticated': request.user.is_authenticated,
        'daily_limit': 3 if not request.user.is_authenticated else 'unlimited'
    })

@api_view(['GET'])
@permission_classes([AllowAny])  # Allow guest job status checking
def detailed_job_status(request, job_id):
    """Enhanced job status with progress simulation and file processing"""
    try:
        # Get job info
        job_info = cache.get(f'job_info:{job_id}')
        if not job_info:
            return Response({'error': 'Job not found'}, status=404)
        
        # Get or create job progress
        progress_key = f"job_progress:{job_id}"
        progress_data = cache.get(progress_key)
        
        if not progress_data:
            # Initialize new job
            progress_data = {
                "progress": 10,
                "stage": "processing",
                "status": "processing",
                "created_at": time.time()
            }
            cache.set(progress_key, progress_data, timeout=3600)
        
        # Simulate progress over time
        elapsed = time.time() - progress_data.get('created_at', time.time())
        
        if elapsed < 3:  # First 3 seconds
            progress_data['progress'] = min(30, 10 + (elapsed * 7))
            progress_data['stage'] = 'uploading'
            progress_data['status'] = 'processing'
        elif elapsed < 8:  # Next 5 seconds
            progress_data['progress'] = min(80, 30 + ((elapsed - 3) * 10))
            progress_data['stage'] = 'converting'
            progress_data['status'] = 'processing'
        elif elapsed < 10:  # Final 2 seconds
            progress_data['progress'] = min(95, 80 + ((elapsed - 8) * 7.5))
            progress_data['stage'] = 'finalizing'
            progress_data['status'] = 'processing'
        else:  # After 10 seconds, complete
            progress_data['progress'] = 100
            progress_data['stage'] = 'completed'
            progress_data['status'] = 'completed'
            
            # Simulate file conversion completion
            if not progress_data.get('output_created'):
                progress_data['output_created'] = True
                # Create a simple "converted" file
                import os
                from django.conf import settings
                
                output_dir = os.path.join(settings.BASE_DIR, 'outputs')
                os.makedirs(output_dir, exist_ok=True)
                
                # Create converted file (just copy for simulation)
                input_path = job_info['input_file_path']
                output_filename = f"converted_{job_info['original_filename']}"
                output_path = os.path.join(output_dir, f"{job_id}_{output_filename}")
                
                if os.path.exists(input_path):
                    import shutil
                    shutil.copy2(input_path, output_path)
                    progress_data['output_path'] = output_path
        
        cache.set(progress_key, progress_data, timeout=3600)
        
        job_data = {
            'job_id': job_id,
            'status': progress_data['status'],
            'progress': int(progress_data['progress']),
            'current_stage': progress_data['stage'],
            'file_info': {
                'original_filename': job_info['original_filename'],
                'file_size': job_info['file_size'],
                'source_format': job_info['source_format'],
                'target_format': job_info['target_format'],
            },
            'created_at': datetime.fromtimestamp(job_info['created_at']).isoformat(),
        }
        
        if job_data['status'] == 'completed':
            job_data['download_url'] = f"/api/conversions/download/{job_id}/"
            job_data['output_file_size'] = job_info['file_size'] * 0.9  # Simulate compression
        
        return Response(job_data)
        
    except Exception as e:
        return Response({'error': f'Job status error: {str(e)}'}, status=500)

@api_view(['POST'])
@permission_classes([AllowAny])
def update_job_progress(request, job_id):
    """Update job progress (called by Go microservice)"""
    try:
        progress = request.data.get('progress', 0)
        stage = request.data.get('stage', 'processing')
        status_val = request.data.get('status', 'processing')
        
        progress_data = {
            'progress': progress,
            'stage': stage,
            'status': status_val,
            'updated_at': datetime.now().isoformat()
        }
        
        # Store in cache
        cache.set(f"job_progress:{job_id}", progress_data, timeout=3600)
        
        return Response({'message': 'Progress updated successfully'})
        
    except Exception as e:
        return Response({'error': str(e)}, status=400)

@api_view(['GET'])
@permission_classes([AllowAny])
def conversion_presets(request):
    """Get available conversion presets"""
    presets = {
        'video': {
            '4k': {'width': 3840, 'height': 2160, 'bitrate': '8000k'},
            '1080p': {'width': 1920, 'height': 1080, 'bitrate': '2000k'},
            '720p': {'width': 1280, 'height': 720, 'bitrate': '1000k'},
            '480p': {'width': 854, 'height': 480, 'bitrate': '500k'}
        },
        'audio': {
            'high': {'bitrate': '320kbps'},
            'standard': {'bitrate': '192kbps'},
            'compressed': {'bitrate': '128kbps'}
        },
        'image': {
            'best_quality': {'quality': 95},
            'balanced': {'quality': 80},
            'best_compression': {'quality': 60}
        }
    }
    return Response(presets)

@api_view(['GET'])
@permission_classes([AllowAny])
def supported_conversions(request):
    """Get supported file conversions"""
    conversions = {
        'images': ['JPEG ↔ PNG', 'PNG → WebP', 'HEIC → JPEG'],
        'videos': ['MP4 ↔ AVI', 'MOV → MP4', 'WMV → MP4'],
        'audio': ['MP3 ↔ WAV', 'FLAC → MP3', 'AAC → MP3']
    }
    return Response(conversions)

@api_view(['GET'])
@permission_classes([AllowAny])
def test_connection(request):
    """Test API connectivity"""
    return Response({'message': 'API is working!', 'timestamp': datetime.now().isoformat()})

@api_view(['GET'])
@permission_classes([AllowAny])
def download_converted_file(request, job_id):
    """Download converted file"""
    try:
        from django.http import FileResponse, Http404
        import os
        from django.conf import settings
        
        # Get job progress to find output file
        progress_data = cache.get(f"job_progress:{job_id}")
        if not progress_data or progress_data.get('status') != 'completed':
            return Response({'error': 'File not ready for download'}, status=400)
        
        output_path = progress_data.get('output_path')
        if not output_path or not os.path.exists(output_path):
            return Response({'error': 'Converted file not found'}, status=404)
        
        # Get job info for filename
        job_info = cache.get(f'job_info:{job_id}')
        filename = f"converted_{job_info['original_filename']}" if job_info else 'converted_file'
        
        response = FileResponse(
            open(output_path, 'rb'),
            as_attachment=True,
            filename=filename
        )
        return response
        
    except Exception as e:
        return Response({'error': str(e)}, status=500)