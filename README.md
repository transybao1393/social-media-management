# Social Media Management
## High level architect
![Alt text](/clean-architect.png "Clean architecture design")

### Interpretation
app/                  → This folder contains general files that applied for all code, including config, connector, logging, utilities

config/        → This folder contains config file in json or .env (optional)

connector/ → This folder contains connector including router handling.

logger/        → This folder contains logger initialize.

utils/            → This folder contains utility helpers.

logs/            → This folder contains log files.

bin/                    → This folder contains air execution or you can install air globally.

domain/             → This folder contains db instances and other models.

dbInstance/ → This folder contains any database instance (redis, mongodb, mysql, postgres,…).

<another models>

tiktok/                 → This folder contains all API management from Tiktok.

delivery/        → This folder contains contains endpoint functions for OAuth,..getting from usecase folder.

repository/    → This folder contains code to communicate and data processing from database.

usecase/       → This folder contains logic/business code, OAuth, data retrieva, update,…test cases for every usecase.

<test files for every usecases>

### Prerequisites
- Go v1.20
### How to run
#### Docker
First time use
```
docker compose up --build
```
Usually, we just need:
```
docker compose up
```
Note: When run using Docker, Docker itseft can be automatically download all dependencies for us.
#### Build from source
```
# Step 1: Download dependencies
go mod download

# Step 2: Build from source with hot reloading
./bin/air -c .air.toml 
```