// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

var xgid = RootGroupID.Extend("xgid")

func isEmpty(g Grouping) bool {
	return g.Columns() == nil
}

func de(x, y interface{}) bool {
	return reflect.DeepEqual(x, y)
}

func equal(g1, g2 Grouping) bool {
	if !de(g1.Columns(), g2.Columns()) ||
		!de(g1.Groups(), g2.Groups()) {
		return false
	}
	for _, gid := range g1.Groups() {
		for _, col := range g1.Columns() {
			if !de(g1.Table(gid).Column(col), g2.Table(gid).Column(col)) {
				return false
			}
		}
	}
	return true
}

func shouldPanic(t *testing.T, re string, f func()) {
	r := regexp.MustCompile(re)
	defer func() {
		err := recover()
		if err == nil {
			t.Fatalf("want panic matching %q; got no panic", re)
		} else if !r.MatchString(fmt.Sprintf("%s", err)) {
			t.Fatalf("want panic matching %q; got %s", re, err)
		}
	}()
	f()
}

func TestEmptyTable(t *testing.T) {
	tab := new(Table)
	if !isEmpty(tab) {
		t.Fatalf("Table{} is not empty")
	}
	tab.Add("x", []int{})
	tab.Add("x", []int{1, 2, 3})
	shouldPanic(t, "not a slice", func() {
		tab.Add("x", 1)
	})
	if v := tab.Len(); v != 0 {
		t.Fatalf("Table{}.Len() should be 0; got %v", v)
	}
	if v := tab.Columns(); v != nil {
		t.Fatalf("Table{}.Columns() should be nil; got %v", v)
	}
	if v := tab.Column("x"); v != nil {
		t.Fatalf("Table{}.Column(\"x\") should be nil; got %v", v)
	}
	shouldPanic(t, "unknown column", func() {
		tab.MustColumn("x")
	})
	if v, w := tab.Groups(), []GroupID{}; !de(v, w) {
		t.Fatalf("Table{}.Groups should be %v; got %v", w, v)
	}
	if v := tab.Table(RootGroupID); v != nil {
		t.Fatalf("Table{}.Table(RootGroupID) should be nil; got %v", v)
	}
	if v := tab.Table(xgid); v != nil {
		t.Fatalf("Table{}.Table(xgid) should be nil; got %v", v)
	}

	tab1 := tab.Add("x", []int{})
	if tab2 := tab.AddTable(RootGroupID, tab); !isEmpty(tab2) {
		t.Fatalf("Table{}.AddTable(RootGroupID, Table{}) should be empty; got %v", tab2)
	}
	if tab2 := tab.AddTable(RootGroupID, tab1); !equal(tab1, tab2) {
		t.Fatalf("tab.AddTable(RootGroupID, tab1) should be tab1; got %v", tab2)
	}
	if tab2 := tab1.AddTable(RootGroupID, tab); !equal(tab1, tab2) {
		t.Fatalf("tab1.AddTable(RootGroupID, tab) should be tab1; got %v", tab2)
	}
}

func TestTable0(t *testing.T) {
	col := []int{}
	tab := new(Table).Add("x", col)
	if isEmpty(tab) {
		t.Fatalf("tab is empty")
	}
	tab.Add("x", []int{1}) // Can override only column.
	shouldPanic(t, "column y.* with 1 .* 0 rows", func() {
		tab.Add("y", []int{1})
	})
	tab.Add("y", []int{})
	if v := tab.Len(); v != 0 {
		t.Fatalf("tab.Len() should be 0; got %v", v)
	}
	if v, w := tab.Columns(), []string{"x"}; !de(v, w) {
		t.Fatalf("tab.Columns() should be %v; got %v", w, v)
	}
	if v := tab.Column("x"); !de(v, col) {
		t.Fatalf("tab.Column(\"x\") should be %v; got %v", col, v)
	}
	if v := tab.Column("y"); v != nil {
		t.Fatalf("tab.Column(\"y\") should be nil; got %v", v)
	}
	if v := tab.MustColumn("x"); !de(v, col) {
		t.Fatalf("tab.MustColumn(\"x\") should be %v; got %v", col, v)
	}
	shouldPanic(t, "unknown column", func() {
		tab.MustColumn("y")
	})
	if v, w := tab.Groups(), []GroupID{RootGroupID}; !de(v, w) {
		t.Fatalf("tab.Groups() should be %v; got %v", w, v)
	}
	if v := tab.Table(RootGroupID); v != tab {
		t.Fatalf("tab.Table(RootGroupID) should be %v; got %v", tab, v)
	}
	if v := tab.Table(xgid); v != nil {
		t.Fatalf("tab.Table(xgid) should be nil; got %v", v)
	}
}

func TestTable1(t *testing.T) {
	col := []int{1}
	tab := new(Table).Add("x", col)
	if isEmpty(tab) {
		t.Fatalf("tab is empty")
	}
	tab.Add("x", []int{}) // Can override only column.
	shouldPanic(t, "column y.* with 2 .* 1 rows", func() {
		tab.Add("y", []int{1, 2})
	})
	tab.Add("y", []int{1})
	if v := tab.Len(); v != 1 {
		t.Fatalf("tab.Len() should be 1; got %v", v)
	}
	if v, w := tab.Columns(), []string{"x"}; !de(v, w) {
		t.Fatalf("tab.Columns() should be %v; got %v", w, v)
	}
	if v := tab.Column("x"); !de(v, col) {
		t.Fatalf("tab.Column(\"x\") should be %v; got %v", col, v)
	}
	if v := tab.Column("y"); v != nil {
		t.Fatalf("tab.Column(\"y\") should be nil; got %v", v)
	}
	if v := tab.MustColumn("x"); !de(v, col) {
		t.Fatalf("tab.MustColumn(\"x\") should be %v; got %v", col, v)
	}
	shouldPanic(t, "unknown column", func() {
		tab.MustColumn("y")
	})
	if v, w := tab.Groups(), []GroupID{RootGroupID}; !de(v, w) {
		t.Fatalf("tab.Groups() should be %v; got %v", w, v)
	}
	if v := tab.Table(RootGroupID); v != tab {
		t.Fatalf("tab.Table(RootGroupID) should be %v; got %v", tab, v)
	}
	if v := tab.Table(xgid); v != nil {
		t.Fatalf("tab.Table(xgid) should be nil; got %v", v)
	}
}

func TestAddTable(t *testing.T) {
	tab0 := new(Table).Add("x", []int{})
	tab1 := new(Table).Add("x", []int{1})
	tabY := new(Table).Add("y", []int{})
	tabXY := new(Table).Add("x", []int{}).Add("y", []int{})

	if v := tab0.AddTable(RootGroupID, tab0); !equal(tab0, v) {
		t.Fatalf("tab0.AddTable(RootGroupID, tab0) should be %v; got %v", tab0, v)
	}
	if v := tab0.AddTable(RootGroupID, tab1); !equal(tab1, v) {
		t.Fatalf("tab0.AddTable(RootGroupID, tab1) should be %v; got %v", tab0, v)
	}
	if v := tab0.AddTable(RootGroupID, tabY); !equal(tabY, v) {
		t.Fatalf("tab0.AddTable(RootGroupID, tabY) should be %v; got %v", tab0, v)
	}
	shouldPanic(t, "table missing column: x", func() {
		tab0.AddTable(xgid, tabY)
	})
	shouldPanic(t, "table has extra column: y", func() {
		tab0.AddTable(xgid, tabXY)
	})

	tab01 := tab0.AddTable(xgid, tab1)
	if v, w := tab01.Columns(), []string{"x"}; !de(v, w) {
		t.Fatalf("tab01.Columns() should be %v; got %v", w, v)
	}
	if v, w := tab01.Groups(), []GroupID{RootGroupID, xgid}; !de(v, w) {
		t.Fatalf("tab01.Groups() should be %v; got %v", w, v)
	}
	if v := tab01.Table(RootGroupID); v != tab0 {
		t.Fatalf("tab01.Table(RootGroupID) should be tab0; got %v", v)
	}
	if v := tab01.Table(xgid); v != tab1 {
		t.Fatalf("tab01.Table(xgid) should be tab1; got %v", v)
	}
	if v := tab01.Table(RootGroupID.Extend("ygid")); v != nil {
		t.Fatalf("tab01.Table(ygid) should be nil; got %v", v)
	}
	if v := tab01.AddTable(RootGroupID, new(Table)); !equal(tab01, v) {
		t.Fatalf("tab01.AddTable(RootGroupID, empty) should be tab01; got %v", v)
	}
}
