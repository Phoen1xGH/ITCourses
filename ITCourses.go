package main

import (
	MODMLOG "ITCourses/mlog"   // присоединение модуля mlog
	MODMPAGE "ITCourses/mpage" // присоединение модуля mpage
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

var host, user, password string

// объявим переменную для создания объекта БД
var dbmysqlg *sql.DB
var err error

// функция init() предназначена для инициализации переменных
// данной сессии
func init() {
	// ---------------------------------------------------------
	MODMLOG.LoggedUserSession.Options = &sessions.Options{
		// необходимо задать адрес вашего домена (может быть и localhost)
		// путь Path, максимальное время сессии MaxAge
		Domain: "localhost",
		Path:   "/",
		MaxAge: 3600 * 3, // 3 часа
		// задаем флаг HttpOnly: true, чтобы код на javaScript не имел доступа к
		// cookies
		HttpOnly: true,
	}
	// ---------------------------------------------------------
}

/* ****************************************************** */
func main() {
	// ---------------------------------------------------------
	// соединяемся с СУБД
	// открываем файл db.txt, он лежит в той же директории, что и activity.go
	// получаем адрес БД, логин и пароль для подключения
	host, user, password = MODMLOG.InfoMyConn("db.txt")
	// выведем на печать считанную информацию
	fmt.Println(host, user, password)
	// создаем соединение с MySQL, при помощи функции Open
	dbmysqlg, err = sql.Open("mysql",
		user+":"+password+"@tcp("+host+":3306)/IT courses?charset=utf8")
	//проверка ошибки, вызов соответствующей функции
	MODMLOG.CheckErr(err, "Не могу открыть БД IT courses")
	// убедимся, что создано соединение нужного типа
	// выведем на печать тип переменной dbmysqlg (тип *sql.DB)
	fmt.Printf("dbmysqlg=%T\n", dbmysqlg)
	// отложили закрытие БД, чтобы иметь возможность ее использования
	// в процессе работы сайта
	defer dbmysqlg.Close()
	// присваиваем переменной из модуля MODMLOG переменную
	// для работы с базой данных
	MODMLOG.Exdbmysqlg = dbmysqlg
	// откладываем закрытие соединения с базой данных
	defer MODMLOG.Exdbmysqlg.Close()
	// присваиваем переменной из модуля MODMPAGE переменную
	// для работы с базой данных

	MODMPAGE.Exdbmysqlg = dbmysqlg

	// откладываем закрытие соединения с базой данных

	defer MODMPAGE.Exdbmysqlg.Close()

	//
	// ---------------------------------------------------------
	// функция Handle служит для организации доступа к директориям
	// StripPrefix() возвращает обработчик, который обслуживает HTTP-запросы
	// функция FileServer() возвращает обработчик, который используется для
	// доступа к статическим файлам (например, CSS) из указанной директории
	// по HTTP протоколу
	http.Handle("/public/", http.StripPrefix("/public/",
		http.FileServer(http.Dir("public"))))
	// каждому маршруту ставим в соответствие функцию, ответственную за
	// обработку запроса по этому маршруту
	http.HandleFunc("/", MODMLOG.LoginPageHandler)
	http.HandleFunc("/index", MODMLOG.Index)
	http.HandleFunc("/logout", MODMLOG.LogoutHandler)
	http.HandleFunc("/searchstudent", MODMPAGE.SearchStudent)
	http.HandleFunc("/searchconference", MODMPAGE.SearchСonference)
	http.HandleFunc("/cityclassifier", MODMPAGE.CityClassifier)
	http.HandleFunc("/searchproject", MODMPAGE.SearchProject)
	http.HandleFunc("/searchpaper", MODMPAGE.SearchPaper)
	http.HandleFunc("/chartbasicbar", MODMPAGE.ChartBasicBar)
	http.HandleFunc("/searchreport", MODMPAGE.SearchReport)

	// вывод на печать результатов соединения с сервером MySQL
	fmt.Printf("Starting server IT courses...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
