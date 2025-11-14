# Fly.io Deployment Guide

## Prerequisites

1. **Fly.io Account**: Sign up at https://fly.io
2. **GitHub Repository**: Push code to GitHub
3. **AWS S3 Bucket**: For file storage (qoala bucket)
4. **Fly.io Redis**: Free tier included with Fly.io

---

## Step 1: Install Fly CLI (Local Testing)

```bash
# Windows (PowerShell)
iwr https://fly.io/install.ps1 -useb | iex

# Mac/Linux
curl -L https://fly.io/install.sh | sh

# Login
flyctl auth login
```

---

## Step 2: Create Fly.io App

```bash
cd backend

# Create app (don't deploy yet)
flyctl launch --no-deploy --name qoal-converter --region iad

# This creates fly.toml (already exists in repo)
```

---

## Step 3: Create Postgres Database (FREE)

```bash
# Create free shared Postgres
flyctl postgres create --name qoal-db --region iad --vm-size shared-cpu-1x --volume-size 1

# Attach to app
flyctl postgres attach qoal-db -a qoal-converter
```

This automatically sets `DATABASE_URL` secret.

---

## Step 4: Set Environment Secrets

```bash
# JWT Secret
flyctl secrets set JWT_SECRET=$(openssl rand -base64 32) -a qoal-converter

# Redis (Fly.io free tier)
flyctl redis create --name qoal-redis -a qoal-converter
# This automatically sets REDIS_URL secret

# AWS S3 Credentials
flyctl secrets set AWS_REGION="us-east-1" -a qoal-converter
flyctl secrets set AWS_ACCESS_KEY_ID="your-access-key" -a qoal-converter
flyctl secrets set AWS_SECRET_ACCESS_KEY="your-secret-key" -a qoal-converter
flyctl secrets set AWS_S3_BUCKET="qoala" -a qoal-converter

# Verify secrets
flyctl secrets list -a qoal-converter
```

---

## Step 5: Create Volume for Temp Files

```bash
# Create 3GB volume (FREE)
flyctl volumes create qoal_storage --size 3 --region iad -a qoal-converter
```

---

## Step 6: Deploy Manually (First Time)

```bash
# Deploy from backend directory
flyctl deploy --config fly.toml

# Check status
flyctl status -a qoal-converter

# View logs
flyctl logs -a qoal-converter
```

---

## Step 7: Setup GitHub Actions (Auto Deploy)

### Add Fly.io Token to GitHub Secrets

1. Get Fly.io API token:
```bash
flyctl auth token
```

2. Go to GitHub repo → Settings → Secrets and variables → Actions
3. Click "New repository secret"
4. Name: `FLY_API_TOKEN`
5. Value: Paste the token from step 1
6. Click "Add secret"

### Workflow is Already Configured

The `.github/workflows/deploy.yml` file is already in the repo. It will:
- Trigger on push to `main` branch
- Deploy backend to Fly.io automatically

---

## Step 8: Test Deployment

```bash
# Get app URL
flyctl info -a qoal-converter

# Test health endpoint
curl https://qoal-converter.fly.dev/health

# Test API
curl https://qoal-converter.fly.dev/api/auth/register \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'
```

---

## Frontend Deployment (Vercel/Netlify)

### Update Frontend API URL

In `frontend/src/services/api.ts`:
```typescript
const API_BASE = 'https://qoal-converter.fly.dev/api';
```

### Deploy to Vercel

```bash
cd frontend
npm install -g vercel
vercel login
vercel --prod
```

### Or Deploy to Netlify

```bash
cd frontend
npm run build
# Upload dist/ folder to Netlify
```

---

## Monitoring & Maintenance

### View Logs
```bash
flyctl logs -a qoal-converter
```

### Check Status
```bash
flyctl status -a qoal-converter
```

### Scale (if needed)
```bash
# Increase RAM (costs money)
flyctl scale memory 512 -a qoal-converter

# Add more instances
flyctl scale count 2 -a qoal-converter
```

### Database Access
```bash
# Connect to Postgres
flyctl postgres connect -a qoal-db

# Run migrations
flyctl ssh console -a qoal-converter
./qoal-backend migrate
```

---

## Cost Breakdown (FREE TIER)

| Service | Plan | Cost |
|---------|------|------|
| Fly.io VM (256MB) | Free | $0 |
| Fly.io Volume (3GB) | Free | $0 |
| Fly.io Postgres (1GB) | Free | $0 |
| Fly.io Redis | Free | $0 |
| AWS S3 (with lifecycle) | Pay-as-you-go | ~$0.50/month |
| Vercel/Netlify Frontend | Free | $0 |
| **Total** | | **~$0.50/month** |

---

## Troubleshooting

### App won't start
```bash
flyctl logs -a qoal-converter
# Check for missing secrets or database connection issues
```

### Database connection failed
```bash
# Verify DATABASE_URL is set
flyctl secrets list -a qoal-converter

# Reconnect database
flyctl postgres attach qoal-db -a qoal-converter
```

### Out of memory
```bash
# Check memory usage
flyctl status -a qoal-converter

# Upgrade to 512MB (costs ~$5/month)
flyctl scale memory 512 -a qoal-converter
```

### Storage full
```bash
# Check volume usage
flyctl ssh console -a qoal-converter
df -h

# Clean up old files manually
rm -rf /app/temp/*
```

---

## GitHub Actions Deployment

Once `FLY_API_TOKEN` is set in GitHub secrets:

1. Push to main branch:
```bash
git add .
git commit -m "Deploy to Fly.io"
git push origin main
```

2. GitHub Actions will automatically:
   - Build Docker image
   - Deploy to Fly.io
   - Run health checks

3. Check deployment status:
   - Go to GitHub repo → Actions tab
   - View workflow run logs

---

## Success Checklist

- [ ] Fly.io app created
- [ ] Postgres database attached
- [ ] All secrets set (JWT, Redis, AWS)
- [ ] Volume created and mounted
- [ ] Manual deployment successful
- [ ] GitHub Actions token added
- [ ] Auto-deployment working
- [ ] Frontend updated with production API URL
- [ ] S3 lifecycle policy configured (24 hours)
- [ ] Health endpoint responding
- [ ] Test file conversion working

---

## Production URL

Backend: `https://qoal-converter.fly.dev`
Frontend: `https://your-app.vercel.app` (or Netlify)

---

## Support

- Fly.io Docs: https://fly.io/docs
- Fly.io Community: https://community.fly.io
- GitHub Issues: Create issue in repo
