package periph

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Pin struct {
	pin         int
	mode        int
	chip        int
	value       int
	value2      int
	RisingEdge  *edge
	FallingEdge *edge
}

type edge struct {
	C       chan struct{}
	repeat  bool
	running bool
}

const (
	LOW  = 0
	HIGH = 1

	FALLING_EDGE = 0
	RISING_EDGE  = 1
)

/*
	MODES
	GPIO = 1
	PWM = 2
*/

/*func main() {

	pwmPin := NewPin(0)
	defer pwmPin.Close()
	fmt.Println(pwmPin.Pwm(0))
	fmt.Println(pwmPin.Freq(1000000))
	fmt.Println(pwmPin.DutyCycle(1000000))

	bandPin1 := NewPin(23)
	defer bandPin1.Close()
	bandPin2 := NewPin(24)
	defer bandPin2.Close()
	fmt.Println(bandPin1.Output())
	fmt.Println(bandPin2.Output())

	fmt.Println(bandPin1.Low())
	fmt.Println(bandPin2.High())

	bandUpDownPin1 := NewPin(26)
	defer bandUpDownPin1.Close()
	bandUpDownPin2 := NewPin(19)
	defer bandUpDownPin2.Close()
	fmt.Println(bandUpDownPin1.Input())
	fmt.Println(bandUpDownPin2.Input())

	fmt.Println(bandUpDownPin1.Read())
	fmt.Println(bandUpDownPin2.Read())

	time.Sleep(time.Second * 5)

}*/

func NewPin(p int) *Pin {
	return &Pin{pin: p}
}

func (p *Pin) Pwm(pwmChip int) (err error) {
	if p.mode != 0 && p.mode != 2 {
		p.Close()
	}
	p.chip = pwmChip
	if _, err = os.Stat("/sys/class/pwm/pwmchip" + strconv.Itoa(int(p.chip)) + "/pwm" + strconv.Itoa(int(p.pin))); os.IsNotExist(err) {
		if err = ioutil.WriteFile("/sys/class/pwm/pwmchip"+strconv.Itoa(int(p.chip))+"/export", []byte(strconv.Itoa(int(p.pin))), 0770); err != nil {
			return err
		}
	}
	p.mode = 2
	if err = ioutil.WriteFile("/sys/class/pwm/pwmchip"+strconv.Itoa(int(p.chip))+"/pwm"+strconv.Itoa(int(p.pin))+"/enable", []byte("1"), 0770); err != nil {
		return
	}
	return nil
}

func (p *Pin) Freq(freq int) (err error) {
	if p.mode != 2 {
		err = errors.New("This pin is not pwm")
		return
	}
	if err = ioutil.WriteFile("/sys/class/pwm/pwmchip"+strconv.Itoa(int(p.chip))+"/pwm"+strconv.Itoa(int(p.pin))+"/period", []byte(strconv.Itoa(int(freq))), 0770); err != nil {
		return
	}
	if err = ioutil.WriteFile("/sys/class/pwm/pwmchip"+strconv.Itoa(int(p.chip))+"/pwm"+strconv.Itoa(int(p.pin))+"/enable", []byte("0"), 0770); err != nil {
		return
	}
	if err = ioutil.WriteFile("/sys/class/pwm/pwmchip"+strconv.Itoa(int(p.chip))+"/pwm"+strconv.Itoa(int(p.pin))+"/enable", []byte("1"), 0770); err != nil {
		return
	}
	p.value = freq
	return
}

func (p *Pin) DutyCycle(dutyLen int, cycLen int) (err error) {
	if p.mode != 2 {
		err = errors.New("This pin is not pwm")
		return
	}
	dutyLen = p.value / cycLen * dutyLen
	if err = ioutil.WriteFile("/sys/class/pwm/pwmchip"+strconv.Itoa(int(p.chip))+"/pwm"+strconv.Itoa(int(p.pin))+"/duty_cycle", []byte(strconv.Itoa(int(dutyLen))), 0770); err != nil {
		return
	}
	if err = ioutil.WriteFile("/sys/class/pwm/pwmchip"+strconv.Itoa(int(p.chip))+"/pwm"+strconv.Itoa(int(p.pin))+"/enable", []byte("0"), 0770); err != nil {
		return
	}
	if err = ioutil.WriteFile("/sys/class/pwm/pwmchip"+strconv.Itoa(int(p.chip))+"/pwm"+strconv.Itoa(int(p.pin))+"/enable", []byte("1"), 0770); err != nil {
		return
	}
	p.value2 = dutyLen
	return
}

func (p *Pin) Output() (err error) {
	if p.mode != 0 && p.mode != 1 {
		p.Close()
	}
	if _, err = os.Stat("/sys/class/gpio/gpio" + strconv.Itoa(int(p.pin))); os.IsNotExist(err) {
		if err = ioutil.WriteFile("/sys/class/gpio/export", []byte(strconv.Itoa(int(p.pin))), 0770); err != nil {
			return
		}
	}
	p.mode = 1
	err = ioutil.WriteFile("/sys/class/gpio/gpio"+strconv.Itoa(int(p.pin))+"/direction", []byte("out"), 0770)
	return
}

func (p *Pin) High() (err error) {
	if p.mode != 1 {
		err = errors.New("This pin is not gpio")
		return
	}
	if err = ioutil.WriteFile("/sys/class/gpio/gpio"+strconv.Itoa(int(p.pin))+"/value", []byte("1"), 0770); err != nil {
		return
	}
	p.value = 1
	return
}

func (p *Pin) Low() (err error) {
	if p.mode != 1 {
		return errors.New("This pin is not gpio")
	}
	if err = ioutil.WriteFile("/sys/class/gpio/gpio"+strconv.Itoa(int(p.pin))+"/value", []byte("0"), 0770); err != nil {
		return
	}
	p.value = 0
	return
}

func (p *Pin) Input() (err error) {
	if p.mode != 0 && p.mode != 1 {
		p.Close()
	}
	if _, err = os.Stat("/sys/class/gpio/gpio" + strconv.Itoa(int(p.pin))); os.IsNotExist(err) {
		if err = ioutil.WriteFile("/sys/class/gpio/export", []byte(strconv.Itoa(int(p.pin))), 0770); err != nil {
			return
		}
	}
	p.mode = 1
	err = ioutil.WriteFile("/sys/class/gpio/gpio"+strconv.Itoa(int(p.pin))+"/direction", []byte("in"), 0770)
	return
}

func (p *Pin) Read() (val int) {
	if p.mode != 1 {
		return 0 //, errors.New("This pin is not gpio")
	}
	var dat []byte
	var err error
	if dat, err = ioutil.ReadFile("/sys/class/gpio/gpio" + strconv.Itoa(p.pin) + "/value"); err != nil {
		return
	}
	val, err = strconv.Atoi(strings.Split(string(dat), "\n")[0])
	p.value = val
	return
}

func (p *Pin) FallingEdgeInit(repeat bool) {
	p.FallingEdge.edgeInit(FALLING_EDGE, repeat, p)
}

func (p *Pin) RisingEdgeInit(repeat bool) {
	p.RisingEdge.edgeInit(RISING_EDGE, repeat, p)
}

func (p *Pin) FallingEdgeClose() {
	p.FallingEdge.edgeClose()
}

func (p *Pin) RisingEdgeClose() {
	p.RisingEdge.edgeClose()
}

func (p *Pin) Close() (err error) {
	switch p.mode {
	case 0:
		return
	case 1:
		if _, err = os.Stat("/sys/class/gpio/gpio" + strconv.Itoa(int(p.pin))); !os.IsNotExist(err) {
			if err = ioutil.WriteFile("/sys/class/gpio/unexport", []byte(strconv.Itoa(int(p.pin))), 0770); err != nil {
				return
			}
		}
	case 2:
		if _, err = os.Stat("/sys/class/pwm/pwmchip" + strconv.Itoa(int(p.chip)) + "/pwm" + strconv.Itoa(int(p.pin))); !os.IsNotExist(err) {
			if err = ioutil.WriteFile("/sys/class/pwm/pwmchip"+strconv.Itoa(int(p.chip))+"/unexport", []byte(strconv.Itoa(int(p.pin))), 0770); err != nil {
				return
			}
		}
	}
	p.mode = 0
	p.chip = 0
	p.value = 0
	p.value2 = 0
	return
}

func (e *edge) edgeInit(edgeType int, repeat bool, p *Pin) {
	e.C = make(chan struct{}, 1)
	e.running = true
	lock := false
	go func() {
		for e.running {
			if p.Read() != edgeType {
				lock = true
			} else if lock {
				lock = false

				select {
				case e.C <- struct{}{}:
				default:
				}

				if !repeat {
					return
				}
			}
		}
	}()
}

func (e *edge) edgeClose() {
	e.running = false
}
