# AI Interview Platform – Backend

This repository contains the **Golang backend API** for the AI Interview Platform.

The backend handles interview logic, AI integration, and communication with the frontend.

---

## 🚀 Features

* REST API built with Golang
* AI-powered interview processing
* Candidate answer evaluation
* Secure API endpoints
* JSON-based communication with frontend

---

## 🛠️ Tech Stack

* Golang
* Gin / Fiber / net/http (depending on your framework)
* REST API
* JSON
* Environment variables (.env)

---

## 📂 Project Structure

```
backend/
 ├── handlers
 ├── routes
 ├── services
 ├── models
 ├── main.go
 └── go.mod
```

---

## ⚙️ Installation

Clone the repository

```
git clone https://github.com/Dharaneeshponnuvel/interview-ai-backend.git
```

Move into the project folder

```
cd interview-ai-backend
```

Install dependencies

```
go mod tidy
```

Run the server

```
go run main.go
```

The backend server will start at

```
http://localhost:8080
```

---

## 📌 Environment Variables

Create a `.env` file

```
PORT=8080
AI_API_KEY=your_api_key
DATABASE_URL=your_database_url
```

---

## 🔗 Frontend Repository

Frontend React application:

https://github.com/Dharaneeshponnuvel/interview-ai-frontend

---

## 👨‍💻 Author

**Dharaneesh P**

GitHub:
https://github.com/Dharaneeshponnuvel
