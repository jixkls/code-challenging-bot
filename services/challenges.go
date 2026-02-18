package services

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"telegram-bot/models"
)

var (
	activeChallenges = make(map[int64]*models.ActiveChallenge)
	challengeMu      sync.Mutex
)

// challengePool holds all seeded challenges organized by difficulty.
var challengePool = map[string][]models.Challenge{
	"easy": {
		{
			ID:         "e1",
			Title:      "Even or Odd",
			Difficulty: "easy",
			Description: "What does this function return when called with n = 7?",
			CodeSnippet: `func check(n int) string {
    if n%2 == 0 {
        return "even"
    }
    return "odd"
}`,
			Options:    [4]string{"A) \"even\"", "B) \"odd\"", "C) 1", "D) Compilation error"},
			CorrectIdx: 1,
			Hint:       "The % operator returns the remainder of division. 7 divided by 2 has remainder 1.",
		},
		{
			ID:         "e2",
			Title:      "Zero Value",
			Difficulty: "easy",
			Description: "What is the zero value of an int variable in Go?",
			CodeSnippet: `var x int
fmt.Println(x)`,
			Options:    [4]string{"A) nil", "B) undefined", "C) 0", "D) Compilation error"},
			CorrectIdx: 2,
			Hint:       "Go initializes all variables to their zero value. For numeric types, think of what 'zero' means.",
		},
		{
			ID:         "e3",
			Title:      "Short Variable Declaration",
			Difficulty: "easy",
			Description: "Which line correctly declares and initializes a string variable in Go?",
			CodeSnippet: "// Pick the correct declaration:",
			Options:    [4]string{"A) var name = string(\"Go\")", "B) name := \"Go\"", "C) string name = \"Go\"", "D) let name = \"Go\""},
			CorrectIdx: 1,
			Hint:       "Go uses := for short variable declarations inside functions. It infers the type automatically.",
		},
		{
			ID:         "e4",
			Title:      "For Loop Output",
			Difficulty: "easy",
			Description: "How many times does \"hello\" get printed?",
			CodeSnippet: `for i := 0; i < 3; i++ {
    fmt.Println("hello")
}`,
			Options:    [4]string{"A) 2", "B) 3", "C) 4", "D) Infinite loop"},
			CorrectIdx: 1,
			Hint:       "The loop starts at i=0 and runs while i < 3. Count: 0, 1, 2.",
		},
		{
			ID:         "e5",
			Title:      "String Length",
			Difficulty: "easy",
			Description: "What does len(s) return?",
			CodeSnippet: `s := "Golang"
fmt.Println(len(s))`,
			Options:    [4]string{"A) 5", "B) 6", "C) 7", "D) Error: len not defined for string"},
			CorrectIdx: 1,
			Hint:       "len() returns the number of bytes in a string. Count the letters in \"Golang\".",
		},
	},
	"medium": {
		{
			ID:         "m1",
			Title:      "Slice Append",
			Difficulty: "medium",
			Description: "What is the output of this code?",
			CodeSnippet: `s := []int{1, 2, 3}
s = append(s, 4, 5)
fmt.Println(len(s), cap(s))`,
			Options:    [4]string{"A) 5 5", "B) 5 6", "C) 3 5", "D) 5 8"},
			CorrectIdx: 1,
			Hint:       "append() grows the slice. The capacity doubles when the underlying array is full (cap goes from 3 to 6).",
		},
		{
			ID:         "m2",
			Title:      "Map Access",
			Difficulty: "medium",
			Description: "What does this code print?",
			CodeSnippet: `m := map[string]int{"a": 1, "b": 2}
val, ok := m["c"]
fmt.Println(val, ok)`,
			Options:    [4]string{"A) 0 false", "B) nil false", "C) Panic: key not found", "D) 0 true"},
			CorrectIdx: 0,
			Hint:       "Accessing a missing map key returns the zero value for the value type plus false for the ok flag.",
		},
		{
			ID:         "m3",
			Title:      "Struct Method",
			Difficulty: "medium",
			Description: "What does r.Area() return?",
			CodeSnippet: `type Rect struct {
    Width, Height float64
}

func (r Rect) Area() float64 {
    return r.Width * r.Height
}

r := Rect{Width: 3, Height: 4}`,
			Options:    [4]string{"A) 7", "B) 12", "C) 12.0", "D) Compilation error"},
			CorrectIdx: 1,
			Hint:       "Area multiplies Width * Height. Both are float64 but 3 * 4 = 12 which prints as 12 (not 12.0) with Println.",
		},
		{
			ID:         "m4",
			Title:      "Error Handling",
			Difficulty: "medium",
			Description: "What is the idiomatic way to handle an error from a function in Go?",
			CodeSnippet: `result, err := doSomething()
// What goes here?`,
			Options: [4]string{
				"A) if err != nil { return err }",
				"B) try { result } catch(err)",
				"C) if err { throw err }",
				"D) err.handle()",
			},
			CorrectIdx: 0,
			Hint:       "Go doesn't have try/catch. Errors are values returned from functions and checked explicitly.",
		},
		{
			ID:         "m5",
			Title:      "Defer Order",
			Difficulty: "medium",
			Description: "What is the output?",
			CodeSnippet: `func main() {
    defer fmt.Println("first")
    defer fmt.Println("second")
    defer fmt.Println("third")
}`,
			Options:    [4]string{"A) first second third", "B) third second first", "C) third first second", "D) first third second"},
			CorrectIdx: 1,
			Hint:       "Deferred calls are executed in LIFO (Last In, First Out) order — like a stack.",
		},
	},
	"hard": {
		{
			ID:         "h1",
			Title:      "Goroutine Output",
			Difficulty: "hard",
			Description: "What is the MOST LIKELY output?",
			CodeSnippet: `func main() {
    go fmt.Println("hello")
    fmt.Println("world")
}`,
			Options:    [4]string{"A) hello world", "B) world hello", "C) world", "D) hello"},
			CorrectIdx: 2,
			Hint:       "The main goroutine doesn't wait for other goroutines to finish. It might exit before the goroutine prints.",
		},
		{
			ID:         "h2",
			Title:      "Channel Deadlock",
			Difficulty: "hard",
			Description: "What happens when this code runs?",
			CodeSnippet: `func main() {
    ch := make(chan int)
    ch <- 42
    fmt.Println(<-ch)
}`,
			Options:    [4]string{"A) Prints 42", "B) Prints 0", "C) Deadlock (fatal error)", "D) Compilation error"},
			CorrectIdx: 2,
			Hint:       "An unbuffered channel blocks on send until another goroutine reads. Here, main is both sender and receiver — it blocks forever.",
		},
		{
			ID:         "h3",
			Title:      "Interface Satisfaction",
			Difficulty: "hard",
			Description: "Does Dog satisfy the Animal interface?",
			CodeSnippet: `type Animal interface {
    Speak() string
}

type Dog struct{}

func (d *Dog) Speak() string {
    return "Woof"
}

var a Animal = Dog{}`,
			Options:    [4]string{"A) Yes, compiles fine", "B) No, compile error: Dog doesn't implement Animal", "C) Runtime panic", "D) Yes, but Speak() returns empty string"},
			CorrectIdx: 1,
			Hint:       "The method is defined on *Dog (pointer receiver), but we're assigning a Dog value. Only *Dog satisfies the interface.",
		},
		{
			ID:         "h4",
			Title:      "Closure Trap",
			Difficulty: "hard",
			Description: "What does this print?",
			CodeSnippet: `funcs := make([]func(), 3)
for i := 0; i < 3; i++ {
    funcs[i] = func() { fmt.Println(i) }
}
for _, f := range funcs {
    f()
}`,
			Options:    [4]string{"A) 0 1 2", "B) 3 3 3", "C) 2 2 2", "D) 0 0 0"},
			CorrectIdx: 1,
			Hint:       "The closure captures the variable i, not its value. By the time the functions run, the loop has finished and i == 3.",
		},
		{
			ID:         "h5",
			Title:      "Select Statement",
			Difficulty: "hard",
			Description: "What happens with this select?",
			CodeSnippet: `ch1 := make(chan string)
ch2 := make(chan string)

go func() {
    ch1 <- "one"
}()

select {
case msg := <-ch1:
    fmt.Println(msg)
case msg := <-ch2:
    fmt.Println(msg)
}`,
			Options:    [4]string{"A) Always prints \"one\"", "B) Deadlock", "C) Prints \"one\" (usually)", "D) Random between ch1 and ch2"},
			CorrectIdx: 2,
			Hint:       "select picks whichever channel is ready first. The goroutine sends on ch1, so that case fires — but the timing isn't guaranteed by the spec.",
		},
	},
}

// GetRandomChallenge returns a challenge for the given difficulty.
// Tries LLM-generated challenge first, falls back to static pool on error.
// Returns nil if the difficulty is not recognized.
func GetRandomChallenge(difficulty string) *models.Challenge {
	// Try LLM-generated challenge first
	challenge, err := GenerateChallenge(difficulty)
	if err == nil {
		return challenge
	}
	log.Printf("Gemini fallback: %v", err)

	// Fall back to static pool
	pool, ok := challengePool[difficulty]
	if !ok || len(pool) == 0 {
		return nil
	}
	idx := rand.Intn(len(pool))
	staticChallenge := pool[idx]
	return &staticChallenge
}

// SetActiveChallenge stores an active challenge session for a user.
func SetActiveChallenge(telegramID int64, challenge *models.Challenge) {
	challengeMu.Lock()
	defer challengeMu.Unlock()

	activeChallenges[telegramID] = &models.ActiveChallenge{
		UserTelegramID: telegramID,
		Challenge:      challenge,
		HintUsed:       false,
		StartedAt:      time.Now(),
	}
}

// GetActiveChallenge retrieves the current active challenge for a user.
// Returns nil if no active challenge exists.
func GetActiveChallenge(telegramID int64) *models.ActiveChallenge {
	challengeMu.Lock()
	defer challengeMu.Unlock()

	return activeChallenges[telegramID]
}

// ClearActiveChallenge removes the active challenge for a user.
func ClearActiveChallenge(telegramID int64) {
	challengeMu.Lock()
	defer challengeMu.Unlock()

	delete(activeChallenges, telegramID)
}
