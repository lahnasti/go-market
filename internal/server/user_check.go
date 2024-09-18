package server

import "regexp"
// проверки регулярных выражений и уникальности логина
const (
	usernameR = `^[a-zA-Z0-9_]{3,20}$`
	passR = `^[A-Za-z\d]{8,}$`
) 
/*
passR:
^ — указывает на начало строки.
(?=.*[A-Za-z]) — утверждение, что строка должна содержать хотя бы одну букву:
[A-Za-z] — любая буква (как заглавная, так и строчная).
.* — позволяет находить буквы в любом месте строки.
(?=.*\d) — утверждение, что строка должна содержать хотя бы одну цифру:
\d — любая цифра.
[A-Za-z\d]{8,} — сама строка, состоящая из букв и цифр, и имеет длину не менее 8 символов.
$ — указывает на конец строки.
*/

func isValidUsername(username string)bool {
	usernameR := regexp.MustCompile(usernameR)
	return usernameR.MatchString(username)
}

func isValidPass(password string)bool {
	passR := regexp.MustCompile(passR)
    return passR.MatchString(password)
}

func(s *Server) isUsernameUnique(username string)(bool, error) {
	return s.Db.IsUsernameUnique(username)
}