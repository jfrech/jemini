package jemini


const (
    StatusInternalErroneous         = -1
    StatusInternalNone              =  0


    StatusInput                     = 10
    StatusSensitiveInput            = 11

    StatusSuccess                   = 20

    StatusTemporaryRedirect         = 30
    StatusPermanentRedirect         = 31

    StatusTemporaryFailure          = 40
    StatusServerUnavailable         = 41
    StatusCGIError                  = 42
    StatusProxyError                = 43
    StatusSlowDown                  = 44

    StatusPermanentFailure          = 50
    StatusNotFound                  = 51
    StatusGone                      = 52
    StatusProxyRequestRefused       = 53
    StatusBadRequest                = 59

    StatusClientCertificateRequired = 60
    StatusCertificateNotAuthorised  = 61
    StatusCertificateNotValid       = 62
)

func ValidStatus(status int) bool {
    for _, valid := range []int{
        StatusInput, StatusSensitiveInput, StatusSuccess,
        StatusTemporaryRedirect, StatusPermanentRedirect,
        StatusTemporaryFailure, StatusServerUnavailable, StatusCGIError,
        StatusProxyError, StatusSlowDown, StatusPermanentFailure,
        StatusNotFound, StatusGone, StatusProxyRequestRefused, StatusBadRequest,
        StatusClientCertificateRequired, StatusCertificateNotAuthorised,
        StatusCertificateNotValid,
    } {
        if status == valid {
            return true
        }
    }
    return false
}
