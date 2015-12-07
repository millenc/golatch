package golatch

const (
	//Latch related constants
	API_URL                                  = "https://latch.elevenpaths.com"
	API_PATH                                 = "/api"
	API_VERSION                              = "1.0"
	API_CHECK_STATUS_ACTION                  = "status"
	API_PAIR_ACTION                          = "pair"
	API_PAIR_WITH_ID_ACTION                  = "pairWithId"
	API_UNPAIR_ACTION                        = "unpair"
	API_LOCK_ACTION                          = "lock"
	API_UNLOCK_ACTION                        = "unlock"
	API_HISTORY_ACTION                       = "history"
	API_OPERATION_ACTION                     = "operation"
	API_APPLICATION_ACTION                   = "application"
	API_SUBSCRIPTION_ACTION                  = "subscription"
	API_NOOTP_SUFFIX                         = "nootp"
	API_SILENT_SUFFIX                        = "silent"
	API_AUTHENTICATION_METHOD                = "11PATHS"
	API_AUTHORIZATION_HEADER_NAME            = "Authorization"
	API_DATE_HEADER_NAME                     = "X-11Paths-Date"
	API_AUTHORIZATION_HEADER_FIELD_SEPARATOR = " "
	API_X_11PATHS_HEADER_PREFIX              = "X-11Paths-"
	API_X_11PATHS_HEADER_SEPARATOR           = ":"
	API_UTC_STRING_FORMAT                    = "2006-01-02 15:04:05" //format layout as defined here: http://golang.org/pkg/time/#pkg-constants

	//Possible values for the Two factor and Lock on request options
	//NOT_SET is used in the UpdateOperation() method to leave the existing value
	MANDATORY = "MANDATORY"
	OPT_IN    = "OPT_IN"
	DISABLED  = "DISABLED"
	NOT_SET   = ""

	//Possible status values for the latch
	LATCH_STATUS_ON  = "on"
	LATCH_STATUS_OFF = "off"

	//HTTP methods
	HTTP_METHOD_POST   = "POST"
	HTTP_METHOD_GET    = "GET"
	HTTP_METHOD_PUT    = "PUT"
	HTTP_METHOD_DELETE = "DELETE"
)
