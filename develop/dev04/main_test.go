package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchAnagram(t *testing.T){
	inData:= []string{"пятак", "пЯтка", "тяпка", "тяпка", "Листок", "слиток", "столик", "кофта"}
	outData:= map[string]*[]string{
		"пятак": {"пятак", "пятка", "тяпка"},
		"листок": {"листок", "слиток", "столик"},
	}
	data:= SearchAnagram(&inData)
	require.Equal(t, outData, *data)
}