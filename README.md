# pips

`pips` is a Go library for building concurrent data processing pipelines with strong typing through generics. It provides a clean and composable way to process streams of data through a series of stages.

## Installation

```bash
go get github.com/zhulik/pips
```

## Features

- Type-safe pipelines using Go generics
- Concurrent processing with goroutines
- Error handling throughout the pipeline
- Composable stages for common operations:
  - Map: Transform data from one type to another
  - Filter: Keep only items that match a predicate
  - Batch: Group items into batches of a specified size
  - Flatten: Expand slices into individual items
  - Each: Apply a function to each item without changing the item
  - Zip: Transform items and pair the results with the original items

## Basic Usage

```go
package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/zhulik/pips"
	"github.com/zhulik/pips/apply"
)

func main() {
	// Create a context that can be used to cancel the pipeline
	ctx := context.Background()

	// Create an input channel
	input := make(chan pips.D[string])

	// Build a pipeline that:
	// 1. Converts strings to integers
	// 2. Filters out odd numbers
	// 3. Batches results in groups of 3
	pipeline := pips.New[string, []int](
		apply.Map(func(ctx context.Context, s string) (int, error) {
			return strconv.Atoi(s)
		}),
		apply.Filter(func(ctx context.Context, n int) (bool, error) {
			return n%2 == 0, nil // Keep only even numbers
		}),
		apply.Batch[int](3),
	)

	// Run the pipeline
	output := pipeline.Run(ctx, input)

	// Feed data into the pipeline
	go func() {
		for _, s := range []string{"1", "2", "3", "4", "5", "6", "8", "10"} {
			input <- pips.NewD(s)
		}
		close(input)
	}()

	// Process results
	for batch := range output {
		numbers, err := batch.Unpack()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Printf("Batch: %v\n", numbers)
	}
}
```

## Use Cases

### Data Transformation

Use `pips` to transform data from one format to another, such as parsing JSON, converting types, or extracting fields.

```go
pipeline := pips.New[string, User](
    apply.Map(func(ctx context.Context, jsonStr string) (User, error) {
        var user User
        err := json.Unmarshal([]byte(jsonStr), &user)
        return user, err
    }),
)
```

### Filtering and Validation

Filter out invalid or unwanted data:

```go
pipeline := pips.New[User, User](
    apply.Filter(func(ctx context.Context, user User) (bool, error) {
        return user.Age >= 18, nil // Keep only adult users
    }),
)
```

### Batching for Bulk Operations

Group items for bulk processing, such as batch database inserts:

```go
pipeline := pips.New[User, []User](
    apply.Batch[User](100), // Group users in batches of 100
)
```

### Flattening Nested Data

Process nested collections by flattening them:

```go
pipeline := pips.New[[]int, int](
    apply.Flatten[int](),
)
```

### Side Effects with Each

Perform side effects on each item without changing the data flow:

```go
pipeline := pips.New[User, User](
    apply.Each(func(ctx context.Context, user User) error {
        // Log user information
        log.Printf("Processing user: %s", user.Name)
        return nil
    }),
)
```

### Pairing Original and Transformed Data with Zip

Transform items while keeping the original data paired with the result:

```go
pipeline := pips.New[string, pips.P[string, int]](
    apply.Zip(func(ctx context.Context, s string) (int, error) {
        return strconv.Atoi(s)
    }),
)

// Later, access both the original string and the parsed integer
for item := range output {
    pair, err := item.Unpack()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        continue
    }

    original, parsed := pair.Unpack()
    fmt.Printf("Original: %s, Parsed: %d\n", original, parsed)
}
```

## Error Handling

Errors are propagated through the pipeline using the `D[T]` wrapper:

```go
pipeline := pips.New[string, int](
    apply.Map(func(ctx context.Context, s string) (int, error) {
        n, err := strconv.Atoi(s)
        if err != nil {
            return 0, fmt.Errorf("invalid number: %w", err)
        }
        return n, nil
    }),
)
```

## License

This project is available under the [MIT License](LICENSE).
