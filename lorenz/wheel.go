/*

https://en.wikipedia.org/wiki/Malkus_waterwheel

*/

package lorenz

import (
	"fmt"
	"math"
	"strings"
)

const nBuckets = 8

const tau = math.Pi * 2.0

// bucket on the lorenz wheel
type bucket struct {
	amount   float64 // amount of water
	capacity float64 // maximum capacity
}

// add water to a bucket
func (b *bucket) add(x float64) {
	b.amount += x
	if b.amount > b.capacity {
		b.amount = b.capacity
	}
	if b.amount < 0 {
		b.amount = 0
	}
}

// Wheel stores the state for a Lorenz Wheel.
type Wheel struct {
	iflow  float64          // inflow of water to top bucket (mass/time)
	oflow  float64          // outlow of water from all buckets (mass/time)
	theta  float64          // wheel position for bucket 0 (radians)
	av     float64          // angular velocity
	arc    float64          // radians per bucket
	bucket [nBuckets]bucket // water in bucket (mass)
}

// New returns a Lorenz Wheel
func New(iflow, oflow float64) *Wheel {

	w := Wheel{
		iflow: iflow,
		oflow: oflow,
		arc:   tau / float64(nBuckets),
		av:    0.01,
	}

	for k := range w.bucket {
		w.bucket[k].capacity = 10.0 * iflow
	}

	return &w
}

func (w *Wheel) String() string {
	n := w.topBucket()
	s := make([]string, nBuckets)
	for i, b := range w.bucket {
		top := " "
		if i == n {
			top = "*"
		}
		s[i] = fmt.Sprintf("%s%.2f", top, b.amount)
	}
	return fmt.Sprintf("theta %.2f av %.2f: %s", w.theta, w.av, strings.Join(s, " "))
}

func (w *Wheel) topBucket() int {
	k := math.Mod(w.theta, tau)
	if k < 0 {
		k += tau
	}
	return int(math.Trunc(k / w.arc))
}

// Run an iteration of the Lorenz Wheel.
func (w *Wheel) Run(dt float64) {

	// add water to top bucket
	n := w.topBucket()
	w.bucket[n].add(w.iflow * dt)

	// remove water from all buckets
	for k := range w.bucket {
		w.bucket[k].add(-w.oflow * dt)
	}

	// work out angular inertia and torque
	ai := 0.0
	torque := 0.0
	angle := w.theta
	for k := range w.bucket {
		ai += w.bucket[k].amount
		torque += w.bucket[k].amount * math.Sin(angle)
		angle += w.arc
	}

	// angular acceleration
	aa := torque / ai
	// new angular velocity
	w.av += aa * dt
	// new position
	w.theta += w.av * dt
}
