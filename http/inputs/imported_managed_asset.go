package inputs

import (
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/smartystreets/detour"
)

type ImportManagedAsset struct {
	Name     string
	MIMEType string
	Size     uint64
	Reader   io.ReadCloser
}

func (this *ImportManagedAsset) Bind(request *http.Request) error {
	reader, header, err := request.FormFile(fileField)
	if err != nil {
		return badAssetError
	}

	this.Name = normalizeFilename(header.Filename)
	this.MIMEType = this.computeMIMEType()
	this.Size = uint64(header.Size)
	this.Reader = reader

	return nil
}

func (this *ImportManagedAsset) Sanitize() {
	this.Name = strings.TrimSpace(this.Name)
}

func (this *ImportManagedAsset) Validate() error {
	var errors detour.Errors
	errors = errors.AppendIf(filenameError, len(this.Name) == 0)
	errors = errors.AppendIf(emptyAssetError, this.Size == 0)
	errors = errors.AppendIf(unsupportedTypeError, len(this.MIMEType) == 0)
	return errors
}

func (this *ImportManagedAsset) computeMIMEType() string {
	switch path.Ext(this.Name) {
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".ods":
		return "application/vnd.oasis.opendocument.spreadsheet"
	case ".odt":
		return "application/vnd.oasis.opendocument.text"
	default:
		return ""
	}
}

func normalizeFilename(value string) string {
	value = strings.TrimSpace(value)
	base := path.Base(value)
	switch base {
	case ".":
		return ""
	case "/":
		return ""
	default:
		return base + strings.ToLower(path.Ext(base))
	}
}

var (
	badAssetError        = fieldError("Unable to receive the managed asset.", fileField)
	emptyAssetError      = fieldError("The manged asset body can't be empty.", fileField)
	filenameError        = fieldError("The managed asset filename can't be empty.", fileField)
	unsupportedTypeError = fieldError("The type of managed asset supplied is unsupported.", fileField)
	DuplicateAssetResult = conflictResult("The asset provided already exists.", fileField)
)
