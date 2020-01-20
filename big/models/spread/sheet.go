package spread


import (
	. "github.com/universe-10th/calculus/v2/big/models"
	. "github.com/universe-10th/calculus/v2/big/expressions"
	"math/big"
)


type EvaluationResult map[Variable]interface{}


// This is an error class giving information of a recursion
// error (this means: interdependent cells that end in an
// infinite recursion - a model under this condition will
// never be able to be solved).
type InfiniteCellRecursionError struct {
	cells Variables
}


// Implementation of the error interface: the error message.
func (infiniteCellRecursionError *InfiniteCellRecursionError) Error() string {
	return "there are cells that incur in an infinite recursion - this model cannot be solved"
}


// Tells the involved cells (they are the cells that could not be reached
// by the execution order algorithm).
func (infiniteCellRecursionError *InfiniteCellRecursionError) Cells() Variables {
	return infiniteCellRecursionError.cells
}


func newInfiniteRecursionError(cells Variables) error {
	return &InfiniteCellRecursionError{cells: cells}
}


// Spread sheet model flows allow regular argument input (i.e.
// common variables that serve as input to a model flow) and
// maintain a bi-dimensional array of per-cell models that can
// be only populated on initialization, but evaluated an
// arbitrary amount of times by filling the arguments. The
// return value of sheet evaluation is a set of "arguments"
// (they are output arguments) being the "used" cells (i.e. the
// cells that were assigned expressions to upon initialization).
type SheetModelFlow struct {
	CustomModelFlow
	content *SheetContent
	executionOrder []usedCell
	zeroCells Variables
}


func NewSheet(name string, width, height uint16, initializer func(*SheetContent)) (*SheetModelFlow, error) {
	/**
	 * First, create the sheet and run the user-provided initializer.
	 * After the initializer, the content cannot be edited anymore.
	 * It is totally discouraged to keep the reference for external
	 * uses and even worse inside some sort of goroutine.
	 */
	content := newSheetContent(name, width, height)
	initializer(content)
	usedCells, zeroCells := content.lock()

	/**
	 * The execution order algorithm will run now.
	 */

	// First, just allocate which cell variables were reached and which
	// ones were still not reached so far.
	reachedVars := Variables{}
	unreachedVars := Variables{}
	for usedCell := range usedCells {
		unreachedVars[usedCell.variable] = true
	}

	// Now, allocate the space to keep the input variables. These are
	// the variables defined as input in each of the per-cell flows but
	// belong to neither the reached nor unreached vars. This means:
	// they are not current spreadsheet variables, and so they are part
	// of the overall model input.
	overallInput := Variables{}
	// Now initialize the new structure: the execution order, and the
	// tracking of the structure's filling.
	usedCellsLength := len(usedCells)
	executionOrder := make([]usedCell, usedCellsLength)
	sortingIndex := 0
	// And now run the loop. The loop will run until:
	// - Failure: Cannot mark another used cell as reachable, and there
	//            are cells still not marked as reachable.
	// - Success: The algorithm could mark all of the cells as reachable
	//            and completely filled the execution order structure.
	for {
		foundReachable := true

		// We run over all the cells, not considering the already
		// reached ones, to try to add them to the execution order.
		// This means: trying to add at least one. If we are successful
		// in the task, then mark "found reachable" = true.
		for usedCell := range usedCells {
			// We will a priori consider the current cell as reachable,
			// and we will perform a test to check whether we must
			// revoke such status.
			reachableCell := true

			if _, unreached := unreachedVars[usedCell.variable]; unreached {
				// This is a new iteration involving a used cell not
				// being marked as reached. Now we get its flow. It
				// will always exist.
				flow := content.flow(usedCell.row, usedCell.column)

				// We must iterate through its variables and distinguish:
				// - When a variable is an already "reached" cell.
				// - When a variable is a not yet "reached" cell.
				// - When a variable is not a cell (i.e. is input).
				// We will think as follow.
				// - Each input will be added to the "overall input".
				// - If a not yet "reached" cell is found, we mark
				//   "reachable cell" = false, and break.
				input := flow.Input()
				for inputVar := range input {
					_, reached := reachedVars[inputVar]
					_, unreached := unreachedVars[inputVar]
					// Always: reached && unreached == false.
					if reached {
						// This is a reached variable. Try the next one.
						continue
					} else if unreached {
						// This is a still unreached variable. We must
						// cut here, and mark "found reachable" = false
						// so this cell is still not marked as reachable.
						reachableCell = false
						break
					} else {
						// This is not a reached / unreached variable.
						// This means: this is not a cell. It is an input
						// variable.
						overallInput[usedCell.variable] = true
					}
				}
			}

			// If the cell, after the test, was not rejected as "not yet reachable",
			// then it is a reachable cell. We must add it to the execution order
			// and also tell that a reachable cell was found in this current (overall)
			// iteration. Also, remove it from unreachedVars and add it to reachedVars.
			if reachableCell {
				foundReachable = true
				executionOrder[sortingIndex] = usedCell
				reachedVars[usedCell.variable] = true
				delete(unreachedVars, usedCell.variable)
				sortingIndex++
			}
		}

		// If we were unsuccessful, we have to return the appropriate error.
		if !foundReachable {
			return nil, newInfiniteRecursionError(unreachedVars)
		}

		// On the other hand, we found a reachable cell, and the sorting index
		// has increased
		if sortingIndex == usedCellsLength {
			// This is a successful break. We will not return anything yet, but
			// can successfully break this loop for good.
			break
		}

		// Otherwise, next iteration.
	}

	// So now we have three things that matter:
	// - The overall input variables, which will serve as spread sheet flow input.
	// - The reached variables, which will serve as spread sheet flow output.
	// - The execution order, which will serve to tell how to run this model and
	//   process the results by back-feeding the arguments into the next iteration
	//   in the execution order.
	// We can now create a custom model, since it is all valid now.
	if customFlow, err := NewCustomModelFlow(overallInput, reachedVars); err != nil {
		return nil, err
	} else {
		return &SheetModelFlow{
			CustomModelFlow: *customFlow,
			content: content,
			executionOrder: executionOrder,
			zeroCells: zeroCells,
		}, nil
	}
}


// The name of the spread sheet is got from the underlying
// spread sheet content.
func (sheet *SheetModelFlow) Name() string {
	return sheet.content.name
}


var zero = big.NewFloat(0)


// Executes the underlying spread sheet given the arguments,
// always considering the discovered execution order, and keeps
// populating the arguments with the return values, to compute
// the latter cells. In the end, returns only the cell results.
func (sheet *SheetModelFlow) Evaluate(args Arguments) (Arguments, error) {
	/**
	 * Evaluation of this sheet is to be done carefully.
	 *
	 * The first step is to get the arguments. Here, we
	 *   must completely ignore the argument keys that
	 *   are cell references in this spreadsheet (we
	 *   track this by caching the sheet.usedCellsSet).
	 *
	 * Now, with the pruned arguments
	 */

	// Init the results copy with the input.
	results := Arguments{}
	for variable, value := range args {
		results[variable] = value
	}
	// Also, put `zero` as argument in each cell that
	// has been referenced but it is not a "used" cell
	// (i.e. never called a Put/PutFlow/Spread on it).
	for variable := range sheet.zeroCells {
		results[variable] = zero
	}
	// Evaluate the model, step by step, considering
	// the appropriate execution order.
	for _, usedCell := range sheet.executionOrder {
		if result, err := sheet.content.flow(usedCell.row, usedCell.column).Evaluate(results); err != nil {
			// An error has occurred. Abort everything and return it.
			return nil, err
		} else {
			// A single output variable will have the executed model.
			// The output variable will be the used cell's variable.
			// We take the value from the result using the variable
			// as key, and put it into the current results.
			results[usedCell.variable] = result[usedCell.variable]
		}
	}
	// Prune the input from the end results.
	for variable := range args {
		delete(results, variable)
	}
	// Return the results, with no error.
	return results, nil
}