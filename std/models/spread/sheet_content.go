package spread

import (
	"fmt"
	"errors"
	. "github.com/universe-10th/calculus/v2/std/expressions"
	. "github.com/universe-10th/calculus/v2/std/models"
)


// Details of each cell being used. This works as a
// sort of internal cache that will serve also to sort
// the execution of the whole spread sheet.
type usedCell struct {
	variable Variable
	row, column uint16
}
type usedCellsSet map[usedCell]bool


var ErrInvalidSize = errors.New("width and height must be both > 0")
var ErrOutOfBounds = errors.New("coordinates out of bounds")
var ErrFlowIsNotSingleOutput = errors.New("in-cell model must have just one output variable being the cell-related variable")
var ErrSheetContentIsLocked = errors.New("cannot edit this sheet content - it is locked / in use")


// SheetContent is the actual content of a spread sheet.
// The user will have ONE opportunity to edit the
// sheet content and this is upon instantiation of the
// actual, public, spread sheet.
type SheetContent struct {
	locked bool
	name string
	width, height uint16
	flows []ModelFlow
	// When locked, this member will empty.
	// Tells which cells are being filled.
	usedCells usedCellsSet
	// When locked, this member will be empty.
	// Tells which cells are being referenced.
	referencedCells Variables
}


// Name of this sheet. The name should be a slug (i.e. no spaces) for
// good practice's sake, but it is not strictly required.
func (sheetContent *SheetContent) Name() string {
	return sheetContent.name
}


// Height of this sheet.
func (sheetContent *SheetContent) Height() uint16 {
	return sheetContent.height
}


// Width of this sheet.
func (sheetContent *SheetContent) Width() uint16 {
	return sheetContent.width
}


// Retrieves a cell-related variable in the bounds of the current
// sheet content. Intended to be used only on parent sheet
// initialization.
func (sheetContent *SheetContent) Cell(row, column uint16) Variable {
	if row >= sheetContent.height || column >= sheetContent.width {
		panic(ErrOutOfBounds)
	}
	cellVar := Var(fmt.Sprintf("%s[%d:%d]", sheetContent.name, row, column))
	sheetContent.referencedCells[cellVar] = true
	return cellVar
}


// Puts an expression inside the current sheet content at a specific
// cell given by its row/column. Intended to be used only on parent
// sheet initialization.
func (sheetContent *SheetContent) Put(row, column uint16, expression Expression) error {
	var initializer func(Variable) (ModelFlow, error) = nil
	if expression != nil {
		initializer = func(cell Variable) (ModelFlow, error) {
			if flow, err := NewExpressionModelFlow(cell, expression); err != nil {
				return nil, err
			} else {
				return flow, nil
			}
		}
	}
	return sheetContent.PutFlow(row, column, initializer)
}


// Puts a custom flow inside the current sheet content at a specific
// cell given by its row/column. Intended to be used only on parent
// sheet initialization. The third function is a model flow creator
// that must return a single-output flow using, as output variable,
// the requested cell to put the model into.
func (sheetContent *SheetContent) PutFlow(row, column uint16, initializer func(Variable) (ModelFlow, error)) error {
	if sheetContent.locked {
		return ErrSheetContentIsLocked
	}

	cellVar := sheetContent.Cell(row, column)
	cellToUse := usedCell{variable: cellVar, row: row, column: column}
	if initializer != nil {
		if flow, err := initializer(cellVar); err != nil {
			return err
		} else if flow != nil {
			output := flow.Output()
			if len(output) != 1 || !output[cellVar] {
				return ErrFlowIsNotSingleOutput
			}
			sheetContent.flows[row * sheetContent.height + column] = flow
			sheetContent.usedCells[cellToUse] = true
		} else {
			sheetContent.flows[row * sheetContent.height + column] = nil
			delete(sheetContent.usedCells, cellToUse)
		}
	} else {
		sheetContent.flows[row * sheetContent.height + column] = nil
		delete(sheetContent.usedCells, cellToUse)
	}
	return nil
}


// Locks the sheetContent so further changes are totally disallowed
// when the initializer function ends. Returns the set of used cells,
// and the set of referenced cells that are not among the used cells,
// so they will be pre-populated with 0 on each evaluation.
func (sheetContent *SheetContent) lock() (usedCellsSet, Variables) {
	usedCells := sheetContent.usedCells
	sheetContent.usedCells = nil
	zeroCells := sheetContent.referencedCells
	sheetContent.referencedCells = nil
	sheetContent.locked = true
	for usedCell := range usedCells {
		delete(zeroCells, usedCell.variable)
	}
	return usedCells, zeroCells
}


// Gets a flow from a specific position (cell). This method is private
// and only makes sense in terms of the model execution.
func (sheetContent *SheetContent) flow(row, column uint16) ModelFlow {
	if row >= sheetContent.height || column >= sheetContent.width {
		panic(ErrOutOfBounds)
	}
	return sheetContent.flows[row * sheetContent.height + column]
}


// This is an initializer for a Spread operation. This means: not
// just a PutFlow operation, but a spread operation (which initializes
// a flow on each cell being traversed). For this to work, it takes
// 5 parameters as follows: row/column of the current cell being
// iterated, delta-row/delta-column of the current cell with respect
// to the start cell, and finally the variable of the cell being
// iterated. This is true: row == int16(int32(startingRow) deltaRow)
// and the same applies to column.
type SpreadCellInitializer func(row, column uint16, deltaRow, deltaColumn int32,
	cell Variable) (ModelFlow, error)


// This is a convenience function to fill several cells at once, like
// Excel does when dragging a cell's boundaries to "spread" (hence the
// name) a formula across several cells. The spread can be done in any
// direction, and the user is responsible of not panicking when trying
// to retrieve a cell that could be invalid (panic reason: coordinates
// out of bounds).
func (sheetContent *SheetContent) Spread(
	startRow, startColumn, endRow, endColumn uint16,
	initializer SpreadCellInitializer,
) error {
	// Consider bounds appropriately.
	if startRow >= sheetContent.height || startColumn >= sheetContent.width ||
		endRow >= sheetContent.height || endColumn >= sheetContent.width {
		panic(ErrOutOfBounds)
	}

	// Iterate forward or backward, depending on the direction.
	if startRow < endRow {
		for row := startRow; row <= endRow; row++ {
			if err := sheetContent.spreadColumns(
				row, startColumn, endColumn,
				int32(row) - int32(startRow), initializer,
			); err != nil {
				return err
			}
		}
	} else {
		for row := startRow; row >= endRow; row-- {
			if err := sheetContent.spreadColumns(
				row, startColumn, endColumn,
				int32(row) - int32(startRow),
				initializer,
			); err != nil {
				return err
			}
		}
	}
	return nil
}


// Second stage of the spread operation.
func (sheetContent *SheetContent) spreadColumns(
	row, startColumn, endColumn uint16, deltaRow int32, initializer SpreadCellInitializer,
) error {
	if startColumn < endColumn {
		for column := startColumn; column <= endColumn; column++ {
			if err := sheetContent.spreadCell(
				row, column, deltaRow,
				int32(column) - int32(startColumn),
				initializer,
			); err != nil {
				return err
			}
		}
	} else {
		for column := startColumn; column >= endColumn; column-- {
			if err := sheetContent.spreadCell(
				row, column, deltaRow,
				int32(column) - int32(startColumn),
				initializer,
			); err != nil {
				return err
			}
		}
	}
	return nil
}


// Cell-wise step of the spread operation.
func (sheetContent *SheetContent) spreadCell(
	row, column uint16, deltaRow, deltaColumn int32, initializer SpreadCellInitializer,
) error {
	return sheetContent.PutFlow(row, column, func(cell Variable) (ModelFlow, error) {
		return initializer(row, column, deltaRow, deltaColumn, cell)
	})
}


// This is a private constructor to make a new spread sheet
// content object. It will serve to be wrapped by the
// spread sheet model.
func newSheetContent(name string, width, height uint16) *SheetContent {
	if width == 0 || height == 0 {
		panic(ErrInvalidSize)
	}
	flows := make([]ModelFlow, height * width)
	return &SheetContent{
		locked:    false,
		name:      name,
		width:     width,
		height:    height,
		flows:     flows,
		usedCells: usedCellsSet{},
		referencedCells: Variables{},
	}
}
