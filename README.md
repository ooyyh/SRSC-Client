# SRSC-Client

A client for the SRSC project, designed to be compatible with S3-compliant object storage services. Built with Go, Wails, and Vue.js.

## Features

*(Add a brief description of the key features here)*

## Prerequisites

- Go (version 1.22 or later)
- Node.js and npm

## Getting Started

### Installation

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd SRSC-Client
    ```
2.  **Install frontend dependencies:**
    ```bash
    cd frontend
    npm install
    cd ..
    ```

### Running in Development Mode

To run the application in development mode with hot-reloading:

```bash
wails dev
```

This will build and run the application, automatically watching for changes in both the Go code and the frontend code.

## Building for Production

To build a production-ready executable:

```bash
wails build
```

This command compiles the Go code and bundles the frontend assets into a single executable located in the `build/bin` directory.

## Technologies Used

-   **Backend:** Go
-   **Frontend:** Vue.js 3, Vite
-   **Framework:** Wails v2
-   **AWS SDK:** aws-sdk-go (for S3 interaction)

## Author

-   **Name:** ooyyh
-   **Email:** laowan345@gmail.com
