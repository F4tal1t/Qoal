from rest_framework import viewsets
from .models import ConversionJob
from .serializers import ConversionJobSerializer

class ConversionViewSet(viewsets.ModelViewSet):
    queryset = ConversionJob.objects.all()
    serializer_class = ConversionJobSerializer