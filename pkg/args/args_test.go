package args

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Parallel()

	t.Run("Minimal", func(t *testing.T) {
		t.Parallel()

		arguments := Parse([]string{"dtsync", "-src", "src", "-dst", "dst"})
		assert.Equal(t, Arguments{
			SrcRootPath:             "src",
			DstRootPath:             "dst",
			ReplaceNotMatchingFiles: false,
			RemoveDstLeftover:       false,
		}, arguments)
	})

	t.Run("Replace", func(t *testing.T) {
		t.Parallel()

		arguments := Parse([]string{"dtsync", "-src", "src", "-dst", "dst", "-replace"})
		assert.Equal(t, Arguments{
			SrcRootPath:             "src",
			DstRootPath:             "dst",
			ReplaceNotMatchingFiles: true,
			RemoveDstLeftover:       false,
		}, arguments)
	})

	t.Run("Remove", func(t *testing.T) {
		t.Parallel()

		arguments := Parse([]string{"dtsync", "-src", "src", "-dst", "dst", "-remove"})
		assert.Equal(t, Arguments{
			SrcRootPath:             "src",
			DstRootPath:             "dst",
			ReplaceNotMatchingFiles: false,
			RemoveDstLeftover:       true,
		}, arguments)
	})

	t.Run("ReplaceAndRemove", func(t *testing.T) {
		t.Parallel()

		arguments := Parse([]string{"dtsync", "-src", "src", "-dst", "dst", "-replace", "-remove"})
		assert.Equal(t, Arguments{
			SrcRootPath:             "src",
			DstRootPath:             "dst",
			ReplaceNotMatchingFiles: true,
			RemoveDstLeftover:       true,
		}, arguments)
	})
}
