import React, { useState } from 'react';

const Profile: React.FC = () => {
  const [user, setUser] = useState({
    username: 'john_doe',
    email: 'john@example.com',
    created_at: '2024-01-01T12:00:00Z',
    total_conversions: 42
  });
  const [isEditing, setIsEditing] = useState(false);
  const [editedUser, setEditedUser] = useState(user);

  const handleSave = () => {
    setUser(editedUser);
    setIsEditing(false);
    // Handle save to API
    console.log('Saving user data:', editedUser);
  };

  const handleCancel = () => {
    setEditedUser(user);
    setIsEditing(false);
  };

  return (
    <div className="flex-center min-h-screen">
      <div className="profile-container">
        <h1>Profile</h1>
        
        <div className="profile-info">
          <div className="info-group">
            <label>Username:</label>
            {isEditing ? (
              <input
                type="text"
                value={editedUser.username}
                onChange={(e) => setEditedUser({...editedUser, username: e.target.value})}
              />
            ) : (
              <span>{user.username}</span>
            )}
          </div>

          <div className="info-group">
            <label>Email:</label>
            {isEditing ? (
              <input
                type="email"
                value={editedUser.email}
                onChange={(e) => setEditedUser({...editedUser, email: e.target.value})}
              />
            ) : (
              <span>{user.email}</span>
            )}
          </div>

          <div className="info-group">
            <label>Member Since:</label>
            <span>{new Date(user.created_at).toLocaleDateString()}</span>
          </div>

          <div className="info-group">
            <label>Total Conversions:</label>
            <span>{user.total_conversions}</span>
          </div>
        </div>

        <div className="profile-actions">
          {isEditing ? (
            <>
              <button onClick={handleSave} className="save-btn">Save</button>
              <button onClick={handleCancel} className="cancel-btn">Cancel</button>
            </>
          ) : (
            <button onClick={() => setIsEditing(true)} className="edit-btn">Edit Profile</button>
          )}
        </div>

        <div className="profile-links">
          <a href="/">Back to Home</a>
          <button className="logout-btn">Logout</button>
        </div>
      </div>
    </div>
  );
};

export default Profile;