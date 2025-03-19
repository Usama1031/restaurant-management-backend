# Restaurant Management Backend System

## Introduction
The **Restaurant Management Backend System** is a RESTful API built using **Golang** and **Gin framework** with **MongoDB** as the database. This backend system manages restaurant operations including **users, orders, menus, tables, and invoices**.

## Features
- **User Authentication** (Admin & Staff login)
- **Menu Management** (Add, update, delete food items)
- **Order Processing** (Create, update, track orders)
- **Table Management** (Assign and manage tables)
- **Invoice Generation** (Track payments and generate invoices)

## Technologies Used
- **Golang** (Gin Framework)
- **MongoDB** (NoSQL Database)
- **JWT Authentication**
- **Postman** (For API testing)

## Installation & Setup
### Prerequisites
- **Go 1.18+**
- **MongoDB** (Local or Cloud)

### Clone the Repository
```sh
$ git clone https://github.com/usama1031/restaurant-management-backend.git
$ cd restaurant-management-backend
```

### Install Dependencies
```sh
$ go mod tidy
```

### Configure Environment Variables
Create a `.env` file and configure database settings:
```
MONGO_URI=mongodb://localhost:27017
DB_NAME=restaurant_db
PORT=8080
JWT_SECRET=your_secret_key
```

### Run the Server
```sh
$ go run main.go
```
Server will start at `http://localhost:8080`

## API Endpoints

### **Authentication**
| Method | Endpoint         | Description              |
|--------|-----------------|--------------------------|
| POST   | /login          | User login               |
| POST   | /register       | User registration        |

### **Menu Management**
| Method | Endpoint         | Description              |
|--------|-----------------|--------------------------|
| GET    | /menus           | Get all menus       |
| GET    | /menus           | Get a specifc menu  |
| POST   | /menus           | Add a new menu      |
| PATCH    | /menus/:menu_id | Update a menu      |

### **Menu Food Management**
| Method | Endpoint         | Description              |
|--------|-----------------|--------------------------|
| GET    | /foods           | Get all food items from all menus |
| GET    | /foods/:food_id  | Get info about a specific food item from a menu |
| POST   | /foods           | Add a new food item to a menu  |
| PATCH    | /foods/:food_id  | Update a food item in a menu |

### **Order Management**
| Method | Endpoint         | Description              |
|--------|-----------------|--------------------------|
| GET    | /orders         | Get all orders           |
| GET    | /orders/:order_id   | Get a specific order |
| POST   | /orders         | Create a new order       |
| PATCH    | /orders/:order_id     | Update an order  |

### **Table Management**
| Method | Endpoint         | Description              |
|--------|-----------------|--------------------------|
| GET    | /tables         | Get all tables           |
| GET    | /tables:table_id  | Get info about a specific table |
| POST   | /tables         | Add a new table          |
| PATCH    | /tables/:table_id     | Update table status      |

### **Invoice **
| Method | Endpoint         | Description              |
|--------|-----------------|--------------------------|
| GET    | /invoices       | Get all invoices         |
| POST   | /invoices       | Generate an invoice      |
| PATCH    | /invoices/:id   | Update invoice status    |

## Testing the API
You can use **Postman** or **cURL** to test the endpoints.
Example request to get all menu items:
```sh
curl -X GET http://localhost:8080/menus
```

## Contribution
1. Fork the repository
2. Create a new branch (`git checkout -b feature-name`)
3. Commit your changes (`git commit -m 'Add new feature'`)
4. Push the branch (`git push origin feature-name`)
5. Create a Pull Request

## License
This project is licensed under the MIT License.

---
Made with ❤️ by Usama Shoukat. 

