# TxApp

## Installation Guide
### 1. Installing Go
Follow the installation instruction depending on your os from [here](https://go.dev/doc/install)

### 2. Setting Up Environment Variables
#### Windows
1. Open System Properties > Advanced > Environment Variables
2. Under System Variables, find PATH and click Edit
3. Add the appropriate paths, for example:
   - `C:\Go\bin`
   - `C:\Program Files\Go\bin`
   - `%USERPROFILE%\go\bin`
4. Click OK to save
5. Restart any open terminal windows

#### macOS
1. Open your shell profile file (`~/.zshrc` for Zsh or `~/.bash_profile` for Bash)
2. Add these lines:
   ```bash
   export GOPATH=$HOME/go
   export PATH=$PATH:$GOPATH/bin
   ```
3. Reload your profile:
   ```bash
   source ~/.zshrc   # for Zsh
   # or
   source ~/.bash_profile   # for Bash
   ```

#### Linux
1. Open your `~/.bashrc` file
2. Add these lines:
   ```bash
   export GOPATH=$HOME/go
   export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
   ```
3. Reload your profile:
   ```bash
   source ~/.bashrc
   ```

### 3. Installing Make

#### Windows
1. Using winget (Recommended):
   ```bash
   winget install GnuWin32.Make
   ```
   
   Alternatively:
   1. Install [Chocolatey](https://chocolatey.org/install)
   2. Run: `choco install make`

#### macOS
Make comes pre-installed with Xcode Command Line Tools. If needed:
```bash
xcode-select --install
```

#### Linux (Ubuntu/Debian)
```bash
sudo apt-get update
sudo apt-get install make
```

### 4. Installing Docker

For Docker and Docker Compose installation instructions specific to your operating system, please follow the official Docker documentation:
[Get Docker](https://docs.docker.com/get-started/get-docker)
> **Notes**:
> - Docker Desktop (Windows/macOS) includes Docker Compose by default
> - For Linux installations, Docker Compose is included as the `docker compose` plugin

## Setup and Running the Application

1. Install dependencies and set up the development environment:
    ```bash
    make setup
    ```

2. Run the web application:
    #### Development mode (with live reload):
    ```bash
    make app/web/live
    ```

    #### Production mode:
    ```bash
    make app/web/build
    make app/web/bin
    ```

    #### Using Docker:
    ```bash
    make docker
    ```

## Development Commands

To view all available commands and their descriptions, run:
```bash
make help
```
