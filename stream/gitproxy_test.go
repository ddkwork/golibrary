package stream

import "testing"

func TestSetGitProxy(t *testing.T) {
	GitProxy(true)
}

func TestUnSetGitProxy(t *testing.T) {
	GitProxy(false)
}
