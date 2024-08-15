# Robin App API Backend

Join RHA Tech on Discord

[![Join RHA Tech on Discord](https://dcbadge.limes.pink/api/server/nvST6yxSEy)](https://discord.gg/nvST6yxSEy)



## Development

### Requirements
- Docker 

### With VSCode and Dev Containers (Recommended)

- clone the repo to your local machine
- open the folder in VSCode, VSCode will detect devcontainers configuration and ask to reopen in container, proceed with this option
- this will build the containers, mount the code in robin-api container and open in VSCode
- on saving any files, code will automatically rebuild (make dev is not required)
- press F5 to debug using delve

**NOTE:** you must run any make commands (for e.g. creating or applying migrations) on the host machine and not in the dev container

### Without Dev Containers

Follow these steps to build and run the application on your local device without using dev containers:

- clone the repo to your local machine
- run `make dev` to build the code and setup the database container
- run `make migrate` to apply any schema changes to the local database
- (Optional) run `make test-data` to load any new test data into the database

### Database schema changes
Follow these steps to make any changes to the database schema:

- Create new up and down migration files using `make create-migration name=<description>`. Here `description` denotes name of the changes which will become part of the filename, for e.g. `make create-migration name=create_cities`. This will add two empty files (`xxx_create_cities.up.sql` and `xxx_create_cities.down.sql`) to the `/db` directory.
- Add the up changes (e.g. CREATE TABLE / ALTER TABLE ADD) to the up file and the corresponding reversible change (e.g. DROP TABLE / ALTER TABLE DROP COLUMN) to the down file.

**NOTE:**

- Keep the changes idempotent, so use `CREATE TABLE IF NOT EXISTS <tablename>` instead of `CREATE TABLE <tablename>`
- Add one change per file. If you need to add multiple changes, create multiple migrations.
- When merging or resolving conflicts, make sure files are ordered correctly. Changes are applied in the seq order (i.e. `xxx_` prefix)

### Postman Collection
- https://www.postman.com/robinhoodarmy/workspace/backend-api/overview