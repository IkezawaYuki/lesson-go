package main

import "fmt"

type Shape interface {
	Render() string
}

type Circle struct {
	Radius float32
}

func (c *Circle) Render() string {
	return fmt.Sprintf("Circle of radius %f", c.Radius)
}

func (c *Circle) Resize(factor float32) {
	c.Radius *= factor
}

type Square struct {
	Side float32
}

func (s *Square) Render() string {
	return fmt.Sprintf("Square with side %f", s.Side)
}

type ColoredShape struct {
	Shape Shape
	Color string
}

func (c *ColoredShape) Render() string {
	return fmt.Sprintf("%s has the color %s",
		c.Shape.Render(), c.Color)
}

type TransparentShape struct {
	Shape            Shape
	TransparentShape float32
}

func (t *TransparentShape) Render() string {
	return fmt.Sprintf("%s has %f%% transparency",
		t.Shape.Render(), t.TransparentShape*100.0)
}

func main() {
	circle := Circle{2}
	fmt.Println(circle.Render())

	redCircle := ColoredShape{
		Shape: &circle,
		Color: "Red",
	}
	fmt.Println(redCircle.Render())

	thsCircle := TransparentShape{&redCircle, 0.5}
	fmt.Println(thsCircle.Render())
}
