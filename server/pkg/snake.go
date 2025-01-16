package pkg

import "math/rand/v2"

var UNITARY_UP = Speed{
	UP,
	1,
}

var UNITARY_DOWN = Speed{
	DOWN,
	1,
}

var UNITARY_LEFT = Speed{
	LEFT,
	1,
}

var UNITARY_RIGHT = Speed{
	RIGHT,
	1,
}

/*
Manages the limits and positions of objects inside the square
*/
type Zone struct {
	//Max width of the square
	x int
	//Max height of the square
	y int
}

/*
Position of an element inside a zone
*/
type Position struct {
	X int
	Y int
}

// Posible moving directions
const (
	UP = iota
	DOWN
	LEFT
	RIGHT
)

/*
Direction in which an object is moving
*/
type Direction = int

type Speed struct {
	Dir     Direction
	Modulus int
}

type Snake struct {
	Positions []*Position
	Color     string
	Name      string
	HeadSpeed Speed
	Alive     bool
	Score     int
	Id        int
}

func RandomPosition() *Position {
	position := NewPosition(rand.Int(), rand.Int())
	return position
}

func RandomSpeed() Speed {
	v := rand.N(4)
	switch v {
	case 1:
		return UNITARY_RIGHT
	case 2:
		return UNITARY_DOWN
	case 3:
		return UNITARY_LEFT
	default:
		return UNITARY_UP
	}
}

func SpawnSnake(z *Zone, name string, color string, id int) *Snake {
	head := RandomPosition()
	z.IntoLimits(head)
	speed := RandomSpeed()
	positions := make([]*Position, 10)
	positions[9] = head
	for i := 8; i >= 0; i-- {
		positions[i] = NewPosition(positions[i+1].X, positions[i+1].Y)
		positions[i].MoveBy(speed)
	}
	return &Snake{
		positions,
		color,
		name,
		speed,
		true,
		0,
		id,
	}
}

/*
	Increases the score of the snake and
	adds the new position associated with it
*/
func (s *Snake) Grow(z *Zone) {
	s.Score += 10
	newX := s.Positions[len(s.Positions)-1].X - s.Positions[len(s.Positions)-2].X
	newY := s.Positions[len(s.Positions)-1].Y - s.Positions[len(s.Positions)-2].Y
	position := NewPosition(newX, newY)
	z.IntoLimits(position)
	s.Positions = append(s.Positions, position)
}

/*
	Moves a Snake based on its speed
*/
func (s *Snake) Move(z *Zone, food *Position) bool {
	if !s.Alive && len(s.Positions) > 0 {
		s.Positions = s.Positions[1:]
		return false
	}
	if !s.Alive {
		return false
	}
	head := NewPosition(s.Positions[0].X, s.Positions[0].Y)
	head.MoveBy(s.HeadSpeed)
	z.IntoLimits(head)
	for i := len(s.Positions) - 1; i > 0; i-- {
		s.Positions[i] = s.Positions[i-1]
		s.Alive = s.Alive && !head.Collided(s.Positions[i])
	}
	s.Positions[0] = head
	if s.Collided(food) {
		s.Grow(z)
	}
	return s.Collided(food)
}

/*
	Changes a snake speed
*/
func (s *Snake) ChangeSpeed(newSpeed Speed) {
	s.HeadSpeed = newSpeed
}

func (p *Position) Collided(o *Position) bool {
	return p.X == o.X && p.Y == o.Y
}

/*
	Checks if the head of the snake has collided with an object
*/
func (s *Snake) Collided(o *Position) bool {
	return s.Positions[0].Collided(o)
}

func NewZone(x, y int) *Zone {
	return &Zone{
		x,
		y,
	}
}

func NewPosition(x, y int) *Position {
	return &Position{
		x,
		y,
	}
}

func MultiplySpeed(s Speed, by int) Speed {
	return Speed{
		s.Dir,
		s.Modulus * by,
	}
}

/*
Changes a position to respect the limits of the zone
by doing a roundabout
*/
func (z *Zone) IntoLimits(p *Position) {
	for p.X < 0 {
		p.X += z.x
	}
	for p.Y < 0 {
		p.Y += z.y
	}
	p.X = p.X % z.x
	p.Y = p.Y % z.y
}

/*
Moves a position in a given direction by a given speed
*/
func (p *Position) MoveBy(s Speed) {
	switch s.Dir {
	case UP:
		p.Y -= s.Modulus
	case DOWN:
		p.Y += s.Modulus
	case RIGHT:
		p.X += s.Modulus
	case LEFT:
		p.X -= s.Modulus
	}
}
