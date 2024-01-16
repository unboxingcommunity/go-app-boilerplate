package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
	errors "stash.bms.bz/merchandise/go-errors"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/xid"
)

// Invoke function that is passed externally to BindCcms
type Invoke func(string) (string, error)

// GetID gets unique xid
func GetID(prefix ...string) (id string) {
	id = xid.New().String()
	id = strings.ToUpper(id)
	id = id[0:5] + "-" + id[5:10] + "-" + id[10:15] + "-" + id[15:20]

	// If prefix is empty
	if len(prefix) == 0 {
		return id
	}

	// Returns new unique id
	return prefix[0] + "-" + id
}

// Bind binds data
func Bind(input interface{}, output interface{}, tagName ...string) (err error) {
	if len(tagName) == 0 {
		// Decodes
		err = mapstructure.Decode(input, output)
		if err != nil {
			return errors.Wrap(Errors.MapstructureDecode, err)
		}
	} else {
		if tagName[0] == "" {
			return errors.New(Errors.MapstructureTagMissing)
		}
		config := &mapstructure.DecoderConfig{TagName: tagName[0], Result: output}
		ms, err := mapstructure.NewDecoder(config)
		if err != nil {
			return errors.Wrap(Errors.MapstructureDecoderConfig, err)
		}

		err = ms.Decode(input)
		if err != nil {
			return errors.Wrap(Errors.MapstructureDecode, err)
		}

	}

	// Returns
	return nil
}

// ReadFile ...
func ReadFile(path string, instance interface{}) (err error) {
	// Forms bytes
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// Unmarshals
	err = json.Unmarshal(bytes, instance)
	if err != nil {
		return err
	}

	// Returns
	return nil
}

func pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func unPad(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

// Encrypt ...
func Encrypt(key string, iv string, plaintext string) (ciphertext string, err error) {
	keyBytes := []byte(key)
	ivBytes := []byte(iv)
	plaintextBytes := []byte(plaintext)
	plaintextBytes = pad(plaintextBytes, 16)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", errors.Wrap(Errors.Encryption, err)
	}
	encrypter := cipher.NewCBCEncrypter(block, ivBytes)

	ciphertextBytes := make([]byte, len(plaintextBytes))
	encrypter.CryptBlocks(ciphertextBytes, plaintextBytes)
	ciphertext = base64.StdEncoding.EncodeToString(ciphertextBytes)

	return ciphertext, nil
}

// Decrypt ...
func Decrypt(key string, iv string, ciphertext string) (plaintext string, err error) {
	keyBytes := []byte(key)
	ivBytes := []byte(iv)
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", errors.Wrap(Errors.Decryption, err)
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", errors.Wrap(Errors.Decryption, err)
	}
	decrypter := cipher.NewCBCDecrypter(block, ivBytes)

	plaintextBytes := make([]byte, len(ciphertextBytes))
	decrypter.CryptBlocks(plaintextBytes, ciphertextBytes)
	plaintextBytes = unPad(plaintextBytes)
	plaintext = string(plaintextBytes)

	return plaintext, nil
}

// BindConfig evalutes struct and assigns ccms values if required
func BindConfig(input []byte, output interface{}, tagName string, fn Invoke) (err error) {
	var inputMap map[interface{}]interface{}

	// Unmarshalls config into map
	err = yaml.Unmarshal(input, &inputMap)
	if err != nil {
		return errors.Wrap(Errors.YAMLMarshalling, err)
	}

	// Performs the famous moonwalk to make it look like magic
	err = moonwalk(inputMap, tagName, fn)
	if err != nil {
		return err
	}

	// Binds result back to config struct
	err = mapstructure.Decode(inputMap, &output)
	if err != nil {
		return errors.Wrap(Errors.MapstructureDecode, err)
	}

	return nil
}

// Recursive walk to iterate through
func moonwalk(m map[interface{}]interface{}, tagName string, fn Invoke) (err error) {
	for k, v := range m {
		// If there is no value set , by default set it to an empty string
		if v == nil {
			m[k] = ""
		} else if reflect.TypeOf(v).Kind() == reflect.String {
			value := v.(string)
			newValue, err := getCcmsValue(value, tagName, fn)
			if err != nil {
				return err
			}
			// if new value returned is same as key the  no ccms exists .. assign default value
			if newValue == value {
				m[k] = value
			} else {
				m[k] = newValue
			}
		} else if reflect.TypeOf(v).Kind() == reflect.Map {
			if v != nil {
				err = moonwalk(v.(map[interface{}]interface{}), tagName, fn)
				if err != nil {
					return err
				}
			}
			m[k] = v
		} else if reflect.TypeOf(v).Kind() == reflect.Slice {
			s := reflect.ValueOf(v)
			if reflect.TypeOf(s.Index(0).Interface()).Kind() == reflect.Map {
				var arrayMap []map[interface{}]interface{}
				for i := 0; i < s.Len(); i++ {
					innerMap := s.Index(i).Interface().(map[interface{}]interface{})
					err = moonwalk(innerMap, tagName, fn)
					if err != nil {
						return err
					}
					arrayMap = append(arrayMap, innerMap)
				}
				m[k] = arrayMap
			} else if reflect.TypeOf(s.Index(0).Interface()).Kind() == reflect.String {
				var arrayStr []interface{}
				for i := 0; i < s.Len(); i++ {
					innerStr := s.Index(i).Interface().(string)
					newStr, err := getCcmsValue(innerStr, tagName, fn)
					if err != nil {
						return err
					}
					arrayStr = append(arrayStr, newStr)
				}
				m[k] = arrayStr
			} else {
				m[k] = v
			}
		} else {
			m[k] = v
		}
	}
	return nil
}

// getCcmsValue , gets appropriate value for the key
func getCcmsValue(key string, tagName string, fn Invoke) (value interface{}, err error) {
	keyArr := strings.Split(key, "|")
	keyArrLen := len(keyArr)
	if keyArrLen > 3 {
		return value, errors.New(Errors.MultiPiped)
	}
	if strings.ToLower(keyArr[0]) == tagName {
		if keyArrLen == 2 {
			value, err = fn(keyArr[1])
			if err != nil {
				return value, err
			}
		} else if keyArrLen == 3 {
			val, err := fn(keyArr[2])
			if err != nil {
				return value, err
			}

			// Gets appropriate value type
			value, err = getTypedValue(val, keyArr[1])
			if err != nil {
				return value, err
			}
		}
	} else {
		value = key
	}

	return value, err
}

// getTypedValue , gets type casted appropriate value type
func getTypedValue(value string, vType string) (result interface{}, err error) {
	vType = strings.ToLower(vType)

	switch vType {
	case "string":
		result = value
		return result, err
	case "bool":
		result, err = strconv.ParseBool(value)
		if err != nil {
			return result, errors.Wrap(Errors.DataConversion, err)
		}
	case "int":
		result, err = strconv.Atoi(value)
		if err != nil {
			return result, errors.Wrap(Errors.DataConversion, err)
		}
	case "int32":
		cValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return result, errors.Wrap(Errors.DataConversion, err)
		}
		result = int32(cValue)
	case "int64":
		result, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return result, errors.Wrap(Errors.DataConversion, err)
		}
	case "float32":
		cValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return result, errors.Wrap(Errors.DataConversion, err)
		}
		result = float32(cValue)
	case "float64":
		result, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return result, errors.Wrap(Errors.DataConversion, err)
		}
	default:
		return result, errors.New(Errors.UnsupportedDataType)
	}

	// Returns
	return result, err
}
