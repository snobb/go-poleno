package termcolour

const reset = "\x1b[0m"

var mapping = map[string]string{
	"white": "\x1b[1m",
	"grey":  "\x1b[2m",

	"black":   "\x1b[30m",
	"red":     "\x1b[31m",
	"green":   "\x1b[32m",
	"yellow":  "\x1b[33m",
	"blue":    "\x1b[34m",
	"magenta": "\x1b[35m",
	"cyan":    "\x1b[36m",
}

func Lookup(colour string) string {
	return mapping[colour]
}

func Reset() string {
	return reset
}

func Names() []string {
	names := make([]string, 0, len(mapping))
	for k := range mapping {
		names = append(names, k)
	}
	return names
}
