package mpage

import (
	MODMLOG "ITCourses/mlog"
	"database/sql"
	"fmt"
	//"github.com/360EntSecGroup-Skylar/excelize"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nguyenthenguyen/docx"
	"github.com/xuri/excelize/v2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/* ****************************************************** */
// функция построения отчетов
func SearchReport(w http.ResponseWriter, r *http.Request) {
	// ---------------------------------------------------------
	// объявляем переменные
	format := ""
	studentid := ""
	confname := ""
	// ---------------------------------------------------------
	// обрабатываем запрос GET для получения файлов формата xlsx и docx
	if r.Method == "GET" {
		// ---------------------------------------------------------
		MODMLOG.CheckLoginGET(w, r)
		// ---------------------------------------------------------
		query := r.URL.Query()
		format = query.Get("format")
		format = strings.TrimSpace(format)
		// если id студента не задан, выполняем поиск для всех студентов
		studentid = query.Get("studentid")
		studentid = strings.TrimSpace(studentid)
		if studentid == "" {
			studentid = "%"
		}
		// если не задано название конференции, проекта или статьи
		// выполняем поиск по всей научной деятельности
		confname = query.Get("confname")
		confname = strings.TrimSpace(confname)
		if confname == "" {
			confname = "%"
		}
	} else {
		// запрос POST используется для вывода HTML отчета
		// ---------------------------------------------------------
		if MODMLOG.CheckLoginPOST(w, r) == 0 {
			fmt.Fprintf(w, "%v", "0####/")
			return
		}
		// ---------------------------------------------------------
		if err := r.ParseMultipartForm(64 << 20); err != nil {
			fmt.Println("ParseForm() err: ", err)
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			MODMLOG.CheckErr(err, "Ошибка запроса POST: SearchReport")
		}
		// ---------------------------------------------------------
		format = r.FormValue("format")
		format = strings.TrimSpace(format)
		// если id студента не задан, выполняем поиск для всех студентов
		studentid = r.FormValue("studentid")
		studentid = strings.TrimSpace(studentid)
		if studentid == "" {
			studentid = "%"
		}
		// если не задано название конференции, проекта или статьи
		// выполняем поиск по всей научной деятельности
		confname = r.FormValue("confname")
		confname = strings.TrimSpace(confname)
		if confname == "" {
			confname = "%"
		}
	}
	// ---------------------------------------------------------
	// подготовка запроса к серверу
	stmt, err := Exdbmysqlg.Prepare("SELECT `IT courses`.DevelopmentCategory.id as courseId, `IT courses`.DevelopmentCategory.CategoryName as devCat, " +
		"`IT courses`.CourseLanguage.ProgrammingLanguage as progLang, `IT courses`.CourseCost.Cost as cost " +
		"FROM  `IT courses`.DevelopmentCategory, `IT courses`.CourseLanguage, `IT courses`.CourseCost " +
		"WHERE `IT courses`.DevelopmentCategory.id = `IT courses`.CourseLanguage.idCategory " +
		"AND `IT courses`.DevelopmentCategory.id = `IT courses`.CourseCost.id " +
		"AND `IT courses`.DevelopmentCategory.CategoryName LIKE ? ;")
	MODMLOG.CheckErr(err, "Не могу подготовить запрос к БД activity")
	defer stmt.Close()
	// выполнение запроса
	rows, err := stmt.Query(confname + "%")
	MODMLOG.CheckErr(err, "Не могу выполнить запрос в БД activity")
	defer rows.Close()
	// ---------------------------------------------------------
	// в зависимости от значения переменной format вызываем различные
	// подпрограммы
	switch format {
	case "HTML":
		ReportHTML(rows, w)
	case "XLSX":
		ReportXLSX(rows, w)
	case "DOCX":
		ReportDOCX(rows, w)
	default:
		ReportHTML(rows, w)
	}
	// ---------------------------------------------------------
}

/* ****************************************************** */
// функция построения отчета в формате html
func ReportHTML(rows *sql.Rows, w http.ResponseWriter) {
	sOut := ""
	// объявляем переменные
	//nItogoConf := 0
	//nItogoStudent := 0
	// объявляем переменные соответствующего типа для считывания полей
	var courseId sql.NullInt64
	var devCat sql.NullString
	var progLang sql.NullString
	var cost sql.NullInt64
	// цикл по всем возвращенным строкам
	for rows.Next() {
		// считывание строки
		err := rows.Scan(&courseId, &devCat, &progLang, &cost)
		// проверяем, что записи считались
		MODMLOG.CheckErr(err, "Не могу прочесть запись")
		// формируем из считанных полей БД строку таблицы HTML
		sOut += "<tr><td>" + devCat.String + "</td><td>&nbsp;&nbsp;</td>" +
			"<td>" + progLang.String + "</td><td>&nbsp;&nbsp;</td>" +
			"<td>" + strconv.FormatInt(cost.Int64, 10) + "</td></tr>\n"
	}
	// POST
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

/* ****************************************************** */
// функция построения отчета в формате xlsx
func ReportXLSX(rows *sql.Rows, w http.ResponseWriter) {
	// объявляем переменные

	// объявляем переменные соответствующего типа для считывания полей
	var courseId sql.NullInt64
	var devCat sql.NullString
	var progLang sql.NullString
	var cost sql.NullInt64
	// ---------------------------------------------------------
	// задаем номер левой верхней ячейки
	// задаем корректирующие коэффициенты для осей X и Y
	nLeftCol := 1
	nUpperRow := 1
	sNameCol := ""
	dCoeffX := 1.018
	dCoeffY := 1.0
	// ---------------------------------------------------------
	xlsx := excelize.NewFile()
	// создаем новый лист
	sNameSheet := "Sheet1"
	index := xlsx.NewSheet(sNameSheet)
	// ориентация страницы
	err := xlsx.SetPageLayout(
		sNameSheet,
		excelize.PageLayoutOrientation(excelize.OrientationLandscape),
	)
	if err != nil {
		panic(err)
	}
	// задаем параметры страницы
	err = xlsx.SetPageMargins(
		sNameSheet,
		excelize.PageMarginBottom(0.2),
		excelize.PageMarginFooter(0.2),
		excelize.PageMarginHeader(0.2),
		excelize.PageMarginLeft(0.2),
		excelize.PageMarginRight(0.2),
		excelize.PageMarginTop(0.2),
	)
	if err != nil {
		panic(err)
	}
	// ---------------------------------------------------------
	// задаем границы и шрифт ячейки
	var vBorder = `[
 { "type": "left", "color": "000000", "style": 1 },
 { "type": "top", "color": "000000", "style": 1 },
 { "type": "bottom", "color": "000000", "style": 1 },
 { "type": "right", "color": "000000", "style": 1 }]`
	styleTNR12, err := xlsx.NewStyle(`{"font":{"family": "Times New
Roman","size":12,"bold":false},
 "alignment":{"wrap_text":true, "vertical":"center", "horizontal":"center"},
 "fill":{"type":"pattern","color":["#FFFFFF"],"pattern":1}, "border":
` + vBorder + `}`)
	styleTNR12B, err := xlsx.NewStyle(`{"font":{"family": "Times New
Roman","size":12,"bold":true},
 "alignment":{"wrap_text":true, "vertical":"center", "horizontal":"center"},
 "fill":{"type":"pattern","color":["#FFFFFF"],"pattern":1}, "border":
` + vBorder + `}`)
	// ---------------------------------------------------------
	// задаем стили и название колонок
	nMeter := 0
	xlsx.SetRowHeight(sNameSheet, nMeter, 30*dCoeffY)
	xlsx.MergeCell(sNameSheet, ExcelXY2Name(nLeftCol, nUpperRow+nMeter),
		ExcelXY2Name(nLeftCol+4, nUpperRow+nMeter))
	err = xlsx.SetCellStyle(sNameSheet, ExcelXY2Name(nLeftCol,
		nUpperRow+nMeter), ExcelXY2Name(nLeftCol, nUpperRow+nMeter),
		styleTNR12B)
	xlsx.SetCellValue(sNameSheet, ExcelXY2Name(nLeftCol, nUpperRow+nMeter),
		"Отчет")
	nMeter++
	xlsx.SetRowHeight(sNameSheet, nMeter, 25*dCoeffY)
	err = xlsx.SetCellStyle(sNameSheet, ExcelXY2Name(nLeftCol,
		nUpperRow+nMeter), ExcelXY2Name(nLeftCol, nUpperRow+nMeter),
		styleTNR12B)
	sNameCol = ExcelX2Name(nLeftCol)
	xlsx.SetColWidth(sNameSheet, sNameCol,
		sNameCol, 35*dCoeffX)

	err = xlsx.SetCellStyle(sNameSheet, ExcelXY2Name(nLeftCol,
		nUpperRow+nMeter), ExcelXY2Name(nLeftCol, nUpperRow+nMeter),
		styleTNR12B)
	sNameCol = ExcelX2Name(nLeftCol)
	xlsx.SetColWidth(sNameSheet,
		sNameCol, sNameCol, 20.5*dCoeffX)
	xlsx.SetCellValue(sNameSheet, ExcelXY2Name(nLeftCol,
		nUpperRow+nMeter), "Категория разработки")
	err = xlsx.SetCellStyle(sNameSheet, ExcelXY2Name(nLeftCol+1,
		nUpperRow+nMeter), ExcelXY2Name(nLeftCol+1, nUpperRow+nMeter),
		styleTNR12B)
	sNameCol = ExcelX2Name(nLeftCol + 1)
	xlsx.SetColWidth(sNameSheet,
		sNameCol, sNameCol, 22.5*dCoeffX)
	xlsx.SetCellValue(sNameSheet, ExcelXY2Name(nLeftCol+1,
		nUpperRow+nMeter), "Язык программирования")
	err = xlsx.SetCellStyle(sNameSheet, ExcelXY2Name(nLeftCol+2,
		nUpperRow+nMeter), ExcelXY2Name(nLeftCol+2, nUpperRow+nMeter),
		styleTNR12B)
	sNameCol = ExcelX2Name(nLeftCol + 2)
	xlsx.SetColWidth(sNameSheet,
		sNameCol, sNameCol, 10*dCoeffX)
	xlsx.SetCellValue(sNameSheet, ExcelXY2Name(nLeftCol+2,
		nUpperRow+nMeter), "Стоимость")
	err = xlsx.SetCellStyle(sNameSheet, ExcelXY2Name(nLeftCol+3,
		nUpperRow+nMeter), ExcelXY2Name(nLeftCol+3, nUpperRow+nMeter),
		styleTNR12B)
	sNameCol = ExcelX2Name(nLeftCol + 3)
	xlsx.SetColWidth(sNameSheet,
		sNameCol, sNameCol, 20.5*dCoeffX)
	// ---------------------------------------------------------
	// формируем строки отчета
	for rows.Next() {
		nMeter++
		err := rows.Scan(&courseId, &devCat,
			&progLang, &cost)
		MODMLOG.CheckErr(err, "Не могу прочесть запись")

		//nItogoStudent = 0

		err = xlsx.SetCellStyle(sNameSheet, ExcelXY2Name(nLeftCol,
			nUpperRow+nMeter), ExcelXY2Name(nLeftCol, nUpperRow+nMeter),
			styleTNR12)
		xlsx.SetCellValue(sNameSheet, ExcelXY2Name(nLeftCol, nUpperRow+nMeter),
			devCat.String)
		//nItogoConf += n
		//nItogoStudent += n

		err = xlsx.SetCellStyle(sNameSheet, ExcelXY2Name(nLeftCol+1,
			nUpperRow+nMeter), ExcelXY2Name(nLeftCol+1, nUpperRow+nMeter),
			styleTNR12)
		xlsx.SetCellValue(sNameSheet, ExcelXY2Name(nLeftCol+1, nUpperRow+nMeter),
			progLang.String)
		//nItogoProject += n
		//nItogoStudent += n

		err = xlsx.SetCellStyle(sNameSheet, ExcelXY2Name(nLeftCol+2,
			nUpperRow+nMeter), ExcelXY2Name(nLeftCol+2, nUpperRow+nMeter),
			styleTNR12)
		xlsx.SetCellValue(sNameSheet, ExcelXY2Name(nLeftCol+2, nUpperRow+nMeter),
			cost.Int64)
		//nItogoPaper += n
		//nItogoStudent += n

	}
	nMeter++
	// формируем строку Итого

	// ---------------------------------------------------------
	xlsx.SetActiveSheet(index)
	file := xlsx
	sDate := time.Now().Format("2006-01-02")
	sDate = MODMLOG.DateToRus(sDate)
	// Устанавливается заголовок для отображения браузером загружаемого файла
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition",
		"attachment;filename="+sDate+"_Отчет.xlsx")
	w.Header().Set("File-Name", sDate+"_Отчет.xlsx") // userInputData.xlsx
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	err = file.Write(w)
	if err != nil {
		fmt.Println(err)
	}
	// ---------------------------------------------------------
}

/* ****************************************************** */
// функция построения отчета в формате docx
func ReportDOCX(rows *sql.Rows, w http.ResponseWriter) {
	// объявляем переменные
	//nItogoConf := 0
	//nItogoProject := 0
	//nItogoPaper := 0
	progLangCount := 0
	arr := make([]string, 0)
	devCatCount := 0
	fullCost := 0
	// объявляем переменные соответствующего типа для считывания полей
	var courseId sql.NullInt64
	var devCat sql.NullString
	var progLang sql.NullString
	var cost sql.NullString
	// формируем значение строки отчета
	for rows.Next() {
		err := rows.Scan(&courseId, &devCat,
			&progLang, &cost)
		MODMLOG.CheckErr(err, "Не могу прочесть запись")

		progLangCount++

		if len(arr) == 0 {
			arr = append(arr, devCat.String)
		}

		if Contains(devCat.String, arr) == false {
			arr = append(arr, devCat.String)
		}
		devCatCount = len(arr)
		if n, err := strconv.Atoi(cost.String); err == nil {
			fullCost += n
		}

	}
	// ---------------------------------------------------------
	// считываем файл образца
	rWord, err := docx.ReadDocxFile("mpage/report.docx")
	MODMLOG.CheckErr(err, "Не могу получить доступ к report.docx")
	docx1 := rWord.Editable()
	// заменяем в образце паттерны на считанные значения
	docx1.Replace("old_1_1", strconv.Itoa(devCatCount), -1)
	docx1.Replace("old_1_2", strconv.Itoa(progLangCount), -1)
	docx1.Replace("old_1_3", strconv.Itoa(fullCost), -1)
	// ---------------------------------------------------------
	file := docx1
	sDate := time.Now().Format("2006-01-02")
	sDate = MODMLOG.DateToRus(sDate)
	// Устанавливается заголовок для отображения браузером загружаемого файла
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=WhatsApp"+sDate+".docx")
	w.Header().Set("File-Name", "WhatsApp "+sDate+".docx")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	err = file.Write(w)
	if err != nil {
		fmt.Println(err)
	}
	rWord.Close()
	// ---------------------------------------------------------
}

func Contains(name string, arr []string) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == name {
			return true
		}
	}
	return false
}

// *********************************************************
// перевод координат ячеек в имя ячейки
func ExcelXY2Name(X int, Y int) string {
	c, err := excelize.CoordinatesToCellName(X, Y)
	MODMLOG.CheckErr(err, "Ошибка в ExcelXY2Name")
	return c
}

// *********************************************************
// перевод номера ячейки по X в имя ячейки по X
func ExcelX2Name(X int) string {
	c, err := excelize.ColumnNumberToName(X)
	MODMLOG.CheckErr(err, "Ошибка в ExcelX2Name")
	return c
}
