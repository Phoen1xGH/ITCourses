package mpage

import (
	MODMLOG "ITCourses/mlog"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
	"strings"
)

/* ****************************************************** */
// функция для поиска и выбора конференции
func SearchСonference(w http.ResponseWriter, r *http.Request) {
	// ---------------------------------------------------------
	// проверяем, что у нас не запрос GET, а запрос POST
	sOut := ""
	// если запрос был искусственно отправлен методом GET, то
	// возвращаем сообщение об этом
	// из программы приходит только запрос POST
	if r.Method == "GET" {
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
		if err := r.ParseMultipartForm(64 << 20); err != nil {
			fmt.Println("ParseForm() err: ", err)
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			MODMLOG.CheckErr(err, "Ошибка запроса POST: SearchСonference")
		}
		// ---------------------------------------------------------
		// если id студента не задан, осуществляем поиск для всех студентов

		// считываем данные из текстовых полей
		confname := r.FormValue("confname")
		confcity := r.FormValue("confcity")
		//confdatestart := r.FormValue("confdatestart")
		//confdateend:= r.FormValue("confdatestart")
		confcity = strings.TrimSpace(confcity)
		if confcity == "" {
			confcity = "%"
		}

		confdatestart := r.FormValue("confdatestart")
		confdatestart = strings.TrimSpace(confdatestart)
		confdateend := r.FormValue("confdateend")
		confdateend = strings.TrimSpace(confdateend)

		//confdatestart, _ := strconv.Atoi("confdatestart")
		//confdateend, _ := strconv.Atoi("confdateend")
		// если не указана конечная дата, делаем ее максимально большой
		// даже если такой не существует
		// ---------------------------------------------------------
		// обращаемся к базе данных, готовим текст запроса
		stmt, err := Exdbmysqlg.Prepare("SELECT `IT courses`.CourseLanguage.ProgrammingLanguage as progLang, " +
			"`IT courses`.DevelopmentCategory.CategoryName as devCat, " +
			"`IT courses`.CourseCost.Cost as cost " +
			"FROM `IT courses`.CourseLanguage " +
			"LEFT JOIN `IT courses`.DevelopmentCategory " +
			"ON `IT courses`.CourseLanguage.idCategory=`IT courses`.DevelopmentCategory.id " +
			"LEFT JOIN  `IT courses`.CourseCost " +
			"ON `IT courses`.CourseCost.idCategory=`IT courses`.DevelopmentCategory.id " +
			"WHERE `IT courses`.CourseLanguage.ProgrammingLanguage LIKE ? " +
			"AND `IT courses`.DevelopmentCategory.id LIKE ? " +
			"AND `IT courses`.CourseCost.Cost > ? AND `IT courses`.CourseCost.Cost <= ? ;")
		// проверяем, что удалось сформировать запрос
		MODMLOG.CheckErr(err, "Не могу подготовить запрос к БД activity")
		defer stmt.Close()
		// выполняем запрос
		rows, err := stmt.Query(
			confname+"%",
			confcity+"%",
			confdatestart,
			confdateend)
		// проверяем, что удалось выполнить запрос
		MODMLOG.CheckErr(err, "Не могу выполнить запрос в БД activity")
		defer rows.Close()
		// объявляем переменные соответствующего типа для считывания полей
		var devCat sql.NullString
		var progLang sql.NullString
		var cost sql.NullInt64
		for rows.Next() {
			// считывание строки
			err = rows.Scan(&progLang, &devCat,
				&cost)
			// проверяем, что записи считались
			MODMLOG.CheckErr(err, "Не могу прочесть запись")
			// переводим дату в соответствующий формат
			// формируем из считанных полей БД строку таблицы HTML
			sOut += "<tr><td>" + progLang.String + "</td><td>&nbsp;&nbsp;</td>" +
				"<td>" + devCat.String + "</td><td>&nbsp;&nbsp;</td>" +
				"<td>" + strconv.FormatInt(cost.Int64, 10) + "</td></tr>\n"
		}
	} // POST
	// формируем заголовок таблицы с использованием Bootstrap
	sOut = "<table class=\"table table-striped\">\n" +
		" <thead>\n" +
		" <tr>\n" +
		" <th scope=\"col\">Язык программирования</th>\n" +
		" <th scope=\"col\">&nbsp;</th>\n" +
		" <th scope=\"col\">Категория разработки</th>\n" +
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
