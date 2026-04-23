package utils

import (
    "testing"
)

func TestNormalizeToE164(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  string
    }{
        {"Российский +7", "+79161234567", "+79161234567"},
        {"Российский 8", "89161234567", "+79161234567"},
        {"Российский 10 цифр", "9161234567", "+79161234567"},
        {"Международный UK", "+44 20 7946 0638", "+442079460638"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := NormalizeToE164(tt.input)
            if got != tt.want {
                t.Errorf("NormalizeToE164() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestGetCountryByCode(t *testing.T) {
    tests := []struct {
        input string
        want  string
    }{
        {"+79161234567", "Россия"},
        {"+442079460638", "Великобритания"},
    }

    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            if got := GetCountryByCode(tt.input); got != tt.want {
                t.Errorf("GetCountryByCode() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestGetRussianRegionAndProvider(t *testing.T) {
    tests := []struct {
        input        string
        wantRegion   string
        wantProvider string
    }{
        {"+79161234567", "Москва", "МТС"},
        {"+79261234567", "Москва", "Мегафон"},
        // Код 930 есть в defMap, поэтому возвращает Москва, Билайн
        {"+79301234567", "Москва", "Билайн"},
    }

    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            region, provider := GetRussianRegionAndProvider(tt.input)
            if region != tt.wantRegion {
                t.Errorf("GetRussianRegionAndProvider() region = %v, want %v", region, tt.wantRegion)
            }
            if provider != tt.wantProvider {
                t.Errorf("GetRussianRegionAndProvider() provider = %v, want %v", provider, tt.wantProvider)
            }
        })
    }
}
