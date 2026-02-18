# Welcome to My Backend Architecture 🚀

A Backend service built with **Go (Golang)** and the **Gin framework**. It uses **PostgreSQL** for the database with **dbmate** for migrations, and **sqlc** for generating type-safe SQL queries.

The project is deployed using **Docker** for containerization, **GitHub Actions** for automated CI/CD, and **Cloudflare Tunnel** for secure public access.

## 🏗️ Architecture

![Architecture Diagram](./template/Architecture.png)

**Key Infrastructure:**
- **Local Server:** Windows Subsystem for Linux (WSL2) Ubuntu.
- **Containerization:** Docker & Docker Compose.
- **Database Migrations:** Dbmate.
- **CI/CD Pipeline:** GitHub Actions with a Self-Hosted Runner.
- **Secure Tunneling:** Cloudflare Zero Trust (No open router ports).

---

## 🛠️ Prerequisites

Before you begin, ensure you have the following installed on your machine:
- [Go](https://golang.org/doc/install) (v1.25 or latest)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [Dbmate](https://github.com/amacneil/dbmate) (For database migrations)
- [Swag](https://github.com/swaggo/swag) (For generating Swagger API docs)

---

## ⚙️ Environment Setup

1. Clone the repository:
   ```bash
   git clone git@github.com:yourusername/be-project.git
   cd be-project
   ```

2. Create a .env file in the root directory and configure your variables (refer to .env.example if available):
   ```bash
   # Server Config
   SERVER_PORT=3000
   ENVIRONMENT=development
   
   # Database Config
   DB_USER=postgres
   DB_PASSWORD=yourpassword
   DB_NAME=your_db
   DB_PORT=5432
   
   # For dbmate / Local Development
   DATABASE_URL="postgres://postgres:yourpassword@host:5432/your_db?sslmode=disable"
   
   # Security
   JWT_SECRET=your_jwt_secret_key
   TOKEN_DURATION=24h
   ```
## 🚀 Running the Project
You can run this project in two ways: Local Mode (for active development) or Docker Mode (for testing/production simulation).

### Option A: Running Locally (Development)
Make sure your local PostgreSQL service is running and port 5432 is available.

1. Run Database Migrations:
   ```bash
   dbmate up
   ```

2. Generate/Update Swagger Docs and SQL Query
   ```bash
   make gen
   ```

3. Start the Server:
   ```bash
   go run cmd/api/main.go
   ```

### Option B: Running with Docker (Isolated Environment)
This will spin up both the Go API and PostgreSQL in isolated containers.

1. Start the Containers:

   ```bash
   docker compose up -d --build
   ```

2. Run Migrations inside the Host (targeting Docker DB):

   ```bash
   dbmate up
   ```

3. Check Logs:

   ```bash
   docker compose logs -f api
   ```

4. Stop the Containers:

   ```bash
   docker compose down
   ```

## 📚 API Documentation (Swagger)
This project uses Swagger for API documentation. Once the server is running (either locally or via Docker), you can access the interactive API docs at:

👉 http://localhost:3000/swagger/index.html

## 🔄 CI/CD & Deployment
This project utilizes an automated CI/CD pipeline via GitHub Actions.

- Continuous Integration (CI): Triggered on every push to the main branch to build and test the Go code.
- Continuous Deployment (CD): A Self-Hosted Runner on the WSL server listens for successful builds, automatically pulls the latest code, and reconstructs the Docker containers.

The API is securely exposed to the internet using Cloudflare Tunnel.

*Notes: I'm always open to collaboration and pull requests! Feel free to contribute.* 🚀