package apperrors

const defaultLang = "ru"

var messages = map[string]map[string]string{
	"ru": {
		"upload.invalid_file_format":   "Недопустимый формат файла. Поддерживаются только PDF и DOCX",
		"upload.file_too_large":        "Размер файла превышает допустимый лимит 20 МБ",
		"upload.empty_file":            "Загруженный файл пуст",
		"upload.missing_file":          "Файл не передан в запросе",
		"documents.invalid_pagination": "Некорректные параметры пагинации",
		"documents.not_found":          "Документ не найден",
		"upload.extraction_failed":     "Не удалось извлечь текст из документа",
		"upload.indexing_failed":       "Не удалось проиндексировать документ",
		"common.internal_error":        "Внутренняя ошибка сервера",
	},
	"en": {
		"upload.invalid_file_format":   "Unsupported file format. Only PDF and DOCX are allowed",
		"upload.file_too_large":        "File size exceeds the 20 MB limit",
		"upload.empty_file":            "Uploaded file is empty",
		"upload.missing_file":          "No file was provided in the request",
		"documents.invalid_pagination": "Invalid pagination parameters",
		"documents.not_found":          "Document not found",
		"upload.extraction_failed":     "Failed to extract text from the document",
		"upload.indexing_failed":       "Failed to index the document",
		"common.internal_error":        "Internal server error",
	},
}

type Translator struct{}

func NewTranslator() *Translator {
	return &Translator{}
}

func (t *Translator) Translate(messageKey, lang string) string {
	byLang, ok := messages[lang]
	if !ok {
		byLang = messages[defaultLang]
	}

	if msg, ok := byLang[messageKey]; ok {
		return msg
	}

	return messages[defaultLang]["common.internal_error"]
}
