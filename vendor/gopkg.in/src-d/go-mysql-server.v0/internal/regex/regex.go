package regex

import errors "gopkg.in/src-d/go-errors.v1"

var (
	// ErrRegexAlreadyRegistered is returned when there is a previously
	// registered regex engine with the same name.
	ErrRegexAlreadyRegistered = errors.NewKind("Regex engine already registered: %s")
	// ErrRegexNameEmpty returned when the name is "".
	ErrRegexNameEmpty = errors.NewKind("Regex engine name cannot be empty")
	// ErrRegexNotFound returned when the regex engine is not registered.
	ErrRegexNotFound = errors.NewKind("Regex engine not found: %s")

	registry      map[string]Constructor
	defaultEngine string
)

// Matcher interface is used to compare regexes with strings.
type Matcher interface {
	// Match returns true if the text matches the regular expression.
	Match(text string) bool
}

// Constructor creates a new Matcher.
type Constructor func(re string) (Matcher, error)

// Register add a new regex engine to the registry.
func Register(name string, c Constructor) error {
	if registry == nil {
		registry = make(map[string]Constructor)
	}

	if name == "" {
		return ErrRegexNameEmpty.New()
	}

	_, ok := registry[name]
	if ok {
		return ErrRegexAlreadyRegistered.New(name)
	}

	registry[name] = c

	return nil
}

// Engines returns the list of regex engines names.
func Engines() []string {
	var names []string

	for n := range registry {
		names = append(names, n)
	}

	return names
}

// New creates a new Matcher with the specified regex engine.
func New(name, re string) (Matcher, error) {
	n, ok := registry[name]
	if !ok {
		return nil, ErrRegexNotFound.New(name)
	}

	return n(re)
}

// Default returns the default regex engine.
func Default() string {
	if defaultEngine != "" {
		return defaultEngine
	}

	_, ok := registry["oniguruma"]
	if ok {
		return "oniguruma"
	}

	return "go"
}

// SetDefault sets the regex engine returned by Default.
func SetDefault(name string) {
	defaultEngine = name
}
