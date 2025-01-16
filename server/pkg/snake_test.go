package pkg_test

import (
	"testing"
	"websocket/pkg"
)

func TestMovingPosition(t *testing.T) {
	position := pkg.NewPosition(50, 50)

	position.MoveBy(pkg.UNITARY_LEFT)

	if position.X != 49 && position.Y != 50 {
		t.Fatalf("Invalid x position")
	}
	position.MoveBy(pkg.UNITARY_UP)

	if position.Y != 49 && position.X != 49 {
		t.Fatalf("Invalid Y position")
	}

	position.MoveBy(pkg.UNITARY_RIGHT)
	if position.X != 50 && position.Y != 49 {
		t.Fatalf("Invalid x position 2")
	}

	position.MoveBy(pkg.UNITARY_DOWN)
	if position.Y != 50 && position.X != 50 {
		t.Fatalf("Invalid Y position")
	}
}

func TestMovingBySpeed(t *testing.T) {
	position := pkg.NewPosition(10, 10)
	speed := pkg.MultiplySpeed(pkg.UNITARY_DOWN, 10)

	position.MoveBy(speed)
	if position.Y != 20 {
		t.Fatalf("Invalid Y position")
	}

	position.MoveBy(pkg.MultiplySpeed(pkg.UNITARY_RIGHT, 10))
	if position.X != 20 {
		t.Fatalf("Invalid X position")
	}

	position.MoveBy(pkg.MultiplySpeed(pkg.UNITARY_UP, 20))
	if position.Y != 0 {
		t.Fatalf("Invalid Y position 2")
	}

	position.MoveBy(pkg.MultiplySpeed(pkg.UNITARY_LEFT, 20))
	if position.X != 0 {
		t.Fatalf("Invalid X position 2")
	}

}

func TestMoveInsideAZoneX(t *testing.T) {
	position := pkg.NewPosition(95, 0)
	speed := pkg.MultiplySpeed(pkg.UNITARY_RIGHT, 6)

	zone := pkg.NewZone(100, 100)

	position.MoveBy(speed)
	zone.IntoLimits(position)

	if position.X != 1 {
		t.Fatalf("Invalid X")
	}

	speed = pkg.MultiplySpeed(pkg.UNITARY_LEFT, 2)
	position.MoveBy(speed)
	zone.IntoLimits(position)

	if position.X != 99 {
		t.Fatalf("Invalid X on left move, expected %d but got %d", 99, position.X)
	}

}

func TestMoveInsideAZoneY(t *testing.T) {
	position := pkg.NewPosition(0, 95)
	speed := pkg.MultiplySpeed(pkg.UNITARY_DOWN, 6)

	zone := pkg.NewZone(100, 100)

	position.MoveBy(speed)
	zone.IntoLimits(position)

	if position.Y != 1 {
		t.Fatalf("Invalid Y, expected %d but got %d", 1, position.Y)
	}

	speed = pkg.MultiplySpeed(pkg.UNITARY_UP, 2)
	position.MoveBy(speed)
	zone.IntoLimits(position)

	if position.Y != 99 {
		t.Fatalf("Invalid X on left move, expected %d but got %d", 99, position.Y)
	}

}
