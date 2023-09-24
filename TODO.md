### Todo

- [ ] Better error management, errors send to Discord are too descriptive
- [ ] Add more logs for each level
- [ ] Add tracing for each request
- [ ] Add API as controller
- [ ] Fix In UpdatePlayer: placeholder $1 already has type int, cannot assign varchar
- [ ] Fix in SearchPlayer: player strikes must be import there not in usecase
- [ ] Add notes to players (for example: player cannot play on wednesday)

### In Progress

- [ ] Add comments to every functions
- [ ] Add tests to every functions

### Done ✓
- [X] Create CI linter
- [X] Create Release Pipeline
- [X] Create Readme
- [X] Add Validate for entities
- [X] Better context timeout management
- [X] Add Log Level to config
- [X] Fix DeleteStrike always return success
- [X] Add In CreatePlayer id to player entity
- [X] Add Usecase: player name is not discord name. Must implement a way to link them
