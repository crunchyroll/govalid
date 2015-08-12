// Sample .v file for the documentation.

package worker

type jobInput struct {
	// Specifying just the type means the field is required.
	jobId uint

	// Specifying a tag "valid", with the item "def" means that the
	// field is optional, and that the default value should be the
	// zero value of the type.
	nodeBlock bool `valid:"def"`

	// You can specify maximum lengths for strings.
	encryptionKey string `valid:"max:128"`
	encryptionIv  string `valid:"max:8"`

	// Maximum lengths also work for URLs.
	destination *url.URL `valid:"max:256"`

	// You can also mandate minimum lengths and set explicit default
	// values.
	language string `valid:"min:2,max:4,def:\"enUS\""`

	// Bounds apply to numeric types as well.
	threads uint `valid:"max:8"`

	// Fields in the input data not mentioned in the struct will not
	// appear in the validated output.
}
