# Record Emitter

### Usage

```golang

file := "/path/to/file.csv"
e := record_emitter.NewEmitter(file)
receiver := e.Start()

for row := range receiver {
   // mess with record
}

// Alternatively

for row := range e.Start() {
   // mess with record
}
```