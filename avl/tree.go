package avl

import (
	"github.com/dploop/avl-vs-rb/base"
	"github.com/dploop/avl-vs-rb/stats"
	"github.com/dploop/avl-vs-rb/types"
)

type Tree struct {
	*base.Tree
}

func New(less types.Less) *Tree {
	return &Tree{base.New(less)}
}

func (t *Tree) Insert(z *base.Node) {
	z.Extra = Balanced
	z.Parent, z.Left, z.Right = nil, nil, nil
	x, childIsLeft := t.End(), true
	for y := x.Left; y != nil; {
		stats.InsertFindLoopCounter++
		x, childIsLeft = y, t.Less(z.Data, y.Data)
		if childIsLeft {
			y = y.Left
		} else {
			y = y.Right
		}
	}
	z.Parent = x
	if childIsLeft {
		x.Left = z
	} else {
		x.Right = z
	}
	if t.Start.Left != nil {
		t.Start = t.Start.Left
	}
	t.balanceAfterInsert(x, childIsLeft)
	t.Size++
}

func (t *Tree) balanceAfterInsert(x *base.Node, childIsLeft bool) {
	for ; x != t.End(); x = x.Parent {
		stats.InsertBalanceLoopCounter++
		if !childIsLeft {
			switch x.Extra {
			case LeftHeavy:
				x.Extra = Balanced
				return
			case RightHeavy:
				if x.Right.Extra == LeftHeavy {
					stats.InsertRotateCounter += 2
					rotateRightLeft(x)
				} else {
					stats.InsertRotateCounter++
					rotateLeft(x)
				}
				return
			default:
				x.Extra = RightHeavy
			}
		} else {
			switch x.Extra {
			case RightHeavy:
				x.Extra = Balanced
				return
			case LeftHeavy:
				if x.Left.Extra == RightHeavy {
					stats.InsertRotateCounter += 2
					rotateLeftRight(x)
				} else {
					stats.InsertRotateCounter++
					rotateRight(x)
				}
				return
			default:
				x.Extra = LeftHeavy
			}
		}
		childIsLeft = x == x.Parent.Left
	}
}

func (t *Tree) Delete(z *base.Node) {
	if t.Start == z {
		t.Start = z.Next()
	}
	x, childIsLeft := z.Parent, z == z.Parent.Left
	switch {
	case z.Left == nil:
		base.Transplant(z, z.Right)
	case z.Right == nil:
		base.Transplant(z, z.Left)
	default:
		if z.Extra == RightHeavy {
			y := base.Minimum(z.Right)
			x, childIsLeft = y, y == y.Parent.Left
			if y.Parent != z {
				x = y.Parent
				base.Transplant(y, y.Right)
				y.Right = z.Right
				y.Right.Parent = y
			}
			base.Transplant(z, y)
			y.Left = z.Left
			y.Left.Parent = y
			y.Extra = z.Extra
		} else {
			y := base.Maximum(z.Left)
			x, childIsLeft = y, y == y.Parent.Left
			if y.Parent != z {
				x = y.Parent
				base.Transplant(y, y.Left)
				y.Left = z.Left
				y.Left.Parent = y
			}
			base.Transplant(z, y)
			y.Right = z.Right
			y.Right.Parent = y
			y.Extra = z.Extra
		}
	}
	t.balanceAfterDelete(x, childIsLeft)
	t.Size--
}

func (t *Tree) balanceAfterDelete(x *base.Node, childIsLeft bool) {
	for ; x != t.End(); x = x.Parent {
		stats.DeleteBalanceLoopCounter++
		if childIsLeft {
			switch x.Extra {
			case Balanced:
				x.Extra = RightHeavy
				return
			case RightHeavy:
				b := x.Right.Extra
				if b == LeftHeavy {
					stats.DeleteRotateCounter += 2
					rotateRightLeft(x)
				} else {
					stats.DeleteRotateCounter++
					rotateLeft(x)
				}
				if b == Balanced {
					return
				}
				x = x.Parent
			default:
				x.Extra = Balanced
			}
		} else {
			switch x.Extra {
			case Balanced:
				x.Extra = LeftHeavy
				return
			case LeftHeavy:
				b := x.Left.Extra
				if b == RightHeavy {
					stats.DeleteRotateCounter += 2
					rotateLeftRight(x)
				} else {
					stats.DeleteRotateCounter++
					rotateRight(x)
				}
				if b == Balanced {
					return
				}
				x = x.Parent
			default:
				x.Extra = Balanced
			}
		}
		childIsLeft = x == x.Parent.Left
	}
}

func rotateLeft(x *base.Node) {
	z := x.Right
	x.Right = z.Left
	if z.Left != nil {
		z.Left.Parent = x
	}
	z.Parent = x.Parent
	if x == x.Parent.Left {
		x.Parent.Left = z
	} else {
		x.Parent.Right = z
	}
	z.Left = x
	x.Parent = z
	if z.Extra == Balanced {
		x.Extra, z.Extra = RightHeavy, LeftHeavy
	} else {
		x.Extra, z.Extra = Balanced, Balanced
	}
}

func rotateRight(x *base.Node) {
	z := x.Left
	x.Left = z.Right
	if z.Right != nil {
		z.Right.Parent = x
	}
	z.Parent = x.Parent
	if x == x.Parent.Right {
		x.Parent.Right = z
	} else {
		x.Parent.Left = z
	}
	z.Right = x
	x.Parent = z
	if z.Extra == Balanced {
		x.Extra, z.Extra = LeftHeavy, RightHeavy
	} else {
		x.Extra, z.Extra = Balanced, Balanced
	}
}

func rotateRightLeft(x *base.Node) {
	z := x.Right
	y := z.Left
	z.Left = y.Right
	if y.Right != nil {
		y.Right.Parent = z
	}
	y.Right = z
	z.Parent = y
	x.Right = y.Left
	if y.Left != nil {
		y.Left.Parent = x
	}
	y.Parent = x.Parent
	if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}
	y.Left = x
	x.Parent = y
	switch y.Extra {
	case RightHeavy:
		x.Extra, z.Extra = LeftHeavy, Balanced
	case LeftHeavy:
		x.Extra, z.Extra = Balanced, RightHeavy
	default:
		x.Extra, z.Extra = Balanced, Balanced
	}
	y.Extra = Balanced
}

func rotateLeftRight(x *base.Node) {
	z := x.Left
	y := z.Right
	z.Right = y.Left
	if y.Left != nil {
		y.Left.Parent = z
	}
	y.Left = z
	z.Parent = y
	x.Left = y.Right
	if y.Right != nil {
		y.Right.Parent = x
	}
	y.Parent = x.Parent
	if x == x.Parent.Right {
		x.Parent.Right = y
	} else {
		x.Parent.Left = y
	}
	y.Right = x
	x.Parent = y
	switch y.Extra {
	case LeftHeavy:
		x.Extra, z.Extra = RightHeavy, Balanced
	case RightHeavy:
		x.Extra, z.Extra = Balanced, LeftHeavy
	default:
		x.Extra, z.Extra = Balanced, Balanced
	}
	y.Extra = Balanced
}
