package ui

import (
	"math"
	"math/rand"

	"gamejam/util"

	"github.com/hajimehoshi/ebiten/v2"
)

// menuAntSize is the drawn size (square) of a decorative menu ant.
const menuAntSize = 32

// menuAntSpacing is the distance, in pixels, kept between ants marching in a
// single-file column.
const menuAntSpacing = 55

// menuRoachChance is the probability that a decorative menu bug spawns as a
// roach instead of an ant.
const menuRoachChance = 0.25

// Walk spritesheets for the decorative menu bugs, loaded once and shared across
// every menu ant/roach so we don't re-decode the PNG for each spawn.
var (
	menuAntWalkSheet   *ebiten.Image
	menuRoachWalkSheet *ebiten.Image
)

func menuBugSheet(roach bool) *ebiten.Image {
	if roach {
		if menuRoachWalkSheet == nil {
			menuRoachWalkSheet = util.LoadImage("units/roaches/roach-walk.png")
		}
		return menuRoachWalkSheet
	}
	if menuAntWalkSheet == nil {
		menuAntWalkSheet = util.LoadImage("units/ants/ant-walk.png")
	}
	return menuAntWalkSheet
}

// menuAnt is a single decorative ant with a position, facing and walk cycle.
// It has no autonomous behaviour of its own; a MenuAntColumn positions it each
// frame, either by wandering (the leader) or by following the ant ahead.
type menuAnt struct {
	anim   *SpriteAnimation
	x, y   float64 // centre position, in logical screen pixels
	angle  float64 // facing, radians (0 = travelling along +X)
	moving bool
}

func newMenuAnt() *menuAnt {
	sheet := menuBugSheet(rand.Float64() < menuRoachChance)
	return &menuAnt{
		anim: NewSpriteAnimation(sheet, 128, 128, 4, 4, true),
	}
}

// draw renders the ant, rotated to face its direction of travel. The source
// sprite faces "up" (head towards -Y), so we offset the heading by 90 degrees.
func (a *menuAnt) draw(screen *ebiten.Image) {
	// Only advance the walk cycle while the ant is actually moving.
	if a.moving {
		a.anim.Update(1)
	}
	frame := a.anim.CurrentFrameImage()
	fw := float64(a.anim.FrameWidth)
	fh := float64(a.anim.FrameHeight)

	op := &ebiten.DrawImageOptions{}
	scale := menuAntSize / fw
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(-menuAntSize/2, -fh*scale/2)
	op.GeoM.Rotate(a.angle + math.Pi/2)
	op.GeoM.Translate(a.x, a.y)
	screen.DrawImage(frame, op)
}

// MenuAntColumn is a line of ants marching single-file across the menu
// background. A leader wanders and bounces off the screen edges, laying down a
// breadcrumb trail; the followers walk along that same trail at a fixed spacing
// behind, so the whole column snakes around like a real ant line.
type MenuAntColumn struct {
	ants []*menuAnt

	// leader movement state
	lx, ly     float64 // leader centre
	vx, vy     float64 // leader velocity (pixels per update)
	steerTimer float64

	// trail is a breadcrumb history of recent leader positions, newest first.
	trail []point
}

type point struct{ x, y float64 }

// NewMenuAntColumn creates a marching column of the given length confined to
// the screen size.
func NewMenuAntColumn(length, screenW, screenH int) *MenuAntColumn {
	if length < 1 {
		length = 1
	}

	speed := 0.8 + rand.Float64()*0.9
	angle := rand.Float64() * 2 * math.Pi

	c := &MenuAntColumn{
		lx:         rand.Float64() * float64(screenW),
		ly:         rand.Float64() * float64(screenH),
		vx:         math.Cos(angle) * speed,
		vy:         math.Sin(angle) * speed,
		steerTimer: 1 + rand.Float64()*2,
	}
	for i := 0; i < length; i++ {
		c.ants = append(c.ants, newMenuAnt())
	}
	// Seed the trail so followers have somewhere to stand on the first frame.
	c.trail = append(c.trail, point{c.lx, c.ly})
	return c
}

// Update advances the leader and repositions the followers along the trail.
func (c *MenuAntColumn) Update(screenW, screenH int) {
	// Occasionally nudge the leader's heading so the line curves organically.
	c.steerTimer -= 1.0 / 60.0
	if c.steerTimer <= 0 {
		c.steerTimer = 1 + rand.Float64()*2
		speed := math.Hypot(c.vx, c.vy)
		angle := math.Atan2(c.vy, c.vx) + (rand.Float64()-0.5)*math.Pi/2
		c.vx = math.Cos(angle) * speed
		c.vy = math.Sin(angle) * speed
	}

	c.lx += c.vx
	c.ly += c.vy

	half := float64(menuAntSize) / 2
	w, h := float64(screenW), float64(screenH)
	if c.lx < half {
		c.lx = half
		c.vx = math.Abs(c.vx)
	} else if c.lx > w-half {
		c.lx = w - half
		c.vx = -math.Abs(c.vx)
	}
	if c.ly < half {
		c.ly = half
		c.vy = math.Abs(c.vy)
	} else if c.ly > h-half {
		c.ly = h - half
		c.vy = -math.Abs(c.vy)
	}

	// Record the new leader position at the front of the trail.
	c.trail = append([]point{{c.lx, c.ly}}, c.trail...)

	// We need enough trail to place every follower at its spacing. Cap the
	// history a little beyond that to avoid unbounded growth.
	maxDist := menuAntSpacing * float64(len(c.ants))
	c.trimTrail(maxDist)

	// Position the leader ant.
	leaderAngle := math.Atan2(c.vy, c.vx)
	c.ants[0].x = c.lx
	c.ants[0].y = c.ly
	c.ants[0].angle = leaderAngle
	c.ants[0].moving = true

	// Position each follower a fixed arc-length back along the trail.
	for i := 1; i < len(c.ants); i++ {
		targetDist := menuAntSpacing * float64(i)
		px, py, ang, ok := c.pointAlongTrail(targetDist)
		if !ok {
			// Not enough trail yet: stack on the last known point.
			px, py, ang = c.lx, c.ly, leaderAngle
		}
		c.ants[i].x = px
		c.ants[i].y = py
		c.ants[i].angle = ang
		c.ants[i].moving = true
	}
}

// pointAlongTrail walks back along the breadcrumb trail a given arc-length and
// returns the position there plus the facing (pointing towards the leader).
func (c *MenuAntColumn) pointAlongTrail(dist float64) (x, y, angle float64, ok bool) {
	acc := 0.0
	for i := 0; i < len(c.trail)-1; i++ {
		a := c.trail[i]
		b := c.trail[i+1]
		seg := math.Hypot(b.x-a.x, b.y-a.y)
		if seg == 0 {
			continue
		}
		if acc+seg >= dist {
			t := (dist - acc) / seg
			x = a.x + (b.x-a.x)*t
			y = a.y + (b.y-a.y)*t
			// Face towards the newer point (a), i.e. towards the leader.
			angle = math.Atan2(a.y-b.y, a.x-b.x)
			return x, y, angle, true
		}
		acc += seg
	}
	return 0, 0, 0, false
}

// trimTrail drops breadcrumbs once the accumulated length exceeds maxDist.
func (c *MenuAntColumn) trimTrail(maxDist float64) {
	acc := 0.0
	for i := 0; i < len(c.trail)-1; i++ {
		acc += math.Hypot(c.trail[i+1].x-c.trail[i].x, c.trail[i+1].y-c.trail[i].y)
		if acc > maxDist {
			c.trail = c.trail[:i+2]
			return
		}
	}
}

// Draw renders every ant in the column.
func (c *MenuAntColumn) Draw(screen *ebiten.Image) {
	for _, a := range c.ants {
		a.draw(screen)
	}
}
