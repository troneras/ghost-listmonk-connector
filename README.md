# Ghost-Listmonk Connector

![Ghost-monk](ghost_monk.png "Ghost Monk")

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Frontend](#frontend)
- [Contributing](#contributing)
- [License](#license)

## Introduction

Ghost-Listmonk Connector is an open-source project that bridges the gap between Ghost CMS and Listmonk, enabling seamless integration of email marketing capabilities with your Ghost blog. This connector allows you to automate subscriber management, trigger email campaigns based on Ghost events, and provides a user-friendly dashboard for monitoring and managing your email marketing efforts.

## Features

- Automatic synchronization of Ghost subscribers with Listmonk
- Trigger-based actions for various Ghost events (e.g., new post published, new member registered)
- Delayed execution of actions. You can use this to create mail chains. For example send a new subscriber emails a day later, a week later, etc.
- Customizable email templates and campaigns (In Listmonk)
- Real-time dashboard for monitoring Son (Subscriber Operations Notifier) performance
- Webhook management for Ghost events
- Caching system for improved performance
- User authentication and authorization

## Demo
(Click the image to see on youtube)
[![GHOST LISTMONK connector video](https://img.youtube.com/vi/XBhrdeZwpqI/0.jpg)](https://www.youtube.com/watch?v=XBhrdeZwpqI)

## Tech Stack

### Backend

- [Gin](https://github.com/gin-gonic/gin): Web framework
- [JWT-Go](https://github.com/golang-jwt/jwt): JSON Web Token authentication
- [MySQL](https://www.mysql.com/): Relational database
- [Redis](https://redis.io/): Caching and task queue
- [Asynq](https://github.com/hibiken/asynq): Distributed task queue
- [Logrus](https://github.com/sirupsen/logrus): Structured logger

### Frontend

- [Next.js](https://nextjs.org/): React framework
- [Tailwind CSS](https://tailwindcss.com/): Utility-first CSS framework
- [shadcn/ui](https://ui.shadcn.com/): UI component library
- [Recharts](https://recharts.org/): Charting library

## Prerequisites

- Go 1.21 or later
- Node.js 18 or later
- MySQL 8.0 or later
- Redis 6.0 or later
- Ghost CMS instance
- Listmonk instance

## Installation

1. Clone the repository:

   ```
   git clone https://github.com/yourusername/ghost-listmonk-connector.git
   cd ghost-listmonk-connector
   ```

2. Install backend dependencies:

   ```
   go mod tidy
   ```

3. Install frontend dependencies:

   ```
   cd ui
   npm install
   ```

4. Build the project:
   ```
   make build-all
   ```

## Configuration

1. Copy the example environment file and edit it with your settings:

   ```
   cp .env.example .env
   ```

2. Set up your database:

   ```
   make migrate-up
   ```

3. Configure your Ghost webhook to point to your connector's webhook endpoint.

4. Set up Listmonk API credentials in the `.env` file.

## Usage

1. Start the server:

   ```
   ./main
   ```

2. Access the dashboard at `http://localhost:8808` (or your configured port).

3. Log in using the credentials set in your `.env` file.

4. Create and manage Sons (Subscriber Operations Notifiers) through the dashboard.

## API Endpoints

- `POST /api/auth/magic-link`: Request a magic link for authentication
- `GET /api/auth/verify`: Verify magic link and authenticate user
- `GET /api/sons`: List all Sons
- `POST /api/sons`: Create a new Son
- `GET /api/sons/:id`: Get details of a specific Son
- `PUT /api/sons/:id`: Update a Son
- `DELETE /api/sons/:id`: Delete a Son
- `GET /api/webhook-logs`: Get webhook logs
- `GET /api/son-execution-logs`: Get Son execution logs
- `GET /api/son-stats`: Get Son performance statistics

For a complete API documentation, please refer to the [API Documentation](./docs/API.md).

## Frontend

The frontend is built with Next.js and provides a user-friendly interface for managing Sons, viewing logs, and monitoring performance. Key features include:

- Dashboard with recent activity and Son performance charts
- Son creation and management interface
- Webhook log viewer
- Son execution log viewer

To start the frontend in development mode:

```
cd ui
npm run dev
```

## Contributing

We welcome contributions to the Ghost-Listmonk Connector! Please follow these steps to contribute:

1. Fork the repository
2. Create a new branch: `git checkout -b feature/your-feature-name`
3. Make your changes and commit them: `git commit -m 'Add some feature'`
4. Push to the branch: `git push origin feature/your-feature-name`
5. Submit a pull request

Please make sure to update tests as appropriate and adhere to the [Code of Conduct](CODE_OF_CONDUCT.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
