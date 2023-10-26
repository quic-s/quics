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

func GetHttp3ServerAddress(ip string, port string) string {
	serverIP := GetViperEnvVariables("REST_SERVER_ADDR") + ":"
	if ip != "" {
		serverIP = ""
		serverIP = ip + ":"
		WriteViperEnvVariables("REST_SERVER_ADDR", ip)
	}

	serverPort := GetViperEnvVariables("REST_SERVER_H3_PORT")
	if port != "" {
		serverPort = ""
		serverPort = port
		WriteViperEnvVariables("REST_SERVER_H3_PORT", port)
	}

	return serverIP + serverPort
}
