# SRSC-Client

A client for the SRSC project, designed to be compatible with S3-compliant object storage services. Built with Go, Wails, and Vue.js.

## Features

You can choose to use the latest release directly, which works well, at least for programmers.
I have tested the commonly used Cloudflare R2 bucket and BackBlaze B2 bucket in the v0.0.1 release, as well as the object storage of Colorful Cloud. They all performed quite well. You only need to add the AccessKey and SecretKey provided by the service provider using the Add Node function.
The AccessKey and SecertKey of the B2 bucket may be called differently, but the AccessKey is simply changed to AplicationID and the SecertKey is changed to ApplicationKey.
Before filling in these, make sure there is a config.json in your directory. If not, it will be created automatically. It will save the login credentials of all your object storage accounts.

Let me give you a simple example
![image](https://github.com/user-attachments/assets/8f4aa0ea-c171-4f7f-bcb2-f0358298dd30)
After filling in the form, click Add.
Then you will get a node.
![image](https://github.com/user-attachments/assets/ce5be80a-542f-47f7-847e-e8a6bc2ce2ba)
If the operation is correct, click View Node, and you can use the key you filled in to enter the warehouse account, and all the warehouses you own under the endpoint you filled in will be listed, as well as detailed information about each of your warehouses
![image](https://github.com/user-attachments/assets/81125176-2bbb-4b57-ac2a-6ecec0053922)
Click on a bucket to enter the bucket and view all objects in the bucket.
![image](https://github.com/user-attachments/assets/83e8ca18-b599-408b-a8cf-a260d6955eb5)
Now you can perform various operations on the objects in these buckets, including uploading and downloading, and viewing more detailed information.
![image](https://github.com/user-attachments/assets/fe04a460-29b5-4431-ae87-a69037328389)

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
