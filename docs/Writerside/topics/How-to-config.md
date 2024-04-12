# How to Set Up Your Development Environment for Cortex

This guide is designed to help you set up the development environment for Cortex, a comprehensive platform for data governance and cataloging. By following these steps, you will configure the necessary tools and services, including Elasticsearch, Kibana, and the Cortex application itself, both backend and frontend components.

## Before You Start

Before diving into the setup process, ensure you have the following prerequisites installed on your machine:

- Docker and Docker Compose: Cortex uses Docker containers for Elasticsearch, Kibana, MariaDB, and optionally for the Cortex application itself.
- Node.js and Yarn: These are required to build and run the Cortex frontend.
- Go: Required for developing and building the Cortex backend.

## Setting Up the Development Environment

### Step 1: Clone the Repository

Start by cloning the Cortex repository to your local machine. Open a terminal and run:

```bash
git clone https://github.com/neural-pulse/cortex.git
cd cortex
```

### Step 2: Start Elasticsearch and Kibana

Cortex relies on Elasticsearch for data storage and search capabilities, and Kibana for data visualization.

1. Navigate to the root of the cloned repository.
2. Start the services using Docker Compose:

   ```bash
   docker-compose up -d elasticsearch kibana
   ```

3. Wait for the services to start up. You can check the status using:

   ```bash
   docker-compose logs -f elasticsearch kibana
   ```

   Exit the log view by pressing `Ctrl+C`.

### Step 3: Configure Kibana (Optional)

Kibana is used for visualizing Elasticsearch data. To configure Kibana:

1. Open Kibana by navigating to `http://localhost:5601` in your web browser.
2. Follow the on-screen instructions to connect Kibana to your Elasticsearch instance. The default settings should work if you haven't changed any configurations in `docker-compose.yml`.

### Step 4: Set Up the Backend

The Cortex backend is a Go application that handles data ingestion and API requests.

1. Ensure Go is installed and set up on your machine.
2. Navigate to the backend directory:

   ```bash
   cd backend/cmd/app
   ```

3. Build the backend application:

   ```bash
   go build -o cortex-backend
   ```

4. Run the backend:

   ```bash
   ./cortex-backend
   ```

### Step 5: Set Up the Frontend

The Cortex frontend is built with Next.js and React.

1. Ensure Node.js and Yarn are installed on your machine.
2. Navigate to the frontend directory:

   ```bash
   cd frontend
   ```

3. Install the dependencies:

   ```bash
   yarn install
   ```

4. Start the development server:

   ```bash
   yarn dev
   ```

5. Access the Cortex web interface at `http://localhost:3000`.

## Contributing

Contributions to Cortex are welcome! Feel free to open issues or submit pull requests for improvements.

## Conclusion

You've now set up your development environment for Cortex, including the backend and frontend components, as well as Elasticsearch and Kibana for data storage and visualization. This setup will allow you to develop and test features locally.