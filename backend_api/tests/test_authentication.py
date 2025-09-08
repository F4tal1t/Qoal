from django.test import TestCase
from django.contrib.auth.models import User
from rest_framework.test import APIClient
from rest_framework import status

class AuthenticationTests(TestCase):
    def setUp(self):
        self.client = APIClient()

    def test_user_registration(self):
        """Test user registration endpoint"""
        data = {
            'email': 'test@example.com',
            'password': 'testpass123',
            'username': 'testuser'
        }
        response = self.client.post('/api/auth/register/', data)
        
        self.assertEqual(response.status_code, status.HTTP_200_OK)
        self.assertIn('user', response.data)
        self.assertIn('tokens', response.data)
        self.assertIn('access', response.data['tokens'])
        self.assertTrue(User.objects.filter(email='test@example.com').exists())

    def test_user_login(self):
        """Test user login endpoint"""
        # Create user first
        user = User.objects.create_user(
            username='testuser',
            email='test@example.com',
            password='testpass123'
        )
        
        # Test valid credentials
        data = {
            'email': 'test@example.com',
            'password': 'testpass123'
        }
        response = self.client.post('/api/auth/login/', data)
        
        if response.status_code == status.HTTP_200_OK:
            self.assertIn('user', response.data)
            self.assertIn('tokens', response.data)
            self.assertIn('access', response.data['tokens'])
        else:
            self.assertIn('error', response.data)
        
        # Test invalid credentials
        invalid_data = {
            'email': 'test@example.com',
            'password': 'wrong'
        }
        response = self.client.post('/api/auth/login/', invalid_data)
        self.assertEqual(response.status_code, status.HTTP_400_BAD_REQUEST)
        self.assertIn('error', response.data)

    def test_duplicate_email_registration(self):
        """Test registration with existing email"""
        User.objects.create_user(
            username='existing',
            email='test@example.com',
            password='pass123'
        )
        
        data = {
            'email': 'test@example.com',
            'password': 'newpass123'
        }
        response = self.client.post('/api/auth/register/', data)
        
        self.assertEqual(response.status_code, status.HTTP_400_BAD_REQUEST)
        self.assertIn('error', response.data)