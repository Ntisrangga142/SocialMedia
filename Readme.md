# Social Media

A social media application built with Go (Golang) that allows users to connect through posts, likes, comments, and follows, and handle account operations with robust authentication and caching support.

## 📋 Project Overview

**Social Media** is a web application built with **Go (Golang)** that allows users to connect through posts, likes, comments, and follows.  
It provides robust **authentication**, efficient **caching with Redis**, and seamless **account management** to ensure secure and scalable social interactions.

### Technologies Used

- **Go (Golang)** - Main programming language
- **Gin** - Web framework
- **PostgreSQL** - Database (via pgx/v5)
- **Redis** - Caching and session management
- **Docker** - Containerization and deployment

### ✨ Features

- ✅  **User Accounts**: Register, login, logout, and manage profiles with profile pictures.  
- ✅  **Posts**: Create, update, delete posts with text and multiple images.  
- ✅  **Likes & Comments**: Interact with posts by liking and commenting.  
- ✅  **Followers**: Follow and unfollow users, see followers and following lists.  
- ✅  **Notifications**: Real-time notifications for likes, comments, and new followers.  
- ✅  **Caching**: Redis caching for faster response on frequently accessed data.  
- ✅  **Secure Authentication**: JWT-based authentication for secure API access.  


## 🚀 Installation

### Prerequisites

- Go 1.25
- PostgreSQL
- Redis
- Docker

### Environment Variables

Create a `.env` file in the root directory:

```env
# env for golang db config
DB_USER=your_sosmed
DB_PASS=your_sosmed
DB_HOST=pg-db
DB_PORT=5432
DB_NAME=your_sosmed

# env for compose rdb
RDBHOST=rdb
RDBPORT=6380

# JWT golang
JWT_SECRET=a-string-secret-at-least-256-bits-long
JWT_ISSUER=your_issue

# env for compose pg-db
POSTGRES_USER=your_sosmed
POSTGRES_PASSWORD=your_sosmed
POSTGRES_DB=your_sosmed
```

### 🛠 Setup Instructions (Docker Only)

Follow these steps to run the Social Media application locally using Docker without Docker Compose.

---

#### 1. Clone the repository

```bash
git clone https://github.com/username/social-media-app.git
cd social-media-app
```


#### 2. Create a .env file

```bash
cp .env.example .env
```

Example .env:
```bash
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=social_media
REDIS_HOST=redis
REDIS_PORT=6379
JWT_SECRET=mysecretkey
```

#### 3. Run PostgreSQL container

```bash
docker run -d \
  --name social_media_db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=social_media \
  -p 5432:5432 \
  postgres:15
```

#### 4. Run Redis container

```bash
docker run -d \
  --name social_media_redis \
  -p 6379:6379 \
  redis:7
```

#### 5. Build the backend Docker image

```bash
docker build -t social-media-backend .
```

#### 6. Run the backend container

```bash
docker run -d \
  --name social-media-backend \
  --env-file .env \
  -p 8080:8080 \
  --link social_media_db:db \
  --link social_media_redis:redis \
  social-media-backend
```

#### 7. Run database migrations (if needed)

```bash
docker exec -it social-media-backend go run migrations/main.go
```

## 📚 API Documentation

### Authentication Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST   | `/auth`          | User login        | ❌ |
| POST   | `/auth/register` | User registration | ❌ |
| DELETE | `/auth`          | User logout       | ✅ |

### User Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET    | `/user`           | Get user profile | ✅ |
| PATCH  | `/user`           | Update profile   | ✅ |
| POST   | `/user/:id`       | Follow           | ✅ |
| DELETE | `/user/:id`       | Unfollow         | ✅ |
| GET    | `/user/follower`  | Get Followers    | ✅ |
| GET    | `/user/following` | Get Following    | ✅ |

### Post Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET    | `/post`             | Get Following Posts        | ✅ |
| GET    | `/post/:id`         | Get Post Detail            | ✅ |
| POST   | `/post`             | Create Post                | ✅ |
| POST   | `/post/:id/like`    | Like Post                  | ✅ |
| DELETE | `/post/:id/like`    | Unlike Post                | ✅ |
| POST   | `/post/comment`     | Create Comment             | ✅ |
| GET    | `/post/:id/comment` | Get All Comments By Post   | ✅ |

### Notification Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET    | `/notif` | Get Unread Notifications | ✅ |


### Static Files

Profile images are served from `/profile/*` directory.
Post images are served from `/post/*` directory.

## 🔐 Authentication

Most endpoints require authentication using JWT tokens. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## 📝 Version History

### Version 1.0.0 (Current)
- Initial release.
- Basic authentication, user profile, posts, likes, comments, and followers features implemented.
- Added post detail, post images upload, caching with Redis, and notification system.
- Implemented follower/following feed and API improvements. 
- Added Swagger documentation and improved Docker setup.

## Verification

Your application should now be running and accessible at:
- Backend: http://localhost:8080
- PostgreSQL: localhost:5432
- Redis: localhost:6380