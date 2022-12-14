package raycast

type sprite struct {
	image     string
	pos       vector
	distance  float64
	height    float64
	animation *animation
}

func NewSprite(imageName string) *sprite {
	return &sprite{
		image:    imageName,
		pos:      vector{},
		distance: -1,
	}
}

func NewAnimatedSprite(imageName string, animation *animation) *sprite {
	animation.Play()
	return &sprite{
		image:     imageName,
		pos:       vector{},
		distance:  -1,
		animation: animation,
	}
}

type animation struct {
	numFrames    int
	currentFrame int
	numTime      float64
	currentTime  float64
	isLoop       bool
	isPlaying    bool
	isReset      bool
}

func (r *animation) Update(delta float64) {
	if !r.isPlaying {
		return
	}
	if r.numFrames == 1 || r.numFrames == 0 {
		return
	}
	r.currentTime += delta
	if r.currentTime > r.numTime {
		r.currentTime -= r.numTime
		r.currentFrame += 1
		if r.currentFrame == r.numFrames {
			if r.isLoop {
				r.currentFrame = 0
			} else {
				if r.isReset {
					r.currentFrame = 0
				}
				r.isPlaying = false
			}
		}
	}
}

func (r *animation) Play() {
	if !r.isPlaying {
		r.isPlaying = true
		r.currentFrame = 0
		r.currentTime = 0
	}
}
