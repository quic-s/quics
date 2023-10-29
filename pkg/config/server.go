package config

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

func SetServerAddress(ip string, port string, port3 string) {
	if ip != "" {
		WriteViperEnvVariables("REST_SERVER_ADDR", ip)
	}

	if port != "" {
		WriteViperEnvVariables("REST_SERVER_PORT", port)
	}

	if port3 != "" {
		WriteViperEnvVariables("REST_SERVER_H3_PORT", port3)
	}
}
