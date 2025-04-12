# Go Starter Code Project Description

This project provides a robust foundation for building web applications using Go. It incorporates a variety of essential features, including user and product management, data transfer objects (DTOs), validation, file uploads, Redis caching, email sending, and robust security practices.

## Key Features

* **Modular Design:** The project is organized into `user` and `product` modules, promoting maintainability and scalability.

* **Data Transfer Objects (DTOs):** DTOs are used to define the structure of data exchanged between the application layers, ensuring data integrity and decoupling.

* **Validation:** Input data is rigorously validated using the `github.com/go-playground/validator/v10` library to prevent errors and security vulnerabilities.

* **File Uploads:** The application supports file uploads, allowing users to upload images, documents, or other files.

* **Redis Caching:** Redis is integrated using `github.com/redis/go-redis/v9` to cache frequently accessed data, improving application performance and reducing database load.

* **Email Sending:** The application can send emails using SMTP, powered by the `github.com/jordan-wright/email` library, useful for user registration, password resets, and notifications.

* **Security:**

    * **Rate Limiting:** `github.com/go-chi/httprate` is used to protect the application from abuse by limiting the number of requests a user can make within a specific time frame.

    * **CORS:** `github.com/go-chi/cors` is implemented to configure Cross-Origin Resource Sharing, allowing controlled access to the API from different domains.

    * **JWT Authentication:** User authentication is handled using JSON Web Tokens (JWT) via `github.com/golang-jwt/jwt/v4`.

* **Logging:** The application uses `github.com/rs/zerolog` for structured logging, providing detailed and easily searchable logs for debugging and monitoring.

* **Database:** PostgreSQL is used as the database, accessed via the GORM ORM (gorm.io/gorm).

* **Routing:** The Chi router (github.com/go-chi/chi/v5) is used for defining API endpoints and handling requests.

## Tech Stack

* **Language:** Go 1.21+

* **Framework:** Chi router (https://github.com/go-chi/chi/v5)

* **Database:** PostgreSQL with GORM (gorm.io/gorm)

* **Authentication:** JWT (https://github.com/golang-jwt/jwt/v4)

* **Caching:** Redis (https://github.com/redis/go-redis/v9)

* **Email:** SMTP with https://github.com/jordan-wright/email

* **Validation:** https://github.com/go-playground/validator/v10

* **File Upload:**

* **Security:**

    * Rate limiting: https://github.com/go-chi/httprate
    * CORS: https://github.com/go-chi/cors

* **Logging:** https://github.com/rs/zerolog

This starter code provides a solid foundation for building scalable, secure, and performant web applications with Go. It handles many common requirements out-of-the-box, allowing developers to focus on business logic.
