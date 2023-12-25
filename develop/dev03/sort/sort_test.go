package sort

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStartSort(t *testing.T) {
	type test struct {
		name    string
		inData  [][]string
		outData [][]string
		keys    keys
	}

	cases := []test{
		{
			name: "successWithoutKeys",
			inData: [][]string{
				{"KKKKKKK", "96", "22488293"},
				{"aaaaaaa", "32", "kjndkjnk"},
				{"Aaaaaaa", "44", "fdggvgfv"},
				{"Bbbbbbb", "30", "kjnkjnkk"},
			},
			outData: [][]string{
				{"Aaaaaaa", "44", "fdggvgfv"},
				{"Bbbbbbb", "30", "kjnkjnkk"},
				{"KKKKKKK", "96", "22488293"},
				{"aaaaaaa", "32", "kjndkjnk"},
			},
			keys: keys{
				n: false,
				r: false,
				u: false,
				k: 0,
			},
		},
		{
			name: "successWithKeyR",
			inData: [][]string{
				{"KKKKKKK", "96", "22488293"},
				{"aaaaaaa", "32", "kjndkjnk"},
				{"Aaaaaaa", "44", "fdggvgfv"},
				{"Bbbbbbb", "30", "kjnkjnkk"},
			},
			outData: [][]string{
				{"aaaaaaa", "32", "kjndkjnk"},
				{"KKKKKKK", "96", "22488293"},
				{"Bbbbbbb", "30", "kjnkjnkk"},
				{"Aaaaaaa", "44", "fdggvgfv"},
			},
			keys: keys{
				n: false,
				r: true,
				u: false,
				k: 0,
			},
		},
		{
			name: "successWithKeyK",
			inData: [][]string{
				{"KKKKKKK", "96", "22488293"},
				{"aaaaaaa", "32", "kjndkjnk"},
				{"Aaaaaaa", "44", "fdggvgfv"},
				{"Bbbbbbb", "50", "kjnkjnkk"},
			},
			outData: [][]string{
				{"aaaaaaa", "32", "kjndkjnk"},
				{"Aaaaaaa", "44", "fdggvgfv"},
				{"Bbbbbbb", "50", "kjnkjnkk"},
				{"KKKKKKK", "96", "22488293"},
			},
			keys: keys{
				n: false,
				r: false,
				u: false,
				k: 1,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data := startSort(c.inData, &c.keys)
			require.Equal(t, data, c.outData)
		})
	}
}

func TestUnique(t *testing.T) {	
	inData  :=[][]string{
		{"KKKKKKK", "96", "22488293"},
		{"aaaaaaa", "32", "kjndkjnk"},
		{"aaaaaaa", "32", "kjndkjnk"},
		{"Aaaaaaa", "44", "fdggvgfv"},
		{"Bbbbbbb", "30", "kjnkjnkk"},
		{"Bbbbbbb", "30", "kjnkjnkk"},
	}
	outData := [][]string{
		{"KKKKKKK", "96", "22488293"},
		{"aaaaaaa", "32", "kjndkjnk"},
		{"Aaaaaaa", "44", "fdggvgfv"},
		{"Bbbbbbb", "30", "kjnkjnkk"},
	}

	data:=unique(inData)
	require.Equal(t, data, outData)
}

func TestValueSort(t *testing.T){
	type test struct {
		name    string
		inData  [][]string
		outData [][]string
		keys    keys
	}

	cases := []test{
		{
			name: "successWithoutKeys",
			inData: [][]string{
				{"KKKKKKK", "96", "22488293"},
				{"aaaaaaa", "32", "kjndkjnk"},
				{"Aaaaaaa", "44", "fdggvgfv"},
				{"Bbbbbbb", "30", "kjnkjnkk"},
			},
			outData: [][]string{
				{"Aaaaaaa", "44", "fdggvgfv"},
				{"Bbbbbbb", "30", "kjnkjnkk"},
				{"KKKKKKK", "96", "22488293"},
				{"aaaaaaa", "32", "kjndkjnk"},
			},
			keys: keys{
				n: false,
				r: false,
				u: false,
				k: 0,
			},
		},
		{
			name: "successWithKeyR",
			inData: [][]string{
				{"KKKKKKK", "96", "22488293"},
				{"aaaaaaa", "32", "kjndkjnk"},
				{"Aaaaaaa", "44", "fdggvgfv"},
				{"Bbbbbbb", "30", "kjnkjnkk"},
			},
			outData: [][]string{
				{"KKKKKKK", "96", "22488293"},
				{"Aaaaaaa", "44", "fdggvgfv"},
				{"aaaaaaa", "32", "kjndkjnk"},
				{"Bbbbbbb", "30", "kjnkjnkk"},
			},
			keys: keys{
				n: false,
				r: true,
				u: false,
				k: 1,
			},
		},
		{
			name: "successWithKeyK",
			inData: [][]string{
				{"KKKKKKK", "96", "22488293"},
				{"aaaaaaa", "32", "kjndkjnk"},
				{"Aaaaaaa", "44", "fdggvgfv"},
				{"Bbbbbbb", "30", "kjnkjnkk"},
			},
			outData: [][]string{
				{"Bbbbbbb", "30", "kjnkjnkk"},
				{"aaaaaaa", "32", "kjndkjnk"},
				{"Aaaaaaa", "44", "fdggvgfv"},
				{"KKKKKKK", "96", "22488293"},
			},
			keys: keys{
				n: false,
				r: false,
				u: false,
				k: 1,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data := valueSort(c.inData, &c.keys)
			require.Equal(t, data, c.outData)
		})
	}
}
