package ccms

import "os"

func getCCMSType() string {
	typeCCMS := SIDECAR
	envCCMSType := os.Getenv("CCMS_CALL_TYPE")
	if envCCMSType != "" {
		typeCCMS = envCCMSType
	}

	return typeCCMS
}

func getCCMSGRPCPORT() string {
	PORTCCMS := "25001"
	envCCMSPORT := os.Getenv("CCMS_GRPC_PORT")
	if envCCMSPORT != "" {
		PORTCCMS = envCCMSPORT
	}

	return PORTCCMS
}

func getCCMSGRPCHOST() string {
	HOSTCCMS := "localhost"
	envCCMSHOST := os.Getenv("CCMS_GRPC_HOST")
	if envCCMSHOST != "" {
		HOSTCCMS = envCCMSHOST
	}

	return HOSTCCMS
}
