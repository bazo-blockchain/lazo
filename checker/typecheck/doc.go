// The package typecheck encapsulates the type checking phase of the checker.
// It performs different type checks on the contract. Examples are:
// - Not Operator (!) can only be applied to boolean expressions
// - Both sites of an assignment are of the same type
// - If Condition always needs to be a boolean expression
// - a.s.o.
package typecheck
