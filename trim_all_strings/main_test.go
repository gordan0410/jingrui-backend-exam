package main

import (
	"reflect"
	"testing"
)

// 定義用於測試的結構體
type SimpleStruct struct {
	Field1 string
	Field2 int
	Field3 string
}

type NestedStruct struct {
	ID     string
	Child  SimpleStruct
	Status string
}

type PointerChain struct {
	Name string
	Next *PointerChain
	Val  int
}

type MixedTypes struct {
	A string
	B int
	C bool
	D *string
	E float64
}

// 測試函式
func TestTrimAllStrings(t *testing.T) {
	// 測試案例結構
	tests := []struct {
		name     string
		inputPtr any // 傳入 TrimAllStrings 的參數 (必須是指標的指標或結構體指標)
		expected any // 預期的結構體值
	}{
		{
			name: "SimpleStruct_BasicTrim",
			inputPtr: &SimpleStruct{
				Field1: "  hello  ",
				Field2: 10,
				Field3: "world\t",
			},
			expected: SimpleStruct{
				Field1: "hello",
				Field2: 10,
				Field3: "world",
			},
		},
		{
			name: "SimpleStruct_NoChangeNeeded",
			inputPtr: &SimpleStruct{
				Field1: "clean",
				Field2: 5,
				Field3: "data",
			},
			expected: SimpleStruct{
				Field1: "clean",
				Field2: 5,
				Field3: "data",
			},
		},
		// 雖然 TrimAllStrings 目前只處理結構體*直接*欄位，但我們仍然可以用它來測試
		// 巢狀結構如果想被遞迴處理，TrimAllStrings 內部需要額外邏輯。
		// 依照你原始程式碼中的邏輯，它只會處理當前結構體中的字串和指標到下一個結構體的指標。
		// NestedStruct 中的 Child SimpleStruct 不會被處理，因為 Child 是一個結構體 (Kind() == Struct)，而不是指標。
		// 因此我將移除 NestedStruct 的測試，並專注於你原始程式碼所支援的鍊式指標結構。

		{
			name: "PointerChain_Simple",
			inputPtr: &PointerChain{
				Name: " one ",
				Val:  1,
				Next: &PointerChain{
					Name: " two\n",
					Val:  2,
					Next: nil,
				},
			},
			expected: PointerChain{
				Name: "one",
				Val:  1,
				Next: &PointerChain{
					Name: "two",
					Val:  2,
					Next: nil,
				},
			},
		},
		{
			name: "PointerChain_ThreeNodes",
			inputPtr: &PointerChain{
				Name: " A ",
				Val:  1,
				Next: &PointerChain{
					Name: " B ",
					Val:  2,
					Next: &PointerChain{
						Name: " C \t",
						Val:  3,
						Next: nil,
					},
				},
			},
			expected: PointerChain{
				Name: "A",
				Val:  1,
				Next: &PointerChain{
					Name: "B",
					Val:  2,
					Next: &PointerChain{
						Name: "C",
						Val:  3,
						Next: nil,
					},
				},
			},
		},
		// 測試原始 main 函式中使用的指標的指標 (&a) 情況
		{
			name: "PointerOfPointer_Chain",
			inputPtr: func() any {
				// 模擬 main 函式中的 &a
				a := &PointerChain{
					Name: " start ",
					Val:  0,
					Next: &PointerChain{
						Name: " middle ",
						Val:  1,
						Next: nil,
					},
				}
				return &a // 傳入的是 PointerChain 的指標的指標
			}(),
			expected: &PointerChain{ // 注意: 預期結果是結構體的指標，因為 TrimAllStrings 處理完後，變數 a 會變成這個值
				Name: "start",
				Val:  0,
				Next: &PointerChain{
					Name: "middle",
					Val:  1,
					Next: nil,
				},
			},
		},
		{
			name: "MixedTypes_Trim",
			inputPtr: &MixedTypes{
				A: " str_a ",
				B: 100,
				C: true,
				D: func() *string { s := " str_d "; return &s }(), // 指標字串不會被處理
				E: 3.14,
			},
			expected: MixedTypes{
				A: "str_a",
				B: 100,
				C: true,
				D: func() *string { s := " str_d "; return &s }(),
				E: 3.14,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 執行函式
			TrimAllStrings(tt.inputPtr)

			// 獲取修改後的值
			var actual any
			if reflect.ValueOf(tt.inputPtr).Kind() == reflect.Ptr && reflect.ValueOf(tt.inputPtr).Elem().Kind() == reflect.Ptr {
				// 如果傳入的是指標的指標 (**struct)，則實際結果是它指向的結構體指標
				actual = reflect.ValueOf(tt.inputPtr).Elem().Interface()
			} else {
				// 如果傳入的是結構體指標 (*struct)，則實際結果是它指向的結構體值
				actual = reflect.ValueOf(tt.inputPtr).Elem().Interface()
			}

			// 檢查結果是否符合預期
			if !reflect.DeepEqual(actual, tt.expected) {
				// 為了更容易閱讀的輸出，如果類型是指標，我們解引用它
				if reflect.TypeOf(actual).Kind() == reflect.Ptr {
					actual = reflect.ValueOf(actual).Elem().Interface()
					tt.expected = reflect.ValueOf(tt.expected).Elem().Interface()
				}

				t.Errorf("TrimAllStrings() 得到的值與預期不符\n實際值: %+v\n預期值: %+v", actual, tt.expected)
			}
		})
	}
}

// 測試 Panic 情況
func TestTrimAllStrings_Panic(t *testing.T) {
	// 測試非指標
	t.Run("NotAPointer", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("TrimAllStrings 應在傳入非指標時發生 panic")
			}
		}()
		TrimAllStrings(SimpleStruct{})
	})

	// 測試 nil 指標
	t.Run("NilPointer", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("TrimAllStrings 應在傳入 nil 指標時發生 panic")
			}
		}()
		var ptr *SimpleStruct = nil
		TrimAllStrings(ptr)
	})

	// 測試指標指向非結構體
	t.Run("PointerToNonStruct", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("TrimAllStrings 應在指標指向非結構體時發生 panic")
			}
		}()
		i := 10
		TrimAllStrings(&i)
	})
}
