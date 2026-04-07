package i18n

const DefaultLang = "uz"

var Translations = map[string]map[string]string{ //nolint:gochecknoglobals // static translation maps
	"en": {
		"INCORRECT_CREDENTIALS": "Incorrect credentials",
		"ADMIN_NOT_FOUND":       "Admin not found",
		"USER_NOT_FOUND":        "User not found",
		"USERNAME_CONFLICT":     "Username already exists",
		"SESSION_NOT_FOUND":     "Session not found",
		"VALIDATION_FAILED":     "Validation failed",
		"UNAUTHORIZED":          "Unauthorized",
		"FORBIDDEN":             "Forbidden",
		"SESSION_EXPIRED":       "Session expired",
		"USER_INACTIVE":         "User account is inactive",
		"UNSPECIFIED":           "Unspecified error",
	},

	"ru": {
		"INCORRECT_CREDENTIALS": "Неверные входные данные",
		"ADMIN_NOT_FOUND":       "Администратор не найден",
		"USER_NOT_FOUND":        "Пользователь не найден",
		"USERNAME_CONFLICT":     "Имя пользователя уже занято",
		"SESSION_NOT_FOUND":     "Сессия не найдена",
		"VALIDATION_FAILED":     "Ошибка валидации",
		"UNAUTHORIZED":          "Не авторизован",
		"FORBIDDEN":             "Доступ запрещён",
		"SESSION_EXPIRED":       "Сессия истекла",
		"USER_INACTIVE":         "Учётная запись неактивна",
		"UNSPECIFIED":           "Неизвестная ошибка",
	},

	"uz": {
		"INCORRECT_CREDENTIALS": "Kirish ma'lumotlari noto'g'ri",
		"ADMIN_NOT_FOUND":       "Administrator topilmadi",
		"USER_NOT_FOUND":        "Foydalanuvchi topilmadi",
		"USERNAME_CONFLICT":     "Foydalanuvchi nomi band",
		"SESSION_NOT_FOUND":     "Sessiya topilmadi",
		"VALIDATION_FAILED":     "Validatsiyadan o'tmadi",
		"UNAUTHORIZED":          "Avtorizatsiyadan o'tilmagan",
		"FORBIDDEN":             "Ruxsat berilmagan",
		"SESSION_EXPIRED":       "Sessiya muddati tugagan",
		"USER_INACTIVE":         "Foydalanuvchi aktiv emas",
		"UNSPECIFIED":           "Kutilmagan xatolik",
	},
}
