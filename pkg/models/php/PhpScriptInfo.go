package pm

const (
	QueryParamExtraPath = "extrapath"
)

type PhpScriptInfo struct {
	// Paths in URL format,
	// i.e. with forward slashes.
	OriginalUrlPath string
	UrlRelPath      string
	UrlExtraPath    string

	// Paths in file system format,
	// i.e. with separators of an operating system.
	FilePath         string
	FileName         string
	FileExt          string
	FileAbsPath      string
	FileAbsExtraPath string

	// A special parameter for storing CGI extra path.
	// This parameter is used to move extra path from path to a query parameter
	// in order to make CGI requests compatible with modern HTTP standard.
	QueryParamExtraPath string
}
