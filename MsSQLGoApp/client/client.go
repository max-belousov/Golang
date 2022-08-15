package main

import (
	"fmt"
	"net"
	"time"
)

func main() {

	protocol, path := "tcp", "127.0.0.1:4545"
	createConnection(protocol, path)
}

func createConnection(protocol, path string) {
	conn, err := net.Dial(protocol, path)
	if err != nil {
		fmt.Println(err)
		return
	}
	handleUpdates(conn)
}

func handleUpdates(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Введите комманду:\n1 - Добавить клиента в базу\n2 - Удалить клиента из базы\n" +
		"3 - Показать всех клиентов\nhelp для вывода этой подсказки")
	for {
		var source string
		_, err := fmt.Scanln(&source)
		if err != nil {
			fmt.Println("Некорректный ввод", err)
			continue
		}

		toServer(conn, source)

		fromServer(conn)

	}
}

func toServer(conn net.Conn, source string) {
	if n, err := conn.Write([]byte(source)); n == 0 || err != nil {
		fmt.Println(err)
		return
	}
}
func fromServer(conn net.Conn) {
	fmt.Println("Ответ от сервера")
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	for {
		buff := make([]byte, 1024*4)
		n, err := conn.Read(buff)
		if err != nil {
			return
		}
		fmt.Print(string(buff[0:n]))
		conn.SetReadDeadline(time.Now().Add(time.Millisecond * 700))
		fmt.Println()
	}
}
