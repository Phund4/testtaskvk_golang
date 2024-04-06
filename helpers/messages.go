package helpers

import (
	"net/http"

	rabbit "github.com/Phund4/testtaskvk_golang/rabbit/RabbitTest"
)

func SendMessages(w http.ResponseWriter, httpStatus int, clientMessage, serverMessage string) {
	w.WriteHeader(httpStatus)
	w.Write([]byte(clientMessage))
	rabbit.SendRabbitMessage(serverMessage)
}
