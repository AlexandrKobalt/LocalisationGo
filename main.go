package main

import (
	"encoding/json"
	"fmt"
	"languageExample/config"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Эта логика уезжает в create.go
	// Загружаем всю (да, всю) локализацию в ОЗУ для дальнейшего использования
	err := loadTranslations()
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	r.Post("/GetError", func(w http.ResponseWriter, r *http.Request) {
		GetError(w, r)
	})

	http.ListenAndServe(":8665", r)
}

// Да, это глобальная переменная, но её можно
// будет передавать с конфигом, а-ля Л2/Л6
var translations map[string]config.Localization

func loadTranslations() error {
	translations = make(map[string]config.Localization)
	files, err := os.ReadDir("config/localization")
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		language := strings.TrimSuffix(file.Name(), ".json")
		localization, err := loadLocalization(language)
		if err != nil {
			return err
		}
		translations[language] = localization
	}
	return nil
}

func loadLocalization(language string) (config.Localization, error) {
	filename := fmt.Sprintf("config/localization/%s.json", language)
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var l config.Localization
	err = json.Unmarshal(file, &l)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func GetTranslation(key, language string) (string, error) {
	l, ok := translations[language]
	if !ok {
		return "", fmt.Errorf("translation not found for language '%s'", language)
	}
	translation, ok := l[key]
	if !ok {
		return "", fmt.Errorf("translation not found for key '%s'", key)
	}
	return translation, nil
}

func GetError(w http.ResponseWriter, r *http.Request) {
	// Получаем язык клиента
	lang := r.Header.Get("Accept-Language")

	// ...
	// *некоторая логика, где сработала ошибка*
	// ...

	// Получаем нужный перевод ошибки по ключу и языку клиента
	errorText, err := GetTranslation("invalidJSON", lang)
	if err != nil {
		fmt.Println(err)
		// Логируем ошибку, в дальнейшем возможно доработать логику, чтобы
		// клиенту  автоматически отправлялось
		// "Sorry, something went wrong. Please, try letter"
	}

	// Отправляем ответ
	fmt.Fprint(w, errorText)
}
