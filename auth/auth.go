package auth

import (
	"encoding/json"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func Auth(login string, pass string, filename string) (bool, error) {

	users, err := LoadUsers(filename)
	if err != nil {
		return false, err
	}

	for _, user := range users {
		if user.Login == login && user.Password == pass {
			return true, nil
		}
	}
	return false, err
}

func Registration(login, password, filename string, w fyne.Window) error {

	users, err := LoadUsers(filename)
	if err != nil {
		fmt.Println("Ошибка LoadUsers")
		return err
	}

	for _, User := range users {
		if User.Login == login {
			dialog.ShowInformation("ошибка", "пользователь с таким логином уже существует", w)
			return err
		}
	}

	users = append(users, User{
		Login:    login,
		Password: password,
	})

	err = saveUsers(filename, users)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении пользователей: %v", err)
	}

	return nil

}

func LoadUsers(filename string) ([]User, error) {
	var users []User

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return users, nil
		}
		return nil, err
	}

	err = json.Unmarshal(data, &users)
	return users, err

}

func saveUsers(filename string, users []User) error {
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
