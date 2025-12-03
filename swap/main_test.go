package main

import (
	"reflect"
	"testing"
)

// 假設您的 swap 函數簽名是: func swap[T any](a, b T)
// 且它在 main 包中。

// 幫助函數: 用於執行實際的交換和檢查
func testSwap(t *testing.T, a interface{}, b interface{}, expectedA interface{}, expectedB interface{}) {
	t.Helper()

	// 執行交換時捕獲可能的 panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unexpected panic: %v", r)
		}
	}()

	// 調用您的泛型 swap 函數 (這裡假設類型 T 可以被正確推導)
	swap(a, b)

	// 獲取交換後的值
	aValue := reflect.ValueOf(a).Elem().Interface()
	bValue := reflect.ValueOf(b).Elem().Interface()

	// 比較結果
	if !reflect.DeepEqual(aValue, expectedA) {
		t.Errorf("After swap, *a got %v (type %T), want %v (type %T)", aValue, aValue, expectedA, expectedA)
	}

	if !reflect.DeepEqual(bValue, expectedB) {
		t.Errorf("After swap, *b got %v (type %T), want %v (type %T)", bValue, bValue, expectedB, expectedB)
	}
}

// =========================================================================
// 基本類型測試 (新增地址 Log)
// =========================================================================

func TestSwap_Int(t *testing.T) {
	a := 10
	b := 20
	// aPtr 和 bPtr 就是傳遞給 swap 的 *int 類型的值 (&a 和 &b)
	aPtr := &a
	bPtr := &b

	t.Logf("--- TestSwap_Int ---")
	t.Logf("Original address: a = %p, b = %p", aPtr, bPtr)
	t.Logf("Before swap: *a = %d, *b = %d", *aPtr, *bPtr)

	testSwap(t, aPtr, bPtr, 20, 10)

	t.Logf("After swap: *a = %d, *b = %d", *aPtr, *bPtr)
	// 驗證地址保持不變，只有值被交換
	t.Logf("Final address: a = %p, b = %p (Addresses must be the same)", aPtr, bPtr)
}

func TestSwap_String(t *testing.T) {
	a := "Hello"
	b := "World"
	aPtr := &a
	bPtr := &b

	t.Logf("--- TestSwap_String ---")
	t.Logf("Original address: a = %p, b = %p", aPtr, bPtr)
	t.Logf("Before swap: *a = %s, *b = %s", *aPtr, *bPtr)

	testSwap(t, aPtr, bPtr, "World", "Hello")

	t.Logf("After swap: *a = %s, *b = %s", *aPtr, *bPtr)
	t.Logf("Final address: a = %p, b = %p", aPtr, bPtr)
}

func TestSwap_Float64(t *testing.T) {
	a := 3.14159
	b := 2.71828
	aPtr := &a
	bPtr := &b

	t.Logf("--- TestSwap_Float64 ---")
	t.Logf("Original address: a = %p, b = %p", aPtr, bPtr)
	t.Logf("Before swap: *a = %f, *b = %f", *aPtr, *bPtr)

	testSwap(t, aPtr, bPtr, 2.71828, 3.14159)

	t.Logf("After swap: *a = %f, *b = %f", *aPtr, *bPtr)
	t.Logf("Final address: a = %p, b = %p", aPtr, bPtr)
}

func TestSwap_Bool(t *testing.T) {
	a := true
	b := false
	aPtr := &a
	bPtr := &b

	t.Logf("--- TestSwap_Bool ---")
	t.Logf("Original address: a = %p, b = %p", aPtr, bPtr)
	t.Logf("Before swap: *a = %t, *b = %t", *aPtr, *bPtr)

	testSwap(t, aPtr, bPtr, false, true)

	t.Logf("After swap: *a = %t, *b = %t", *aPtr, *bPtr)
	t.Logf("Final address: a = %p, b = %p", aPtr, bPtr)
}

// =========================================================================
// 複雜類型測試 (新增地址 Log)
// =========================================================================

type Point struct {
	X       int
	Y       string
	IsFound bool
	Weight  float64
}

func TestSwap_Struct(t *testing.T) {
	a := Point{X: 1, Y: "A", IsFound: true, Weight: 1.0}
	b := Point{X: 3, Y: "B", IsFound: false, Weight: 3.3}
	expectedA := Point{X: 3, Y: "B", IsFound: false, Weight: 3.3}
	expectedB := Point{X: 1, Y: "A", IsFound: true, Weight: 1.0}

	aPtr := &a
	bPtr := &b

	t.Logf("--- TestSwap_Struct ---")
	t.Logf("Original address: a = %p, b = %p", aPtr, bPtr)
	t.Logf("Before swap: *a = %+v, *b = %+v", *aPtr, *bPtr)

	testSwap(t, aPtr, bPtr, expectedA, expectedB)

	t.Logf("After swap: *a = %+v, *b = %+v", *aPtr, *bPtr)
	t.Logf("Final address: a = %p, b = %p", aPtr, bPtr)
}

func TestSwap_Slice(t *testing.T) {
	a := []int{1, 2}
	b := []int{3, 4}

	aPtr := &a
	bPtr := &b

	t.Logf("--- TestSwap_Slice ---")
	t.Logf("Original address: a = %p, b = %p", aPtr, bPtr)
	t.Logf("Before swap: *a = %v, *b = %v", *aPtr, *bPtr)

	testSwap(t, aPtr, bPtr, []int{3, 4}, []int{1, 2})

	t.Logf("After swap: *a = %v, *b = %v", *aPtr, *bPtr)
	t.Logf("Final address: a = %p, b = %p", aPtr, bPtr)
}

// =========================================================================
// 錯誤處理測試 (保持不變)
// =========================================================================

func TestSwap_NonPointerPanic(t *testing.T) {
	t.Log("Testing Non-Pointer Panic...")
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic when 'a' was not a pointer")
		} else if r != "a is not a pointer/interface" {
			t.Errorf("Unexpected panic message: %v", r)
		}
	}()

	a := 10
	b := 20

	swap(a, b)
}

func TestSwap_NilPointerPanic(t *testing.T) {
	t.Log("Testing Nil-Pointer Panic...")
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic when 'a' was nil")
		} else if r != "a or b is nil" {
			t.Errorf("Unexpected panic message: %v", r)
		}
	}()

	var a *int
	b := 20
	swap(a, &b)
}

func TestSwap_NilInterfacePanic(t *testing.T) {
	t.Log("Testing Nil-Interface Panic...")
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic when 'a' was nil")
		} else if r != "a is nil" {
			t.Errorf("Unexpected panic message: %v", r)
		}
	}()

	var a, b interface{}
	
	swap(a, b)
}
