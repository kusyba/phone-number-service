package utils

import (
    "regexp"
    "strings"
)

var countryCodes = map[string]string{
    "+7":   "Россия",
    "+1":   "США",
    "+44":  "Великобритания",
    "+49":  "Германия",
    "+33":  "Франция",
    "+86":  "Китай",
    "+81":  "Япония",
    "+82":  "Южная Корея",
    "+55":  "Бразилия",
    "+91":  "Индия",
}

var defMap = map[string]struct {
    Region   string
    Provider string
}{
    "910": {"Москва", "МТС"},
    "911": {"Москва", "МТС"},
    "912": {"Екатеринбург", "МТС"},
    "913": {"Новосибирск", "МТС"},
    "916": {"Москва", "МТС"},
    "920": {"Москва", "Мегафон"},
    "921": {"Санкт-Петербург", "Мегафон"},
    "926": {"Москва", "Мегафон"},
    "930": {"Москва", "Билайн"},
    "950": {"Москва", "Tele2"},
}

func NormalizeToE164(number string) string {
    number = strings.ReplaceAll(number, " ", "")
    number = strings.ReplaceAll(number, "-", "")
    number = strings.ReplaceAll(number, "(", "")
    number = strings.ReplaceAll(number, ")", "")
    
    if strings.HasPrefix(number, "+") {
        if matched, _ := regexp.MatchString(`^\+\d{10,15}$`, number); matched {
            return number
        }
        return ""
    }
    
    if strings.HasPrefix(number, "8") && len(number) == 11 {
        return "+7" + number[1:]
    }
    
    if len(number) == 10 && regexp.MustCompile(`^\d{10}$`).MatchString(number) {
        return "+7" + number
    }
    
    if strings.HasPrefix(number, "7") && len(number) == 11 {
        return "+" + number
    }
    
    return ""
}

func ValidateE164(number string) bool {
    matched, _ := regexp.MatchString(`^\+\d{10,15}$`, number)
    return matched
}

func GetCountryByCode(number string) string {
    for code, country := range countryCodes {
        if strings.HasPrefix(number, code) {
            return country
        }
    }
    return "Неизвестно"
}

func GetRussianRegionAndProvider(number string) (string, string) {
    if !strings.HasPrefix(number, "+7") || len(number) < 5 {
        return "Россия", "Неизвестно"
    }
    
    def := number[2:5]
    
    if info, exists := defMap[def]; exists {
        return info.Region, info.Provider
    }
    
    return "Россия", "Другой оператор"
}
