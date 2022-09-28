package raxos

type Recorder struct {
}

// instantiate a new Recorder

func NewRecorder() *Recorder {

	re := Recorder{}

	return &re
}

// start listening to gRPC connection

func (r Recorder) NetworkInit() {

}
