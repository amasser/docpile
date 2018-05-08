package inputs

import (
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/smartystreets/detour"
)

type AddFile struct {
	Filename      string
	Reader        io.ReadCloser
	MIMEType      string
	contentLength uint64
}

func (this *AddFile) Bind(request *http.Request) error {
	reader, file, err := request.FormFile(fileField)
	if err != nil {
		return badAssetError
	}

	this.Filename = file.Filename
	this.MIMEType = this.computeMIMEType()
	this.Reader = reader
	this.contentLength = uint64(request.ContentLength)

	return nil
}

func (this *AddFile) Sanitize() {
	this.Filename = strings.TrimSpace(this.Filename)
}

func (this *AddFile) Validate() error {
	var errors detour.Errors
	errors = errors.AppendIf(filenameError, len(this.Filename) == 0)
	errors = errors.AppendIf(emptyAssetError, this.contentLength == 0)
	errors = errors.AppendIf(unsupportedTypeError, len(this.MIMEType) == 0)
	return errors
}

func (this *AddFile) computeMIMEType() string {
	switch path.Ext(strings.ToLower(this.Filename)) {
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

var (
	badAssetError        = fieldError("Unable to receive the managed asset.", fileField)
	emptyAssetError      = fieldError("The manged asset body can't be empty.", fileField)
	filenameError        = fieldError("The managed asset filename can't be empty.", fileField)
	unsupportedTypeError = fieldError("The type of managed asset supplied is unsupported.", fileField)
	DuplicateAssetResult = conflictResult("The asset provided already exists.", fileField)
)
