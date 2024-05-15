package util

// Create a Discard struct that implements the io.Writer interface
type Discard struct{}

func (d *Discard) Write(p []byte) (n int, err error) {
	return len(p), nil
}
