package domain

import "errors"

var (
	TagAlreadyExistsError   = errors.New("tag already exists")
	AssetAlreadyExistsError = errors.New("asset already exists")
	AssetNotFoundError      = errors.New("asset not found")
	TagNotFoundError        = errors.New("tag not found")
	DocumentNotFoundError   = errors.New("document not found")
	StoreAssetError         = errors.New("unable to storage asset")
)
