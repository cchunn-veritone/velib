# Velib
Veritone Engine Library

## Usage

Below is a very basic example of how to use the Velib library.

```go

package main

import (
	"github.com/cchunn-veritone/velib"
)

func main() {
	velib.Init() // Initialize the Velib library with the default configuration.

  // If you need to do some sort of checks to determine if the engine is ready to process,
  // you can spawn a go routine to do those checks and then call the `Ready` function
	
  go func() {
		time.Sleep(5 * time.Second)
		velib.Ready() // Set the state of the engine as ready to process
	}()
  
  // Otherwise, you can call the `Ready` function directly.
  // velib.Ready()

  // Set the function to be called when the engine is ready to process.
	velib.Process(process)

  // Run the server (this is blocking, no code after this will execute)
	velib.Run()
}

// This function is called when /process is hit
// Inside the library, heartbeat and initial parsing of form data  are already handled by default
// See https://pkg.go.dev/net/http#ServeMux.HandleFunc for more information
func process(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Println("Processing content!")
}

```