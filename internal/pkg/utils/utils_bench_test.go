package utils

import (
	"os"
	"testing"
)

var testFileContent = ` Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed nisi lorem, elementum at odio et, aliquet consectetur enim. Ut laoreet mattis nunc id convallis. Proin id erat tempor, faucibus enim sed, vehicula nulla. Mauris scelerisque dui in turpis ultricies, ut maximus ligula pellentesque. Sed turpis nisl, finibus sit amet scelerisque sit amet, suscipit eu ante. Nunc varius malesuada risus consequat fringilla. Etiam neque mi, congue ut mauris ac, gravida varius nulla. Etiam bibendum massa est, ac facilisis tortor dictum nec. Sed ut rhoncus enim. Sed ac malesuada augue. Quisque nec tempor turpis. In hac habitasse platea dictumst. Aliquam elementum lacus tincidunt vehicula pretium.

Cras finibus purus vitae massa tristique, in feugiat erat pulvinar. Nunc neque nisl, semper vitae ipsum sed, ullamcorper sodales tortor. Maecenas turpis risus, congue eu dignissim in, consectetur convallis neque. Nullam ultrices metus nunc, non sollicitudin enim rutrum a. Ut et accumsan augue, eu blandit tortor. Etiam fringilla elit vitae lacus elementum dignissim. Vivamus vestibulum nulla id sem eleifend euismod. Donec eget consequat tortor.

Aliquam aliquam metus ac metus mattis pulvinar. Fusce pharetra vehicula tincidunt. Morbi urna urna, dictum nec vehicula vel, efficitur id ipsum. Fusce cursus, lacus sit amet viverra condimentum, nunc metus venenatis nunc, eu fermentum odio ex nec lacus. Etiam sed quam nibh. Vestibulum facilisis vehicula fermentum. Nullam vitae ipsum condimentum, euismod augue non, consectetur lacus. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Aliquam at consequat arcu, quis lacinia mauris. Integer at fringilla lorem. Quisque id odio mi. In quis dolor finibus, finibus tellus vitae, gravida nisl. Sed urna felis, rhoncus id convallis vel, vehicula non turpis. Aliquam efficitur leo eget leo tempus commodo. Interdum et malesuada fames ac ante ipsum primis in faucibus. In cursus ultrices enim sit amet facilisis. `

func BenchmarkCopyFileContent(b *testing.B) {
	src, _ := os.CreateTemp("", "src")
	defer os.RemoveAll(src.Name())
	dst, _ := os.CreateTemp("", "dst")
	defer os.RemoveAll(dst.Name())
	src.WriteString(testFileContent)
	for i := 0; i < b.N; i++ {
		CopyFileContent(src.Name(), dst.Name())
	}
}
