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

### How to run
#### Docker
...
#### Build from source
...