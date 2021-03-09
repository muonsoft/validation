package examples

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	languagepkg "github.com/muonsoft/language"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message/translations/russian"
	"golang.org/x/text/language"
)

type Book struct {
	Title    string   `json:"title"`
	Author   string   `json:"author"`
	Keywords []string `json:"keywords"`
}

func (b Book) Validate(validator *validation.Validator) error {
	return validator.Validate(
		validation.StringProperty("title", &b.Title, it.IsNotBlank()),
		validation.StringProperty("author", &b.Author, it.IsNotBlank()),
		validation.CountableProperty("keywords", len(b.Keywords), it.HasCountBetween(1, 10)),
		validation.EachStringProperty("keywords", b.Keywords, it.IsNotBlank()),
	)
}

func HandleBooks(writer http.ResponseWriter, request *http.Request) {
	var book Book
	err := json.NewDecoder(request.Body).Decode(&book)
	if err != nil {
		http.Error(writer, "invalid request", http.StatusBadRequest)
		return
	}

	// setting up validator
	validator, err := validation.NewValidator(validation.Translations(russian.Messages))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = validator.WithContext(request.Context()).ValidateValidatable(book)
	if err != nil {
		violations, ok := validation.UnwrapViolationList(err)
		if ok {
			response, _ := json.Marshal(violations)
			writer.WriteHeader(http.StatusUnprocessableEntity)
			writer.Header().Set("Content-Type", "application/json")
			writer.Write(response)
			return
		}

		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// handle valid book

	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte("ok"))
}

func ExampleValidator_Validate_httpHandler() {
	var handler http.Handler
	handler = http.HandlerFunc(HandleBooks)
	// middleware set up: we need to set supported languages
	// detected language will be passed via request context
	handler = languagepkg.NewMiddleware(handler, languagepkg.SupportedLanguages(language.English, language.Russian))
	server := httptest.NewServer(handler)
	defer server.Close()

	// creating request with the language-specific header
	request, _ := http.NewRequest(http.MethodPost, server.URL, strings.NewReader(`{}`))
	request.Header.Set("Accept-Language", "ru")

	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)

	// response should contain array of violations
	fmt.Println(string(responseBody))
	// Output:
	// [{"code":"notBlank","message":"Значение не должно быть пустым.","propertyPath":"title"},{"code":"notBlank","message":"Значение не должно быть пустым.","propertyPath":"author"},{"code":"countTooFew","message":"Эта коллекция должна содержать 1 элемент или больше.","propertyPath":"keywords"}]
}
