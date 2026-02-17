# Rootstock — Engineering Rules

## Metalogic and behavior
- If you cannot provide an example or a reference, you are guessing. Stop and look for something verifiable to build on.

## Patterns and Consistency
- Consistency is one of the most important design principles of this project. Verify all changes against design principles.

## Go Concurrency
- Never use `sync.Mutex` or `sync.RWMutex`. Use goroutines and channels for all concurrency.


