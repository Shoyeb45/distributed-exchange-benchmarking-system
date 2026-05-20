# Distributed Exchange Benchmarking System

- Distributed benchmarking and stress-testing platform for trading infrastructure with sandboxed execution, massive bot fleets, real-time telemetry, and live leaderboard analytics.

## Development Setup

- This project uses taskfile orchestration in this monorepo. 
- Each services have their own toolchain, and the central `Taskfile.yml` controls everything. 
- Each apps/services has their own `Taskfile.yml`.


- Following are the set of commands which is present in almost all the services:

    - `dev`: Run this command for development, the services may have hot reload
    - `start`: Run this command for production ready execution of services
    - `build`: Optimized build for services
    - `lint`: To format the code
    - `test`: To test the code

- Following are the project names which can be used to run the specific command for the specific service(`service_name:command`):

    - **sandbox-engine**: Rust based sandbox engine to execute the contestants code
    - **web**: The web interface of the entire system
    - **server**: Central api server and orchestrator