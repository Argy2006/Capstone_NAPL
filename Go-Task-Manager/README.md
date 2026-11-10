# TaskFlow – Personal Task Manager

A clean, modular personal task management application built in Go (Golang) as part of B.Tech Semester V Capstone Project (SEC042 – New Age Programming Languages Lab).

---

## Overview

TaskFlow is designed to help users organize, track, and manage their daily tasks efficiently. The project demonstrates core principles of modern Go development:
- Clean modular design separating models, business services, storage, and HTTP handlers.
- Pure Go standard library networking (`net/http`) and template rendering (`html/template`).
- File-based persistent storage using structured JSON with thread safety.
- An organic, distraction-free visual design language featuring warm ivory tones, soft cream surfaces, sage green accents, and terracotta highlights.

---

## Features (Planned & Pending Implementation)

- [x] **Modular Go Foundation:** Scalable layered architecture.
- [x] **Google Stitch UI Integration:** Successfully imported 3 Stitch screens (Dashboard, All Tasks Workspace, and Empty State).
- [ ] **Task Creation:** Add new tasks with title, description, priority, category, and due date.
- [ ] **Task Tracking & Updates:** Mark tasks as pending, in-progress, or completed.
- [ ] **Task Deletion:** Remove completed or obsolete tasks.
- [ ] **Task Filtering & Categorization:** Organize tasks by status, priority, or category.
- [ ] **Persistent Storage:** Safe local persistence in `data/tasks.json`.

*(Full backend CRUD logic and complete UI interactivity are pending Phase II implementation confirmation.)*

---

## Project Structure

```
Go-Task-Manager/
├── main.go                     # Entry point & HTTP server
├── models/                     # Data models (Task, Priority, Status)
├── handlers/                   # HTTP handlers & API endpoints
├── services/                   # Business logic layer
├── storage/                    # JSON file storage manager
├── templates/                  # HTML templates (index.html)
├── static/                     # Static files (CSS, JS, Assets)
│   ├── css/
│   ├── js/
│   └── assets/
├── data/                       # Persistent JSON storage (tasks.json)
├── tests/                      # Automated test suite
├── docs/                       # Project documentation
├── screenshots/                # Visual verification captures
├── go.mod                      # Go module definition
├── .gitignore                  # Git ignore definitions
└── LICENSE                     # MIT License
```

---

## Installation & Setup

### Prerequisites
- [Go](https://golang.org/dl/) version 1.20 or newer (tested on Go 1.26+)

### Running the Application

1. Open your terminal and navigate to the project directory:
   ```bash
   cd Go-Task-Manager
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

3. Open your browser and navigate to:
   - Dashboard: `http://localhost:8081/`
   - All Tasks Workspace: `http://localhost:8081/tasks`
   - Empty State Onboarding: `http://localhost:8081/empty`

### Running Tests

To run the automated tests:
```bash
cd Go-Task-Manager
go test -v ./...
```

---

## Status

- **Foundation & Structure:** Completed
- **Stitch UI Import:** Completed (Dashboard, All Tasks Workspace, and Empty State integrated)
- **Backend CRUD Logic:** Ready for implementation pending Phase II confirmation

