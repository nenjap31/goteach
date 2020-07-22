# goteach

What is goteach?
--------------
goteach is golang project template, use [echo](https://echo.labstack.com) golang framework.

Installation for Development
----------------------------
Requirement:
- Go Programming language
- MySql Database

Installation:
1. Clone from github repository: git clone https://github.com/nenjap31/goteach.git
2. Install [Dep](https://golang.github.io/dep/docs/introduction.html) (Dependency Management Tool): curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
3. Run Command to get the dependencies: dep ensure
4. Run Command to update existing dependencies: dep ensure -update
5. Set Enviroment Variable from source code directory:
    - export goteach_ENV=development
6. Create mysql database schema.
7. copy cfg/example.json to cfg/dev.json and configure it
8. Run mysql database migration: go run main.go migrate up
9. Run mysql data seeder: go run main.go seed
10. Run it without compile:
    - HTTP Service: go run main.go start
12. Compile it with command: go build
12. Run binary with command:
    - HTTP Service: ./goteach start

Folder Structure
----------------
    - cfg: config files location
    - cmd: location for code to make CLI command
    - config: code location for config files reader
    - controllers: location for Accepts input and converts it to commands for the model or view.
    - database: location for migration and seeder code
    - logger: location for log handling
    - models: location for database data model & data logic code
    - modules: location for modular logic code
    - presenter: location for api response data model
    - routes: location for http routing
    - script: location for CLI script
    - utils: location for other utility code

Testing
-------
1. Create mysql database schema for testing
2. copy cfg/example.json to cfg/testing.json and configure it
3. Run test script from project root:
    - All Test: ./script/testing.sh ./...
    - Single Test: ./script/testing.sh ./folder/folder/
4. Or run command manually from project root:
    - Set Enviroment Variable from source code directory:
        - export GOTEACH_ENV=testing
        - export GOTEACH_APP_PATH=$(pwd)
    - Reset testing database: go run main.go migrate reset
    - Database migration: go run main.go migrate up
    - Run data seeder: go run main.go seed
    - Run test command:
        - All Test: go test ./...
        - Single Test: go test ./folder/folder/