from django.shortcuts import render
from rest_framework import status
from rest_framework.decorators import api_view, permission_classes
from rest_framework.permissions import AllowAny
from rest_framework.response import Response
from rest_framework_simplejwt.tokens import RefreshToken
from django.contrib.auth import authenticate
from django.contrib.auth.models import User

@api_view(['POST'])
@permission_classes([AllowAny])
def register_user(request):
    email = request.data.get('email')
    password = request.data.get('password')
    username = request.data.get('username', email)
    
    if not email or not password:
        return Response({'error': 'Email and password are required'}, status=400)
    
    if User.objects.filter(email=email).exists():
        return Response({'error': 'Email already exists'}, status=400)
    
    if User.objects.filter(username=username).exists():
        return Response({'error': 'Username already exists'}, status=400)
    
    user = User.objects.create_user(username=username, email=email, password=password)
    print(f"Created user: {user.username}, email: {user.email}")
    
    refresh = RefreshToken.for_user(user)
    
    return Response({
        'user': {'id': user.id, 'email': user.email, 'username': user.username},
        'tokens': {
            'refresh': str(refresh),
            'access': str(refresh.access_token),
        }
    })

@api_view(['POST'])
@permission_classes([AllowAny])
def login_user(request):
    email = request.data.get('email')
    password = request.data.get('password')
    
    if not email or not password:
        return Response({'error': 'Email and password are required'}, status=400)
    
    # Try to find user by email first
    try:
        user_obj = User.objects.get(email=email)
        print(f"Found user: {user_obj.username}, email: {user_obj.email}")
        
        # Check password manually first
        if user_obj.check_password(password):
            user = user_obj
            print("Password check passed")
        else:
            user = None
            print("Password check failed")
            
    except User.DoesNotExist:
        print(f"No user found with email: {email}")
        return Response({'error': 'User not found'}, status=400)
    
    if not user:
        print("Authentication failed")
        return Response({'error': 'Invalid password'}, status=400)
    
    if not user.is_active:
        return Response({'error': 'Account is disabled'}, status=400)
    
    refresh = RefreshToken.for_user(user)
    return Response({
        'user': {'id': user.id, 'email': user.email, 'username': user.username},
        'tokens': {
            'refresh': str(refresh),
            'access': str(refresh.access_token),
        }
    })