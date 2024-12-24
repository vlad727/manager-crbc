package getsa

import (
	"log"
	"net/http"
	"os"
	"sigs.k8s.io/yaml"
	"text/template"
	"webapp/getsacollect"
	"webapp/home/loggeduser"
)

/*
Функция `GetSa` выполняет следующие шаги:

1. **Получение информации о пользователе**: Функция вызывает `loggeduser.LoggedUserRun` для получения карты с именем пользователя и группами, к которым он принадлежит.
2. **Форматирование данных**: После получения карты, вы извлекаете имя пользователя и отправляете эту информацию вместе с группой в другую функцию `getsacollect.GetSaCollect`, которая возвращает две переменные: `M3` (карта) и `Sl1` (срез), который вы игнорируете.
3. **Преобразование данных в YAML**: Вы конвертируете карту `M3` в YAML и сохраняете результат в переменную `str`.
4. **Рендеринг шаблона**: Вы загружаете HTML-шаблон и выполняете его рендеринг с использованием структуры `Msg`, содержащей сообщение и имя пользователя.
*/
var (
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
)

func GetSa(w http.ResponseWriter, r *http.Request) {
	defer logger.Println("INFO: Func GetSa finished")
	logger.Println("INFO: Func GetSa started")
	// send request to parse, trim and decode jwt, get map with user and groups
	UserAndGroups := loggeduser.LoggedUserRun(r) // get logged user and groups, ex:  map[ose.test.user:[ipausers tuz-endless]]

	logger.Println("INFO: Got message form LoggedUserRun")
	log.Println(UserAndGroups)

	// name of logged user
	var username string
	// get logged user name from map
	for k, _ := range UserAndGroups {
		username = k
	}

	// Создаем новую карту для каждого вызова функции таким образом мы очистим данные из mNsAndSa
	mNsAndSa := make(map[string][]string)

	// send map to func GetSaCollect and return M3 map and Sl1 slice
	mNsAndSa, _ = getsacollect.GetSaCollect(UserAndGroups)
	// Sl1 will be skipped  because we don't need here, Sl1 it's slice with namespace name and service account name example below:
	// my-test-ns: my-test-sa
	// M3 it's map with namespace name and service account name like below:
	// ose-test-ns:
	// - default
	// - ose-sa

	// Marshal to yaml for out to web page
	yamlFile, err := yaml.Marshal(mNsAndSa)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// convert to string for struct if you do not convert it will be in bytes
	str := string(yamlFile)
	// parse html
	t, err := template.ParseFiles("tmpl/getsa.html")
	if err != nil {
		log.Printf("Error parsing template: %v\n", err)
		return
	}
	// init struct and var
	Msg := struct {
		Message           string `yaml:"message"`
		MessageLoggedUser string
	}{
		Message:           str,
		MessageLoggedUser: username,
	}
	// execute
	err = t.Execute(w, Msg)
	if err != nil {
		return
	}
	mNsAndSa = nil
}
