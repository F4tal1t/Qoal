<img width="1584" height="396" alt="Dibby (2)" src="https://github.com/user-attachments/assets/f7194852-be7e-444e-a36e-e8b8fc1b84c7" />
<br></br>

### A full-stack web application for converting files between multiple formats including images, videos, audio, documents, and archives


## Live Demo

- Frontend: [https://qoal.onrender.com](https://qoal.onrender.com)
- Backend API: [https://qoalbackend.onrender.com](https://qoalbackend.onrender.com)

## Features

- User authentication with JWT tokens
- Support for 20+ file formats across 5 categories
- Real-time conversion status tracking
- Cloud storage integration with AWS S3
- Responsive design with modern UI
- Secure file handling and processing

## Tech Stack

### Frontend
- React 18 with TypeScript
- React Router for navigation
- Vite for build tooling
- TailwindCSS for styling
- Custom animation components

### Backend
- Go (Golang) with Gin framework
- PostgreSQL database
- Redis for job queue management
- AWS S3 for file storage
- JWT authentication
- GORM for database operations

### Infrastructure
- Docker containerization
- Render for backend hosting
- Netlify for frontend hosting
- GitHub Actions for CI/CD

## Supported File Formats

### Images
PNG, JPG, JPEG, WEBP, GIF, BMP, TIFF, SVG, ICO, HEIC

### Videos
MP4, AVI, MOV, MKV, WEBM, FLV, WMV, M4V

### Audio
MP3, WAV, AAC, FLAC, OGG, M4A, WMA, AIFF

### Documents
PDF, DOCX, DOC, TXT, RTF, ODT, HTML, MD

### Archives
ZIP, RAR, TAR, GZ, 7Z, BZ2


## API Endpoints

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /api/auth/profile` - Get user profile

### File Operations
- `POST /api/upload` - Upload and queue file for conversion
- `GET /api/jobs/:id` - Get conversion job status
- `GET /api/download/:id` - Download converted file
- `GET /api/jobs` - List user's conversion jobs

## Local Development

### Prerequisites
- Go 1.24+
- Node.js 18+
- PostgreSQL 14+
- Redis 7+
- AWS S3 account

### Backend Setup
```bash
cd backend
cp .env.template .env
# Configure environment variables
go mod download
go run main.go
```

### Frontend Setup
```bash
cd frontend
npm install
npm run dev
```

## Environment Variables

### Backend
```
DATABASE_URL=postgresql://user:password@localhost:5432/qoal
REDIS_URL=localhost:6379
JWT_SECRET=your-secret-key
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_S3_BUCKET=your-bucket-name
PORT=8000
```

### Frontend
```
VITE_API_URL=http://localhost:8000/api
```

## Deployment

The application uses automated deployment:
- Backend deploys to Render via Docker
- Frontend deploys to Netlify via GitHub integration
- Database hosted on Render PostgreSQL
- Redis hosted on Render Redis

## Security Features

- Password hashing with bcrypt
- JWT token-based authentication
- CORS protection
- File type validation
- Size limit enforcement (30MB)
- Secure file storage with S3

## Performance

- Asynchronous job processing with Redis queue
- Efficient file streaming
- Database connection pooling
- Optimized build with Vite

## Project Structure

```
qoal/
├── backend/
│   ├── handlers/       # HTTP request handlers
│   ├── services/       # Business logic
│   ├── models/         # Data models
│   ├── middleware/     # Auth middleware
│   ├── storage/        # S3 and local storage
│   ├── worker/         # Background job processor
│   └── main.go         # Application entry point
├── frontend/
│   ├── src/
│   │   ├── components/ # React components
│   │   ├── pages/      # Page components
│   │   ├── services/   # API client
│   │   └── App.tsx     # Main app component
│   └── public/         # Static assets
└── docker-compose.yml  # Local development setup
```

## License

MIT License

## Contact

For questions or feedback, please open an issue on GitHub.
