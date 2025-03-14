package main

import (
	"evm/auth"
	"evm/calculating"
	dataanalyze "evm/data-analyze"
	"evm/database"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/joho/godotenv"
)

func main() {

	a := app.New()
	w := a.NewWindow("Числа Фибоначчи и Простые числа")

	loginWindow(w)
	w.ShowAndRun()

}

func loginWindow(w fyne.Window) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	regToken := os.Getenv("TOKEN")
	filename := "auth/accountsDB/users.json"
	w.Resize(fyne.NewSize(400, 200))

	loginEntry := widget.NewEntry()
	loginEntry.SetPlaceHolder("Введите логин")

	passwordEntry := widget.NewEntry()
	passwordEntry.SetPlaceHolder("Введите пароль")

	loginButton := widget.NewButton("Войти", func() {
		login := loginEntry.Text
		password := passwordEntry.Text
		isAuthed, err := auth.Auth(login, password, filename)
		if err != nil {
			dialog.ShowError(fmt.Errorf("auth err"), w)
		}
		if isAuthed {
			mainApp(w)
		} else {
			dialog.ShowError(fmt.Errorf("неверный логин или пароль"), w)
		}
	})
	createAccount := widget.NewButton("Зарегестрироваться", func() {
		createNewAccount(w, regToken)
	})

	authForm := container.NewVBox(
		loginEntry,
		passwordEntry,
		createAccount,
		loginButton,
	)
	w.SetContent(authForm)

}

func createNewAccount(w fyne.Window, regToken string) {
	filename := "auth/accountsDB/users.json"
	w.Resize(fyne.NewSize(400, 200))
	createAccountLoginEntry := widget.NewEntry()
	createAccountLoginEntry.SetPlaceHolder("Введите логин")

	createAccountPassworEntry := widget.NewEntry()
	createAccountPassworEntry.SetPlaceHolder("Введите пароль")

	createAccountTokenEntry := widget.NewEntry()
	createAccountTokenEntry.SetPlaceHolder("Введите код приглашения в систему")

	sendData := widget.NewButton("Зарегестрироваться", func() {
		login := createAccountLoginEntry.Text
		password := createAccountPassworEntry.Text
		token := createAccountTokenEntry.Text

		if token != regToken {
			dialog.ShowError(fmt.Errorf("неверный код-приглашение"), w)
			return
		} else {
			if login != "" || password != "" {
				err := auth.Registration(login, password, filename)
				if err != nil {
					dialog.ShowError(fmt.Errorf("ошибка регистрации"), w)
				} else {
					dialog.ShowInformation("Успех", "Регистрация прошла успешно!", w)
					loginWindow(w)
				}

			}
		}

	})

	createNewAccountForm := container.NewVBox(
		createAccountLoginEntry,
		createAccountPassworEntry,
		createAccountTokenEntry,
		sendData,
	)

	w.SetContent(createNewAccountForm)

}

func mainApp(w fyne.Window) {
	w.Resize(fyne.NewSize(760, 300))

	fibonacciStepsEntry := widget.NewEntry()
	primeStepsEntry := widget.NewEntry()

	// Создаем поле для вывода результатов
	output := widget.NewLabel("")
	output.Wrapping = fyne.TextWrapWord // Перенос текста на новую строку

	//Второе поле резов
	output2 := widget.NewLabel("")
	output2.Wrapping = fyne.TextWrapWord

	// скачать резы отдельно

	// Обертываем поле вывода в контейнер с прокруткой
	scrollContainer := container.NewScroll(output)
	scrollContainer.SetMinSize(fyne.NewSize(380, 300)) // Фиксируем размер контейнера с прокруткой

	// Обвертка второго
	scrollContainer2 := container.NewScroll(output2)
	scrollContainer2.SetMinSize(fyne.NewSize(380, 300))

	//1
	row1 := container.NewHBox(
		layout.NewSpacer(),
		scrollContainer,
	)
	//2
	row2 := container.NewHBox(
		layout.NewSpacer(),
		scrollContainer2,
	)

	rowContainer := container.NewHBox(
		row1,
		row2,
	)

	downloadButton := widget.NewButton("Выгрузить результаты в отдельный файл", func() {
	})
	downloadButton.Resize(fyne.NewSize(760, 40))

	filenameEntry := widget.NewEntry()
	filenameEntry.SetPlaceHolder("Введите желаемое название файла, он сохранится в папке database в формате .txt")
	filenameEntry.Resize(fyne.NewSize(760, 40))

	downloadButton.OnTapped = func() {

		filename := filenameEntry.Text
		if filename == "" {
			dialog.ShowError(fmt.Errorf("неверное название файла"), w)
			return
		}
		err := database.UniqueSave(filename)
		if err != nil {
			rowContainer.Show()
			output.SetText("err")
		}
		downloadButton.SetIcon(theme.CheckButtonCheckedIcon())
		downloadButton.Importance = widget.HighImportance
		downloadButton.Refresh()
		go func() {
			time.Sleep(2 * time.Second)
			downloadButton.Importance = widget.MediumImportance
			downloadButton.SetIcon(nil)
			downloadButton.Refresh()
		}()
	}

	//анализ виджет
	analyzewidget := widget.NewLabel("")
	analyzewidget.Wrapping = fyne.TextWrapWord
	analyzecontainer := container.NewScroll(analyzewidget)
	analyzecontainer.SetMinSize(fyne.NewSize(760, 300))
	downloadButton.Hide()
	filenameEntry.Hide()

	analyzeContainerWithDButton := container.NewHBox(
		analyzecontainer,
	)

	// Канал для передачи данных о найденных числах
	resultsChan := make(chan string)

	// Горутина для обновления интерфейса в реальном времени
	go func() {
		for result := range resultsChan {
			// Разделяем вывод: числа Фибоначчи в output1, простые числа в output2
			if strings.Contains(result, "Фибоначчи") {
				output.SetText(output.Text + "\n" + result)
				scrollContainer.ScrollToBottom() // Автоматическая прокрутка вниз
			} else if strings.Contains(result, "простых чисел") {
				output2.SetText(output2.Text + "\n" + result)
				scrollContainer2.ScrollToBottom() // Автоматическая прокрутка вниз
			} else {
				// Если сообщение не подходит ни под одно условие, выводим его в output1
				output.SetText(output.Text + "\n" + result)
				scrollContainer.ScrollToBottom()
			}
		}
	}()

	startButton := widget.NewButton("Начать поиск", func() {
		analyzeContainerWithDButton.Hide()
		downloadButton.Hide()
		filenameEntry.Hide()
		if rowContainer.Visible() {
		} else {
			// Если новый виджет скрыт, показываем его и скрываем старые
			rowContainer.Show()
		}

		fibonacciSteps, err := strconv.Atoi(fibonacciStepsEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("неверное количество шагов для чисел Фибоначчи"), w)
			return
		}

		primeSteps, err := strconv.Atoi(primeStepsEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("неверное количество шагов для простых чисел"), w)
			return
		}

		output.SetText("Поиск начат...")
		output2.SetText("Поиск начат...")
		go func() {
			calculating.FAndPfoundind(fibonacciSteps, primeSteps, resultsChan)
		}()
	})

	analyzeButton := widget.NewButton("Анализировать данные", func() {
		analyzeContainerWithDButton.Show()
		filenameEntry.Show()
		downloadButton.Show()
		rowContainer.Hide()
		dataanalyze.DataAnalyze(analyzewidget)
		// dataanalyze.DataAnalyze(output)
	})

	// Создаем основной контейнер с элементами интерфейса
	content := container.NewVBox(
		widget.NewLabel("Количество шагов для чисел Фибоначчи:"),
		fibonacciStepsEntry,
		widget.NewLabel("Количество шагов для простых чисел:"),
		primeStepsEntry,
		startButton,
		analyzeButton,
		rowContainer,
		analyzeContainerWithDButton,
		filenameEntry,
		downloadButton,
	)

	w.SetContent(content)

}
