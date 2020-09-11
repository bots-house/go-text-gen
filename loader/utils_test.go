package loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTemplateVars(t *testing.T) {
	for _, test := range []struct {
		Text string
		Vars []string
	}{
		{"", []string{}},
		{"{var}", []string{"var"}},
		{"{x} {y} {z}", []string{"x", "y", "z"}},
	} {
		assert.Equal(t,
			test.Vars,
			getTemplateVars(test.Text),
		)
	}
}
