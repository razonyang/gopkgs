package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackageDocsURL(t *testing.T) {
	pkg := &Package{
		Prefix: "foo.bar/fizz",
	}
	assert.Equal(t, "https://pkg.go.dev/foo.bar/fizz?tab=doc", pkg.DocsURL())

	docsURL := "https://docs.foo.bar/fizz"
	pkg.Docs = docsURL
	assert.Equal(t, docsURL, pkg.DocsURL())
}
