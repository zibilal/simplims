package response

import (
	"github.com/zibilal/simpleapi/api"
)

const (
	unknownError      = "Terjadi kesalahan, silahkan dicoba kembali beberapa saat lagi"
	validationError   = "Data tidak sesuai"
	invalidPin        = "PIN yang and masukkan salah"
	invalidData       = "Data yang and masukkan tidak valid"
	errorMakePayment  = "Pembayaran gagal. Silahkan ulangi beberapa saat lagi"
	methodAlreadyUsed = "Metode pembayaran ini sudah digunakan. Silahkan bayar terlebih dahulu atau pilih metode pembayaran lain"
)

var (
	APIErrorUnknown         = NewVersionOneBaseResponse(api.ErrorUnknownCode, unknownError)
	APIErrorValidation      = NewVersionOneBaseResponse(api.ErrorValidationCode, validationError)
	APIErrorInvalidPIN      = NewVersionOneBaseResponse(api.ErrorInvalidPassword, invalidPin)
	APIErrorInvalidData     = NewVersionOneBaseResponse(api.ErrorInvalidData, invalidData)
	APIErrorMakePaymentCode = NewVersionOneBaseResponse(api.APIErrorPaymentFailed, errorMakePayment)
)
