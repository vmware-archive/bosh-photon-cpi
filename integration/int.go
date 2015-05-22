package inttests

import (
	"math/rand"
	"time"
)

func init() {
	// Initialize random seed used for naming tenant, resticket, etc
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
}
