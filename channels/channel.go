package channels

import (
	"github.com/iamsayantan/larasockets"
	"strings"
)

// NewChannel is a factory method that returns the appropriate channel based on
// the channel name.
func NewChannel(name string) larasockets.Channel {
	if strings.HasPrefix(name, "private-") {
		return newPrivateChannel(name)
	}

	return newPublicChannel(name)
}
