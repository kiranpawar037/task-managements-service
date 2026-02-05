# task-managements-service
So the flow becomes:
1. Clone repo
2. `go mod tidy`
3. Setup `env.yaml`
4. Run service

📌 Task Management Service (Backend)

A role-based Task Management Backend built with Go (Gin + GORM).
Supports Admin & User roles, JWT authentication, task assignment, and automatic task completion using background workers.

🚀 Features
🔐 Authentication

Signup (creates Pending User)

First Login → moves user from pending_users to users

JWT-based authentication

Role-based access (admin, user)

👤 Admin Capabilities

View all users

Create tasks

Assign tasks to users

View all tasks

Delete tasks

👷 User Capabilities

View assigned tasks

Update task status (pending, in_progress, completed)

⏱ Background Worker

Automatically marks tasks as completed after a configured time


🧱 Tech Stack

Language: Go

Framework: Gin

ORM: GORM

Database: MySQL

Auth: JWT

Config: YAML (env.yaml)

Worker: Goroutines + Channels


📁 Project Structure (High Level)
task-managements-service/
├── cmd/
│   └── main.go
├── pkg/
│   ├── admin/
│   │   └── task/
│   ├── auth/
│   │   ├── signin/
│   │   └── signup/
│   ├── config/
│   ├── database/
│   ├── helper/
│   ├── models/
│   ├── user/
│   │   └── userflow/
│   └── worker/
├── routes/
│   └── useradmin/
├── env.yaml
└── README.md

⚙️ Environment Configuration (env.yaml)
--------------------------------------------------------------------------------------------
Create a file named env.yaml in a directory and export its path using ENV_PATH.

✅ env.yaml
env: development

database:
  host: localhost
  port: "3306"
  username: root
  password: your_db_password
  databaseName: task_management

task:
  autoCompleteMinutes: 30 

------------------------------------------------------  
▶️ Start the Backend Service

This backend supports multiple services.
Currently available:

user-admin

🔥 Run Command
go run cmd/main.go user-admin
--------------------------------------------------------------------------
🔑 Authentication Flow
Signup

User signs up → data stored in pending_users

No JWT returned

First Login

Credentials verified from pending_users

User moved to users table

JWT generated

Next Logins

Direct login from users

JWT generated
-------------------------------------------------------------

🧪 API Overview
Host: http://localhost

Port: 10001

API Version Prefix: /v1

So Base URL:

http://localhost:10001/v1

🩺 Health Check
GET http://localhost:10001/v1/user

🔐 Authentication APIs
1️⃣ User Signup
POST http://localhost:10001/v1/user/sign-up


Body (JSON):

{
  "email": "admin@example.com",
  "password": "Admin@123",
  "role": "admin"
}

2️⃣ User Login
POST http://localhost:10001/v1/user/sign-in


Body (JSON):

{
  "email": "admin@example.com",
  "password": "Admin@123"
}


Response:

{
  "message": "Login successful",
  "token": "JWT_TOKEN_HERE"
}

Authorization in Headers 

All routes below require:
key          :    value         
Authorization:  JWT_TOKEN

👑 Admin APIs
3️⃣ Create Task
POST http://localhost:10001/v1/admin/create-task


Body:

{
  "title": "Prepare Report",
  "description": "Prepare weekly sales report"
}

4️⃣ Get All Tasks
GET http://localhost:10001/v1/admin/get-all-tasks

5️⃣ Get Task By ID
GET http://localhost:10001/v1/admin/get-task/1

6️⃣ Delete Task
DELETE http://localhost:10001/v1/admin/delete-task/1

7️⃣ Assign Task to User
PUT http://localhost:10001/v1/admin/assign-task/1


Body:

{
  "user_id": 2
}

8️⃣ Get All Users
GET http://localhost:10001/v1/admin/get-all-users

👷 User APIs
9️⃣ Get My Tasks
GET http://localhost:10001/v1/user/get-my-tasks

🔟 Update Task Status
PATCH http://localhost:10001/v1/user/update-task-status/1


Body:

{
  "status": "in_progress"
}


Allowed values:

pending | in_progress | completed