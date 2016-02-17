package gosyncmodules

func Initialrun(ADHost, AD_Port, ADUsername, ADPassword, ADBaseDN, ADFilter string, ADAttribute []string, ADPage int, ADConnTimeout int, shutdownChannel chan bool)  {
	connect := ConnectToAD(ADHost, AD_Port, ADUsername, ADPassword, ADConnTimeout)
	defer func() {shutdownChannel <- true}()
	defer Info.Println("closed")
	defer connect.Close()
	defer Info.Println("Closing connection")
	ADElements := GetFromAD(connect, ADBaseDN, ADFilter, ADAttribute, uint32(ADPage))
	Info.Println(ADElements)

}
