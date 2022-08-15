package main

import (
	"database/sql"
	"fmt"
	"net"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Client struct {
	id         int
	FirstName  string
	SecondName string
	Phone      string
	Email      string
}

func main() {
	protocol, path := "tcp", ":4545"
	createListener(protocol, path)
}

func OpenConn() sql.DB {
	connStr := "user=postgres password=password dbname=PostgreStudDb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	return *db
}

func AddToDb(addingClient []string) {
	db := OpenConn()
	result, err := db.Exec("insert into clients (firstname, secondname, phone, email) values ($1, $2, $3, $4)",
		addingClient[0], addingClient[1], addingClient[2], addingClient[3])
	if err != nil {
		panic(err)
	}

	fmt.Println(result.RowsAffected())
}

func GetAllFromDb() (clients []Client) {
	db := OpenConn()
	rows, err := db.Query("select * from clients")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	//clients := []Client{}

	for rows.Next() {
		p := Client{}
		err := rows.Scan(&p.id, &p.FirstName, &p.SecondName, &p.Phone, &p.Email)
		if err != nil {
			fmt.Println(err)
			continue
		}
		clients = append(clients, p)
	}
	for _, p := range clients {
		fmt.Println(p.id, p.FirstName, p.SecondName, p.Phone, p.Email)
	}
	return
}

func DeleteItemById(id int) (isDeleted bool) {
	isDeleted = false
	db := OpenConn()
	result, err := db.Exec("delete from clients where id = $1", id)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())
	if n, err := result.RowsAffected(); n != 0 && err == nil {
		isDeleted = true
	}
	return
}

func createListener(protocol, path string) {
	listener, err := net.Listen(protocol, path)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		go handleConnection(conn) // запускаем горутину для обработки запроса
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		input := make([]byte, (1024 * 4))
		n, err := conn.Read(input)
		if n == 0 || err != nil {
			fmt.Println("Read error:", err)
			break
		}
		source := string(input[0:n])
		switch source {
		case "1":
			conn.Write([]byte("Введите Имя, Фамилию, Телефон и Email через запятую"))
			n, err = conn.Read(input)
			if n == 0 || err != nil {
				fmt.Println("Read error:", err)
				break
			}
			source = string(input[0:n])
			add := strings.Split(source, ",")
			fmt.Println(add[0], add[1], add[2], add[3])
			AddToDb(add)
			conn.Write([]byte("Запись добавлена"))
		case "2":
			conn.Write([]byte("Введите Id удаляемой записи"))
			n, err = conn.Read(input)
			if n == 0 || err != nil {
				fmt.Println("Read error:", err)
				break
			}
			source = string(input[0:n])
			if id, er := strconv.Atoi(source); er == nil {
				if flag := DeleteItemById(id); flag {
					conn.Write([]byte("Запись удалена"))
				} else {
					conn.Write([]byte("В базе данных нет такого индекса"))
				}
			} else {
				mes := source + " не является целым числом"
				conn.Write([]byte(mes))
			}
		case "3":
			clients := GetAllFromDb()
			var result string = "Список всех клиентов\n"
			for _, p := range clients {
				result += fmt.Sprint(p.id) + "\t" +
					p.FirstName + "\t" +
					p.SecondName + "\t" +
					p.Phone + "\t" +
					p.Email + "\n"
			}
			conn.Write([]byte(result))
		case "help":
			conn.Write([]byte("Введите комманду:\n1 - Добавить клиента в базу\n2 - Удалить клиента из базы\n" +
				"3 - Показать всех клиентов\nhelp для вывода этой подсказки"))
		default:
			conn.Write([]byte("Неизвестная команда\n"))
		}
	}
}
