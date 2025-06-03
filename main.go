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
	w := a.NewWindow("Тельянов ЭВМ 16в: Числа Фибоначчи и Простые числа")

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
			mainApp(w, login)
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
				err := auth.Registration(login, password, filename, w)
				if err != nil {
					dialog.ShowError(fmt.Errorf("ошибка регистрации"), w)
				} else {
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

func mainApp(w fyne.Window, login string) {
	w.Resize(fyne.NewSize(760, 0))
	w.SetFixedSize(true)


	usedAccount := widget.NewLabel("Вы управляете: " + login)

	fibonacciStepsEntry := widget.NewEntry()
	primeStepsEntry := widget.NewEntry()

	output := widget.NewLabel("")
	output.Wrapping = fyne.TextWrapWord

	output2 := widget.NewLabel("")
	output2.Wrapping = fyne.TextWrapWord

	scrollContainer := container.NewScroll(output)
	scrollContainer.SetMinSize(fyne.NewSize(380, 300))

	scrollContainer2 := container.NewScroll(output2)
	scrollContainer2.SetMinSize(fyne.NewSize(380, 300))

	row1 := container.NewHBox(
		layout.NewSpacer(),
		scrollContainer,
	)
	row2 := container.NewHBox(
		layout.NewSpacer(),
		scrollContainer2,
	)

	rowContainer := container.NewHBox(
		row1,
		row2,
	)

	backToLoginButton := widget.NewButton("Сменить учетную запись", func() {
		loginWindow(w)
	})
	backToLoginButton.Resize(fyne.NewSize(760, 40))

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

	analyzewidget := widget.NewLabel("")
	analyzewidget.Wrapping = fyne.TextWrapWord
	analyzecontainer := container.NewScroll(analyzewidget)
	analyzecontainer.SetMinSize(fyne.NewSize(760, 200))
	downloadButton.Hide()
	filenameEntry.Hide()

	analyzeContainerWithDButton := container.NewHBox(
		analyzecontainer,
	)

	analyzeButton := widget.NewButton("Анализировать данные", func() {
		backToLoginButton.Show()
		analyzeContainerWithDButton.Show()
		filenameEntry.Show()
		downloadButton.Show()
		rowContainer.Hide()
		dataanalyze.DataAnalyze(analyzewidget)
		// dataanalyze.DataAnalyze(output)
	})

	analyzeButton.Disable()

	startButton := widget.NewButton("Начать поиск", nil)

	startButton.OnTapped = func() {
		backToLoginButton.Hide()
		analyzeContainerWithDButton.Hide()
		downloadButton.Hide()
		filenameEntry.Hide()

		if rowContainer.Visible() {
		} else {
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

		if fibonacciSteps < 0 || primeSteps < 0 {
			dialog.ShowError(fmt.Errorf("неверное количество шагов для простых чисел"), w)
			return
		}

		output.SetText("Поиск начат...")
		output2.SetText("Поиск начат...")

		analyzeButton.Disable()
		startButton.Disable()
		startButton.SetText("Вычисление...")

		go func() {
			resultsChan := make(chan string)

			
			go func() {
				for result := range resultsChan {
					if strings.Contains(result, "Фибоначчи") {
						output.SetText(output.Text + "\n" + result)
						scrollContainer.ScrollToBottom()
					} else if strings.Contains(result, "простых чисел") {
						output2.SetText(output2.Text + "\n" + result)
						scrollContainer2.ScrollToBottom()
					} else {
						output.SetText(output.Text + "\n" + result)
						scrollContainer.ScrollToBottom()
					}
				}
			}()
			calculating.FAndPfoundind(fibonacciSteps, primeSteps, resultsChan, startButton, analyzeButton)

		}()
	}

	content := container.NewVBox(
		usedAccount,
		backToLoginButton,
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
