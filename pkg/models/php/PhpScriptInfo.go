package pm

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
}
