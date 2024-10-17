<img src="./web/src/lib/assets/logo.png" width="80">

# SSH Nexus

Welcome to **SSH Nexus**, your one-stop solution for managing SSH access with ease and security. This application combines a powerful backend using Pocketbase & Go with a sleek and responsive frontend built with Sveltekit, delivering a user-friendly experience for both administrators and users.

## Features

- **User Authentication**

  - **OAuth**: Seamless sign-up and login using popular OAuth providers.
  - **Email**: Traditional email-based sign-up and login.

- **Admin Capabilities**

  - **User Management**: Easily manage user accounts and permissions.
  - **Machine Management**: Add, edit, and remove machines from the network.
  - **Group Management**: Organize machines and users into groups for streamlined access control.

- **User Access Control**

  - **Machine Assignment**: Assign users to machines directly or through groups.
  - **Group Assignment**: Add users to groups to inherit access to all machines within the group.

- **SSH Certificate Management**

  - **Automated User CA**: Automatically add a user Certificate Authority to each machine.
  - **Principal IDs**: Manage principal IDs for users seamlessly.
  - **Key Signing**: Users can sign their own SSH keys with customizable expiry settings via the UI.

- **Automatic Updates and Clean-up**

  - **Agents**: Updates and configurations are handled via small agents that are installed on the machines automatically.
  - **Self-Destructing Agents**: When a machine is removed from the server, the agent destroys itself and all associated files, ensuring the machine remains clean.

## Installation

### Prerequisites

- Ensure your server has SSH access to the target machines for initial agent installation. SSH Nexus will generate a random ssh key on startup, which you can change afterwards. This single key serves as a way to install agents on new machines and will be used as the user certificate authority.

### Backend Setup

#### Method 1: Direct Installation

1. **Download the latest server binary** for your OS from the [releases page](https://github.com/MizuchiLabs/ssh-nexus/releases).
   Or use the install script:
   ```bash
   # For the latest server release
   curl -sSL https://github.com/MizuchiLabs/ssh-nexus/raw/branch/main/install.sh | bash

   # For the agent
   curl -sSL https://github.com/MizuchiLabs/ssh-nexus/raw/branch/main/install.sh | bash -s agent
   ```
1. **Environment variables**: Before starting you will need to set at least the following 2 environment variables:
   - **PB_ADMIN_PASSWORD**: The password for the admin user.
   - **PB_ENCRYPTION_KEY**: The encryption key for the sqlite3 database.
1. **Running the server**:
   ```bash
   ./nexus serve
   ```
1. **Running the agent** on a machine:
   ```bash
   ./nexus-agent --server <your-server-address>
   ```

#### Method 2: Docker

1. **Use the docker compose file or manually below**

1. **Pull the Docker image**:

   ```bash
   docker pull ghcr.io/mizuchilabs/ssh-nexus:latest
   ```

1. **Run the Docker container**:

   ```bash
   docker run --name ssh-nexus -d -p 8090:8090 -p 8091:8091 \
     --env PB_ADMIN_PASSWORD="password" --env PB_ENCRYPTION_KEY="some-random-key" \
      ghcr.io/mizuchilabs/ssh-nexus:latest
   ```

## Usage

### Admin Panel

The admin panel will be available at [http://localhost:8090/\_/](http://localhost:8090/_/) which uses pocketbase as the backend. You will rarely need to interact with it since everything can be done via the frontend and env variables. Be careful though when interacting with the generated collections!

### User Access

The web ui will be available at [http://localhost:8090](http://localhost:8090)

1. **Login** If you set up an OAuth provider it will show up and users can  sign in immediately. Admins can create new users manually for plain username/email and password authentication.
1. **Machines**: See a list of machines you have access to.
1. **Users**: See a list of users and their permissions.
1. **Sign**: Generate and sign your own SSH keys with an optional expiry time.
1. **System**: View various settings, tokens used by agents, keys and certificates.

## Contributing

We welcome contributions to improve SSH Nexus. To get started, fork the repository and create a new branch for your feature or bug fix.

1. **Fork the repository**
1. **Create a new branch**:
   ```bash
   git checkout -b feature-name
   ```
1. **Commit your changes**:
   ```bash
   git commit -m 'Add some feature'
   ```
1. **Push to the branch**:
   ```bash
   git push origin feature-name
   ```
1. **Open a pull request**

## License

SSH Nexus is released under the Apache 2.0 License. See the [LICENSE](LICENSE) file for more details.

# Default Project Environment Variables

```env
# General settings
export PB_APP_URL=""            
export PB_LOG_MAX_DAYS="30"       
export PB_ADMIN_EMAIL="root@nexus.local"        
export PB_ADMIN_PASSWORD="required!"     
export PB_SENDER_NAME="SSH Nexus"        
export PB_SENDER_EMAIL="no-reply@nexus.local"       
export PB_SMTP_ENABLED="false"       
export PB_SMTP_HOST=""          
export PB_SMTP_PORT="587"          
export PB_SMTP_USER=""          
export PB_SMTP_PASSWORD=""      
export PB_SMTP_TLS="true"           
export PB_S3_ENABLED="false"         
export PB_S3_ENDPOINT=""        
export PB_S3_REGION=""          
export PB_S3_BUCKET=""          
export PB_S3_SECRET=""          
export PB_S3_ACCESS_KEY=""      
export PB_S3_FORCE_PATH_STYLE=""
export PB_OIDC_URL="" # Custom oidc endpoint
export PB_OIDC_NAME=""         
export PB_OIDC_REALM="master" # Only used for Keycloak
export PB_OIDC_CLIENT_ID=""     
export PB_OIDC_CLIENT_SECRET="" 
export PB_ENCRYPTION_KEY="required!"

# Custom Repo (for private repos and forks)
export PB_REPO_URL=""
export PB_REPO_OWNER=""
export PB_REPO_NAME="ssh-nexus"
export PB_REPO_TOKEN=""
```

## Roadmap:

- Add bastion host mode
- Add request system
- Add host certificate management
- UI overhaul (e.g. charts, status, etc)
- Better user views
- Better error/debug views
- Fix realtime subscriptions
