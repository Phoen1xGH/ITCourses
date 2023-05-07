package mpage

import (
	MODMLOG "ITCourses/mlog"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
	_ "strconv"
)

// переменная которой было в пакете main файла activity.go присвоено
// значение, необходимое для работы с базой данных
var Exdbmysqlg *sql.DB

/* ****************************************************** */
// функция для поиска и выбора студента
// выбор студента осуществляется двойным щелчком мыши по записи
func SearchStudent(w http.ResponseWriter, r *http.Request) {
	// ---------------------------------------------------------
	// проверяем, что у нас не запрос GET, а запрос POST
	// в данной программе ограничиваемся этими двумя запросами
	sOut := ""
	if r.Method == "GET" {
		// если запрос был искусственно отправлен методом GET, то
		// возвращаем сообщение об этом
		// из программы приходит только запрос POST
		fmt.Fprintf(w, "%v", sOut+" GET ")
		return
	} else {
		// ---------------------------------------------------------
		// проверяем, если есть проблемы с входом в сессию методом POST,
		// если есть, возвращаем сообщение об ошибке
		if MODMLOG.CheckLoginPOST(w, r) == 0 {
			fmt.Fprintf(w, "%v", "0####/")
			return
		}
		// ---------------------------------------------------------
		// считываем данные формы
		// функция ParseMultipartForm пакета net/http анализирует тело запроса
		// как multipart/form-data
		// multipart/form-data чаще всего используется для отправки HTML-форм с
		// бинарными данными методом POST протокола HTTP
		if err := r.ParseMultipartForm(64 << 20); err != nil {
			fmt.Println("ParseForm() err: ", err)
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			MODMLOG.CheckErr(err, "Ошибка запроса POST: SearchCourse")
		}
		// ---------------------------------------------------------
		// считываем данные из текстового поля
		searchstr := r.FormValue("searchstr")
		// ---------------------------------------------------------
		// обращаемся к базе данных, готовим текст запроса
		stmt, err := Exdbmysqlg.Prepare("SELECT `IT courses`.DevelopmentCategory.id as courseId, `IT courses`.DevelopmentCategory.CategoryName as devCat, " +
			"`IT courses`.CourseLanguage.ProgrammingLanguage as progLang, `IT courses`.CourseCost.Cost as cost " +
			"FROM  `IT courses`.DevelopmentCategory, `IT courses`.CourseLanguage, `IT courses`.CourseCost " +
			"WHERE `IT courses`.DevelopmentCategory.id = `IT courses`.CourseLanguage.idCategory " +
			"AND `IT courses`.DevelopmentCategory.id = `IT courses`.CourseCost.id " +
			"AND `IT courses`.DevelopmentCategory.CategoryName LIKE ? ;")
		// проверяем, что удалось сформировать запрос
		MODMLOG.CheckErr(err, "Не могу подготовить запрос к БД IT courses")
		defer stmt.Close()
		// выполняем запрос
		rows, err := stmt.Query(searchstr + "%")

		// проверяем, что удалось выполнить запрос
		MODMLOG.CheckErr(err, "Не могу выполнить запрос в БД IT courses")
		defer rows.Close()
		// объявляем переменные соответствующего типа для считывания полей
		var courseId sql.NullInt64
		var devCat sql.NullString
		var progLang sql.NullString
		var cost sql.NullInt64
		// цикл по всем возвращенным строкам
		for rows.Next() {
			// считывание строки
			err = rows.Scan(&courseId, &devCat, &progLang, &cost)
			// проверяем, что записи считались
			MODMLOG.CheckErr(err, "Не могу прочесть запись")
			// формируем из считанных полей БД строку таблицы HTML
			sOut += "<tr><td>" + devCat.String + "</td><td>&nbsp;&nbsp;</td>" +
				"<td>" + progLang.String + "</td><td>&nbsp;&nbsp;</td>" +
				"<td>" + strconv.FormatInt(cost.Int64, 10) + "</td></tr>\n"
		}
	} // POST
	// формируем заголовок таблицы с использованием Bootstrap
	sOut = "<table class=\"table table-striped\">\n" +
		" <thead>\n" +
		" <tr>\n" +
		" <th scope=\"col\">Категория разработки</th>\n" +
		" <th scope=\"col\">&nbsp;</th>\n" +
		" <th scope=\"col\">Язык программирования</th>\n" +
		" <th scope=\"col\">&nbsp;</th>\n" +
		" <th scope=\"col\">Стоимость</th>\n" +
		" </tr>\n" +
		" </thead>\n" +
		" <tbody>\n" +
		sOut +
		" </tbody>" +
		"</table>\n"
	fmt.Fprintf(w, "%v", sOut)
}
