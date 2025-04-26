//go:build tools
// +build tools

package generated

// https://github.com/khan/genqlient/blob/main/docs/faq.ms#genqlient-fails-after-go-mod-tidy
import (
	_ "github.com/Khan/genqlient"
)
