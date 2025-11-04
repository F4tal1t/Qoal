from django.urls import path
from . import views

urlpatterns = [
    path('job/<str:job_id>/status/', views.detailed_job_status, name='detailed_job_status'),
    path('job/<str:job_id>/progress/', views.update_job_progress, name='update_job_progress'),
    path('presets/', views.conversion_presets, name='conversion_presets'),
    path('supported/', views.supported_conversions, name='supported_conversions'),
    path('upload/', views.upload_file_for_conversion, name='upload_file'),
    path('guest-status/', views.guest_conversion_status, name='guest_status'),
    path('test/', views.test_connection, name='test_connection'),
    path('download/<str:job_id>/', views.download_converted_file, name='download_file'),
]