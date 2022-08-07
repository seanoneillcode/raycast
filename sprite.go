package raycast

type sprite struct {
	image     string
	pos       vector
	distance  float64
	height    float64
	animation *animation
}

type animation struct {
	numFrames    int
	currentFrame int
	numTime      float64
	currentTime  float64
}

func (r *animation) Update(delta float64) {
	r.currentTime += delta
	if r.currentTime > r.numTime {
		r.currentTime -= r.numTime
		//if r.currentTime > r.numTime {
		//	r.currentTime = 0
		//}
		r.currentFrame += 1
		if r.currentFrame == r.numFrames {
			r.currentFrame = 0
		}
	}
}
