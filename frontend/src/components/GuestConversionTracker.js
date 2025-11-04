import React, { useState, useEffect } from 'react';
import axios from 'axios';

// Configure axios base URL
axios.defaults.baseURL = process.env.REACT_APP_API_URL || 'http://localhost:8000';

const GuestConversionTracker = ({ onLimitReached, onConversionUsed }) => {
  const [guestStatus, setGuestStatus] = useState({
    can_convert: true,
    remaining_conversions: 3,
    is_authenticated: false
  });

  useEffect(() => {
    checkGuestStatus();
  }, []);

  const checkGuestStatus = async () => {
    try {
      console.log('Checking guest status...');
      const response = await axios.get('/api/conversions/guest-status/');
      console.log('Guest status response:', response.data);
      setGuestStatus(response.data);
      
      if (!response.data.can_convert) {
        onLimitReached && onLimitReached();
      }
    } catch (error) {
      console.error('Error checking guest status:', error.response?.data || error.message);
    }
  };

  const updateAfterConversion = () => {
    if (!guestStatus.is_authenticated) {
      const newRemaining = Math.max(0, guestStatus.remaining_conversions - 1);
      setGuestStatus(prev => ({
        ...prev,
        remaining_conversions: newRemaining,
        can_convert: newRemaining > 0
      }));
      
      onConversionUsed && onConversionUsed(newRemaining);
      
      if (newRemaining === 0) {
        onLimitReached && onLimitReached();
      }
    }
  };

  return {
    guestStatus,
    checkGuestStatus,
    updateAfterConversion,
    canConvert: guestStatus.can_convert,
    remainingConversions: guestStatus.remaining_conversions,
    isAuthenticated: guestStatus.is_authenticated
  };
};

export default GuestConversionTracker;