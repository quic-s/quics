package config

import "errors"

func GetRestServerAddress() string {
	serverIP := GetViperEnvVariables("REST_SERVER_ADDR") + ":"
	serverPort := GetViperEnvVariables("REST_SERVER_PORT")

	return serverIP + serverPort
}

func GetRestServerH3Address() string {
	serverIP := GetViperEnvVariables("REST_SERVER_ADDR") + ":"
	serverPort := GetViperEnvVariables("REST_SERVER_H3_PORT")

	return serverIP + serverPort
}

func SetServerAddress(ip string, port string, port3 string) error {
	if ip != "" {
		err := WriteViperEnvVariables("REST_SERVER_ADDR", ip)
		if err != nil {
			err = errors.New("while setting server address: " + err.Error())
			return err
		}
	}

	if port != "" {
		err := WriteViperEnvVariables("REST_SERVER_PORT", port)
		if err != nil {
			err = errors.New("while setting server address: " + err.Error())
			return err
		}
	}

	if port3 != "" {
		err := WriteViperEnvVariables("REST_SERVER_H3_PORT", port3)
		if err != nil {
			err = errors.New("while setting server address: " + err.Error())
			return err
		}
	}
	return nil
}
