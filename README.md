# 💸 Xpense

**Xpense** is a full-stack web application that allows users to manage and visualize their expenses.  
The backend is a REST API built with the **Gin** framework, while the frontend is developed using **React** for smooth interaction and data visualization.

---

## 🛠️ Backend Requirements

**Tech Stack:**  
- Framework: Gin (Go)  
- Database: PostgreSQL

### 🔧 REST API Features

- ✅ Create a new expense
- 📄 Read the list of expenses
- 🔍 Read a single expense in detail (not implemented to FE)
- ✏️ Update an existing expense
- ❌ Delete an expense

### 🧾 Expense Object Schema

| Field      | Type     | Description                         |
|------------|----------|-------------------------------------|
| `id`       | UUID     | Unique identifier                   |
| `title`    | string   | Description of the expense          |
| `amount`   | number   | Expense amount                      |
| `category` | string   | e.g., food, travel, utilities, etc. |
| `date`     | ISO date | Date of the expense                 |

### 📡 API Guidelines

- RESTful architecture
- JSON responses
- Proper HTTP status codes (`200`, `201`, `400`, `404`, etc.)

---

## 🌐 Frontend Requirements

**Tech Stack:**  
- Framework: React  
- Charting Library: Recharts

### 🖥️ Features

- 📋 Display a list of expenses
- ➕ Add new expenses
- 📝 Edit and delete existing expenses
- 🔍 Filter expenses by category
- 📆 Summary pages for:
  - Monthly expenses
  - Weekly expenses
- 📊 Visualize expenses using:
  - Pie Chart or Bar Chart

---

## ⚙️ Bonus Technical Add-ons

- 🔐 Authentication (JWT or OAuth)
- ⛑️ TypeScript support for React

---

*Feel free to build and scale as you go—Xpense is designed to be as lightweight or feature-rich as you need.*
