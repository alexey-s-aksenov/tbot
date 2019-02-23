package joke

// Getter interface to decouple different sources
type Getter interface {
	GetJoke() (string, error)
}

// NewJokeGetter returns interface to get joke
func NewJokeGetter() Getter {
	return &bashorg{}
}
