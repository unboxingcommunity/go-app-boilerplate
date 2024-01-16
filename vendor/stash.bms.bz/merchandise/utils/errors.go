package utils

//Error ...
type Error struct {
	JSONMarshalling           string
	JSONUnmarshalling         string
	MapstructureDecode        string
	MapstructureDecoderConfig string
	MapstructureTagMissing    string
	MapstructureResultMissing string
	Encryption                string
	Decryption                string
	MultiPiped                string
	DataConversion            string
	UnsupportedDataType       string
	YAMLMarshalling           string
}

//Errors ...
var Errors = Error{
	"UTILS.BIND.JSON_MARSHALLING",
	"UTILS.BIND.JSON_UNMARSHALLING",
	"UTILS.BIND.MAPSTRUCTURE_DECODE",
	"UTILS.BIND.MAPSTRUCTURE_DECODER_CONFIG",
	"UTILS.BIND.MAPSTRUCTURE_TAG_MISSING",
	"UTILS.BIND.MAPSTRUCTURE_RESULT_MISSING",
	"UTILS.ENCRYPT.ENCRYPTION",
	"UTILS.DECRYPT.DECRYPTION",
	"UTILS.CONFIG.BIND",
	"UTILS.CONFIG.DATA_CONVERSION",
	"UTILS.CONFIG.UNSUPPORTED_DATA_TYPE",
	"UTILS.BIND.YAML_MARSHALLING",
}

var errorConfig = map[string]interface{}{
	"UTILS.BIND.JSON_MARSHALLING": map[string]interface{}{
		"code":        "UTILS.BIND.JSON_MARSHALLING_ERROR",
		"description": "Utils error",
		"message":     "Oops! Something went wrong",
		"severity":    "Level 1",
		"action":      "refresh",
		"statusCode":  400,
	},
	"UTILS.BIND.JSON_UNMARSHALLING": map[string]interface{}{
		"code":        "UTILS.BIND.JSON_UNMARSHALLING_ERROR",
		"description": "Utils error",
		"message":     "Oops! Something went wrong",
		"severity":    "Level 1",
		"action":      "refresh",
		"statusCode":  400,
	},
	"UTILS.BIND.MAPSTRUCTURE_DECODE": map[string]interface{}{
		"code":        "UTILS.BIND.MAPSTRUCTURE_DECODE_ERROR",
		"description": "Utils error",
		"message":     "Oops! Something went wrong",
		"severity":    "Level 1",
		"action":      "refresh",
		"statusCode":  400,
	},
	"UTILS.BIND.MAPSTRUCTURE_DECODER_CONFIG": map[string]interface{}{
		"code":        "UTILS.BIND.MAPSTRUCTURE_DECODER_CONFIG_ERROR",
		"description": "Utils error",
		"message":     "Oops! Something went wrong",
		"severity":    "Level 1",
		"action":      "refresh",
		"statusCode":  400,
	},
	"UTILS.BIND.MAPSTRUCTURE_TAG_MISSING": map[string]interface{}{
		"code":        "UTILS.BIND.MAPSTRUCTURE_TAG_MISSING_ERROR",
		"description": "Utils error",
		"message":     "Oops! Something went wrong",
		"severity":    "Level 1",
		"action":      "refresh",
		"statusCode":  400,
	},
	"UTILS.BIND.MAPSTRUCTURE_RESULT_MISSING": map[string]interface{}{
		"code":        "UTILS.BIND.MAPSTRUCTURE_RESULT_MISSING_ERROR",
		"description": "Utils error",
		"message":     "Oops! Something went wrong",
		"severity":    "Level 1",
		"action":      "refresh",
		"statusCode":  400,
	},
	"UTILS.ENCRYPT.ENCRYPTION": map[string]interface{}{
		"code":        "UTILS.ENCRYPT.ENCRYPTION_ERROR",
		"description": "Utils error",
		"message":     "Oops! Something went wrong",
		"severity":    "Level 1",
		"action":      "refresh",
		"statusCode":  400,
	},
	"UTILS.DECRYPT.DECRYPTION": map[string]interface{}{
		"code":        "UTILS.DECRYPT.DECRYPTION_ERROR",
		"description": "Utils error",
		"message":     "Oops! Something went wrong",
		"severity":    "Level 1",
		"action":      "refresh",
		"statusCode":  400,
	},
	"UTILS.CONFIG.BIND": map[string]interface{}{
		"code":        "UTILS.CONFIG.BIND_ERROR",
		"description": "Utils error : multiple pipes existing in config value",
		"message":     "Oops! Something went wrong",
		"severity":    "priority 1",
		"action":      "refresh",
		"statusCode":  400,
	},
	"UTILS.CONFIG.DATA_CONVERSION": map[string]interface{}{
		"code":        "UTILS.CONFIG.DATA_CONVERSION_ERROR",
		"description": "Utils error : error converting data from one type to another",
		"message":     "Oops! Something went wrong",
		"severity":    "priority 1",
		"action":      "refresh",
		"statusCode":  400,
	},
	"UTILS.CONFIG.UNSUPPORTED_DATA_TYPE": map[string]interface{}{
		"code":        "UTILS.CONFIG.UNSUPPORTED_DATA_TYPE_ERROR",
		"description": "Utils error : Data type provided is not supported for conversion",
		"message":     "Oops! Something went wrong",
		"severity":    "priority 1",
		"action":      "refresh",
		"statusCode":  400,
	},
	"UTILS.BIND.YAML_MARSHALLING": map[string]interface{}{
		"code":        "UTILS.BIND.YAML_MARSHALLING_ERROR",
		"description": "Utils error",
		"message":     "Oops! Something went wrong",
		"severity":    "Level 1",
		"action":      "refresh",
		"statusCode":  400,
	},
}

// GetErrors ...
func GetErrors() (errMap map[string]interface{}) {
	return errorConfig
}

// GetError ...
func (er Error) GetError(code string) (err map[string]interface{}) {
	if errorConfig[code] == nil {
		return nil
	}
	return errorConfig[code].(map[string]interface{})
}
