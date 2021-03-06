/*
Copyright 2019 by ofunc

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
*/

package lmodxlsx

import (
	"ofunc/lua"

	"github.com/plandem/xlsx/options"
)

func metaRow(l *lua.State, mcell int) int {
	l.NewTable(0, 8)
	idx := l.AbsIndex(-1)

	l.Push("cell")
	l.PushClosure(lRowCell, mcell)
	l.SetTableRaw(idx)

	l.Push("cells")
	l.PushClosure(lRowCells, mcell)
	l.SetTableRaw(idx)

	l.Push("clear")
	l.Push(lRowClear)
	l.SetTableRaw(idx)

	l.Push("reset")
	l.Push(lRowReset)
	l.SetTableRaw(idx)

	l.Push("copyto")
	l.Push(lRowCopyTo)
	l.SetTableRaw(idx)

	l.Push("set")
	l.Push(lRowSet)
	l.SetTableRaw(idx)

	l.Push("__index")
	l.PushClosure(lRowIndex, idx)
	l.SetTableRaw(idx)

	l.Push("__newindex")
	l.Push(lRowNewIndex)
	l.SetTableRaw(idx)

	return idx
}

func lRowCell(l *lua.State) int {
	i := int(l.ToInteger(2)) - 1
	l.Push(toRow(l, 1).Cell(i))
	l.PushIndex(lua.FirstUpVal - 1)
	l.SetMetaTable(-2)
	return 1
}

func lRowCells(l *lua.State) int {
	iter := toRow(l, 1).Cells()
	l.PushClosure(func(l *lua.State) int {
		if !iter.HasNext() {
			return 0
		}
		cidx, _, cell := iter.Next()
		l.Push(cidx + 1)
		l.Push(cell)
		l.PushIndex(lua.FirstUpVal - 1)
		l.SetMetaTable(-2)
		return 2
	}, lua.FirstUpVal-1)
	return 1
}

func lRowClear(l *lua.State) int {
	toRow(l, 1).Clear()
	return 0
}

func lRowReset(l *lua.State) int {
	toRow(l, 1).Reset()
	return 0
}

func lRowCopyTo(l *lua.State) int {
	i := int(l.ToInteger(2)) - 1
	o := l.ToBoolean(2)
	toRow(l, 1).CopyTo(i, o)
	return 0
}

func lRowSet(l *lua.State) int {
	l.Push("level")
	l.GetTable(2)
	level := uint8(l.OptInteger(-1, 1) - 1)

	l.Push("collapsed")
	l.GetTable(2)
	collapsed := l.ToBoolean(-1)

	l.Push("phonetic")
	l.GetTable(2)
	phonetic := l.ToBoolean(-1)

	l.Push("hidden")
	l.GetTable(2)
	hidden := l.ToBoolean(-1)

	l.Push("height")
	l.GetTable(2)
	height := float32(l.OptFloat(-1, 0))

	toRow(l, 1).Set(&options.RowOptions{
		OutlineLevel: level,
		Collapsed:    collapsed,
		Phonetic:     phonetic,
		Hidden:       hidden,
		Height:       height,
	})
	return 0
}

func lRowIndex(l *lua.State) int {
	row := toRow(l, 1)
	switch key := l.ToString(2); key {
	case "format":
		l.Push(row.Formatting())
	case "index":
		l.Push(row.Bounds().FromRow + 1)
	default:
		l.Push(key)
		l.GetTableRaw(lua.FirstUpVal - 1)
	}
	return 1
}

func lRowNewIndex(l *lua.State) int {
	row := toRow(l, 1)
	switch key := l.ToString(2); key {
	case "format":
		row.SetFormatting(toStyle(l, 3))
	default:
		panic("xlsx.row: invalid field: " + key)
	}
	return 1
}
